package feishu

import (
	"container/list"
	"context"
	"github.com/Mmx233/daoUtil"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/backoff"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func NewUserSyncBackoff(stat redis.SyncStat, schedule cron.Schedule) backoff.Backoff {
	return backoff.New(backoff.Conf{
		Content: stat.Inject(schedule, func() error {
			sync := UserSyncProcessor{}
			err := sync.Run()
			if err != nil {
				log.Errorf("同步飞书用户列表失败: %v", err)
			} else {
				log.Infof("飞书用户列表同步成功，差异处理耗时 %dms", sync.Cost.Milliseconds())
				sync.PrintSyncResult()
			}

			return err
		}),
		MaxRetryDelay: time.Minute * 60,
	})
}

type UserSyncProcessor struct {
	tx *gorm.DB

	Cost time.Duration

	createdUser  int
	unFrozenUser int
	frozenUser   int

	createdUserGroup int
	deletedUserGroup int
}

type RelatedUserInfo struct {
	Data        dao.User
	Departments []uint
}

func (a *UserSyncProcessor) Run() error {
	userDataRaw, err := Api.LoadUserList()
	if err != nil {
		return err
	}

	var startAt = time.Now()

	// 去除重复用户与无效用户，转换数据为操作结构体
	var userMapLength int
	for _, userSlice := range userDataRaw {
		userMapLength += len(userSlice)
	}
	var userMap = make(map[string]struct{}, userMapLength)
	var userDataList list.List
	for _, userSlice := range userDataRaw {
		for _, userData := range userSlice {
			if _, ok := userMap[userData.UserId]; !ok {
				userMap[userData.UserId] = struct{}{}
				userData := userData
				newUser := NewUser(&userData)
				if !newUser.IsInvalid() {
					userDataList.PushBack(newUser)
				}
			}
		}
	}
	userList := make([]*User, userDataList.Len())
	for i, el := 0, userDataList.Front(); el != nil; i, el = i+1, el.Next() {
		userList[i] = el.Value.(*User)
	}

	// 开启事务，开始数据库操作
	a.tx = dao.DB.Begin()
	defer a.tx.Rollback()

	// 获取飞书部门 OpenID 与基础组的映射关系
	groupMap, err := (service.FeishuGroupsSrv{DB: a.tx}).GetGroupMap(daoUtil.LockForShare)

	// 锁定 User UserGroup 表
	err = a.tx.Session(&gorm.Session{
		PrepareStmt: false,
	}).Exec("LOCK TABLES `users` WRITE, `user_groups` WRITE").Error
	if err != nil {
		return err
	}

	// 同步用户列表
	if err = a.doSyncUsers(userList); err != nil {
		return err
	}

	// 同步用户组
	if err = a.doSyncUserGroups(userList, groupMap); err != nil {
		return err
	}

	// 提交事务
	if err = a.tx.Commit().Error; err != nil {
		return err
	}

	// 计算同步耗时
	a.Cost = time.Now().Sub(startAt)
	return nil
}

func (a *UserSyncProcessor) PrintSyncResult() {
	log.Infof("创建用户 %d 个，冻结用户 %d 个，解冻用户 %d 个，添加用户组 %d 个，移除用户组 %d 个",
		a.createdUser, a.frozenUser, a.unFrozenUser, a.createdUserGroup, a.deletedUserGroup)
}

// 数据库操作：创建不存在的用户，解冻已冻结用户，冻结不在列表中的用户
func (a *UserSyncProcessor) doSyncUsers(userList []*User) error {
	// 读取已存在用户
	var allPhone = make([]string, len(userList))
	for i, user := range userList {
		allPhone[i] = user.Data.Mobile
	}
	userSrv := service.UserSrv{DB: a.tx}
	// 此处 sql 指定了返回数据顺序与输入号码数组顺序一致，见 dao 函数
	existUsers, err := userSrv.GetUnscopedUserByPhoneSlice(allPhone)
	if err != nil {
		return err
	}

	// 冻结不在列表但在数据库中未冻结的用户
	invalidUsers, err := userSrv.GetUserNotInPhoneSlice(allPhone)
	if err != nil {
		return err
	}
	if len(invalidUsers) > 0 {
		var invalidUID = make([]uint, len(invalidUsers))
		for i, user := range invalidUsers {
			invalidUID[i] = user.ID
		}
		if err = userSrv.FrozeByIDSlice(invalidUID); err != nil {
			return err
		}
		redisJwt := redis.NewUserJwt()
		for _, uid := range invalidUID {
			err = redisJwt.NewOperator(uid).Del(context.Background())
			if err != nil {
				return err
			}
		}
		a.frozenUser = len(invalidUID)
	}

	// 对比数据。此处不考虑手机号重复的情况，届时将 panic
	// 此处依赖 []*User 中的元素为指针
	a.createdUser = len(allPhone) - len(existUsers)
	var userToCreate []*User
	var userToUnFroze list.List // uint
	if a.createdUser > 0 {
		userToCreate = make([]*User, a.createdUser)
	}
	if len(existUsers) == 0 {
		for i, user := range userList {
			userToCreate[i] = user
		}
	} else {
		var exUserIndex int
		var userToCreateIndex int
		userModel := &existUsers[0]
		for _, user := range userList {
			if userModel != nil && user.Data.Mobile == userModel.Phone {
				user.ID = userModel.ID
				if userModel.DeletedAt.Valid {
					userToUnFroze.PushBack(userModel.ID)
				}

				exUserIndex++
				if exUserIndex >= len(existUsers) {
					userModel = nil
				} else {
					userModel = &existUsers[exUserIndex]
				}
			} else {
				userToCreate[userToCreateIndex] = user
				userToCreateIndex++
			}
		}
	}

	// 将对比结果写入数据库
	if len(userToCreate) > 0 {
		userModelToCreate := make([]dao.User, len(userToCreate))
		for i, user := range userToCreate {
			userModelToCreate[i] = user.Model()
		}
		if err = userSrv.CreateAll(userModelToCreate); err != nil {
			return err
		}
		// 回填新用户 uid
		for i, userModel := range userModelToCreate {
			userToCreate[i].ID = userModel.ID
		}
	}
	if userToUnFroze.Len() > 0 {
		a.unFrozenUser = userToUnFroze.Len()
		var idSlice = make([]uint, userToUnFroze.Len())
		for i, el := 0, userToUnFroze.Front(); el != nil; i, el = i+1, el.Next() {
			idSlice[i] = el.Value.(uint)
		}
		if err = userSrv.UnFrozeByIDSlice(idSlice); err != nil {
			return err
		}
	}

	return nil
}

// 数据库操作：同步用户部门关系
func (a *UserSyncProcessor) doSyncUserGroups(userList []*User, groupMap map[string]uint) error {
	userGroupSrv := service.UserGroupsSrv{DB: a.tx}
	// 此处数据已按照 uid,gid 排序
	existUserGroups, err := userGroupSrv.GetAll()
	if err != nil {
		return err
	}
	var exUserGroupMap map[uint][]dao.UserGroups

	// 处理特殊情况
	if len(existUserGroups) == 0 {
		var length int
		var modelSlice = make([][]dao.UserGroups, len(userList))
		for i, user := range userList {
			userDepartmentModels := user.Departments(groupMap).Models(user.ID)
			length += len(userDepartmentModels)
			modelSlice[i] = userDepartmentModels
		}
		a.createdUserGroup = length
		models := make([]dao.UserGroups, length)
		length = 0
		for _, modelSliceEl := range modelSlice {
			for _, userGroup := range modelSliceEl {
				models[length] = userGroup
				length++
			}
		}
		return userGroupSrv.CreateAll(models)
	}

	// 转换数据库数据为 uid 映射
	exUserGroupMap = make(map[uint][]dao.UserGroups, len(userList)-a.createdUser-a.frozenUser)
	var beginIndex int
	var lastUID uint
	lastUID = existUserGroups[0].UID
	for i, userGroup := range existUserGroups {
		if userGroup.UID != lastUID {
			exUserGroupMap[lastUID] = existUserGroups[beginIndex:i]
			beginIndex = i
			lastUID = userGroup.UID
		}
	}
	exUserGroupMap[lastUID] = existUserGroups[beginIndex:]

	// 计算差异
	var userGroupsToAdd = list.New()    // dao.UserGroups
	var userGroupsToDelete = list.New() // uint
	redisJwt := redis.NewUserJwt()
	for _, user := range userList {
		thisUserExistGroups := exUserGroupMap[user.ID]
		currentUserGroups := user.Departments(groupMap).Ids()
		var userGroupChanged bool
		for _, gid := range currentUserGroups {
			for _, exUserGroup := range thisUserExistGroups {
				if gid == exUserGroup.GID {
					goto nextUserGroup
				}
			}
			userGroupsToAdd.PushBack(dao.UserGroups{
				UID: user.ID,
				GID: gid,
			})
			userGroupChanged = true
		nextUserGroup:
		}
		for _, exUserGroup := range thisUserExistGroups {
			for _, gid := range currentUserGroups {
				if exUserGroup.GID == gid {
					goto nextExUserGroup
				}
			}
			userGroupsToDelete.PushBack(exUserGroup.ID)
			userGroupChanged = true
		nextExUserGroup:
		}
		if userGroupChanged {
			redisUserOperator := redisJwt.NewOperator(user.ID)
			exist, err := redisUserOperator.Exist(context.Background())
			if err != nil {
				return err
			} else if exist {
				_, err = redisUserOperator.ChangeOperateID(context.Background())
				if err != nil {
					return err
				}
			}
		}
	}

	// 存储计算结果
	if userGroupsToAdd.Len() != 0 {
		a.createdUserGroup = userGroupsToAdd.Len()
		userGroupsToAddModels := make([]dao.UserGroups, userGroupsToAdd.Len())
		for i, el := 0, userGroupsToAdd.Front(); el != nil; i, el = i+1, el.Next() {
			userGroupsToAddModels[i] = el.Value.(dao.UserGroups)
		}
		if err = userGroupSrv.CreateAll(userGroupsToAddModels); err != nil {
			return err
		}
	}
	if userGroupsToDelete.Len() != 0 {
		a.deletedUserGroup = userGroupsToDelete.Len()
		userGroupsToDeleteSlice := make([]uint, userGroupsToDelete.Len())
		for i, el := 0, userGroupsToDelete.Front(); el != nil; i, el = i+1, el.Next() {
			userGroupsToDeleteSlice[i] = el.Value.(uint)
		}
		if err = userGroupSrv.DeleteByIDSlice(userGroupsToDeleteSlice); err != nil {
			return err
		}
	}
	return nil
}
