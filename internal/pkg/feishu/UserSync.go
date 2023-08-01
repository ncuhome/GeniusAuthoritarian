package feishu

import (
	"container/list"
	"github.com/Mmx233/daoUtil"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/backoff"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

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
	GroupOpenIdToFeishuUserSliceMap, e := a.downloadUserList()
	if e != nil {
		return e
	}

	var startAt = time.Now()
	a.tx = dao.DB.Begin()
	defer a.tx.Rollback()

	// 过滤无效数据
	validFeishuGroupIdMap, e := a.filterInvalidGroups(GroupOpenIdToFeishuUserSliceMap)
	if e != nil {
		return e
	}
	a.filterInvalidUsers(GroupOpenIdToFeishuUserSliceMap)

	// 转换数据
	GroupDaoIdToFeishuUserSliceMap := a.convertApiDataToGroupIdKeyMap(GroupOpenIdToFeishuUserSliceMap, validFeishuGroupIdMap)
	UserPhoneToRelatedUserInfoMap := a.convertReverseMap(GroupDaoIdToFeishuUserSliceMap)

	// 将数据同步入数据库
	if e = a.doSyncUsers(UserPhoneToRelatedUserInfoMap); e != nil {
		return e
	}
	if e = a.doSyncUserGroups(UserPhoneToRelatedUserInfoMap); e != nil {
		return e
	}
	if e = a.tx.Commit().Error; e != nil {
		return e
	}

	a.Cost = time.Now().Sub(startAt)
	return nil
}

func (a *UserSyncProcessor) PrintSyncResult() {
	log.Infof("创建用户 %d 个，解冻用户 %d 个，冻结用户 %d 个，添加用户组 %d 个，移除用户组 %d 个",
		a.createdUser, a.frozenUser, a.unFrozenUser, a.createdUserGroup, a.deletedUserGroup)
}

func (a *UserSyncProcessor) downloadUserList() (map[string][]feishuApi.ListUserContent, error) {
	return Api.LoadUserList()
}

// 过滤无效组，并返回有效组映射 飞书 OpenID ==> dao.BaseGroup.ID
func (a *UserSyncProcessor) filterInvalidGroups(feishuUserList map[string][]feishuApi.ListUserContent) (map[string]uint, error) {
	var openID = make([]string, len(feishuUserList))
	i := 0
	for k := range feishuUserList {
		openID[i] = k
		i++
	}
	validGroups, e := service.FeishuGroupsSrv{DB: a.tx}.GetByOpenID(openID, daoUtil.LockForShare)
	if e != nil {
		return nil, e
	}
	var validGroupsMap = make(map[string]uint, len(validGroups))
	for _, group := range validGroups {
		validGroupsMap[group.OpenDepartmentId] = group.GID
	}
	var invalidOpenID = list.New() // string
	for k := range feishuUserList {
		if _, ok := validGroupsMap[k]; ok {
			goto nextGroup
		}
		invalidOpenID.PushBack(k)
	nextGroup:
	}
	el := invalidOpenID.Front()
	for el != nil {
		delete(feishuUserList, el.Value.(string))
		el = el.Next()
	}
	return validGroupsMap, nil
}

// 过滤无效用户
func (a *UserSyncProcessor) filterInvalidUsers(feishuUserList map[string][]feishuApi.ListUserContent) {
	for k, users := range feishuUserList {
		var lens int
		for _, user := range users {
			if !user.Status.IsActivated || user.Status.IsFrozen || !user.MobileVisible {
				user.MobileVisible = false
			} else {
				lens++
			}
		}
		var filteredList = make([]feishuApi.ListUserContent, lens)
		lens = 0
		for _, user := range users {
			if user.MobileVisible {
				filteredList[lens] = user
				lens++
			}
		}
		feishuUserList[k] = filteredList
	}
}

// 转换数据为 dao.BaseGroup.ID ==> []dao.User
func (a *UserSyncProcessor) convertApiDataToGroupIdKeyMap(feishuUserList map[string][]feishuApi.ListUserContent, validGroupsMap map[string]uint) map[uint][]dao.User {
	var filteredData = make(map[uint][]dao.User, len(feishuUserList))
	for openID, users := range feishuUserList {
		dbUserList := make([]dao.User, len(users))
		for i, user := range users {
			dbUserList[i] = dao.User{
				Name:  user.Name,
				Phone: user.Mobile,
			}
		}
		filteredData[validGroupsMap[openID]] = dbUserList
	}
	return filteredData
}

// 反转映射关系，以 dao.User.Phone 为 key
func (a *UserSyncProcessor) convertReverseMap(filteredData map[uint][]dao.User) map[string]*RelatedUserInfo {
	var lens int
	for _, users := range filteredData {
		lens += len(users)
	}
	var reserveData = make(map[string]*RelatedUserInfo, lens)
	for gid, users := range filteredData {
		for _, user := range users {
			if _, ok := reserveData[user.Phone]; ok {
				reserveData[user.Phone].Departments = append(reserveData[user.Phone].Departments, gid)
			} else {
				reserveData[user.Phone] = &RelatedUserInfo{
					Data:        user,
					Departments: []uint{gid},
				}
			}
		}
	}
	return reserveData
}

// 数据库操作：创建不存在的用户，解冻已冻结用户，冻结不在列表中的用户
func (a *UserSyncProcessor) doSyncUsers(reserveData map[string]*RelatedUserInfo) error {
	var allPhone = make([]string, len(reserveData))
	i := 0
	for phone := range reserveData {
		allPhone[i] = phone
		i++
	}
	userSrv := service.UserSrv{DB: a.tx}
	existUsers, e := userSrv.GetUnscopedUserByPhoneSlice(allPhone)
	if e != nil {
		return e
	}
	a.createdUser = len(allPhone) - len(existUsers)
	for _, exUser := range existUsers {
		if exUser.DeletedAt.Valid {
			a.unFrozenUser++
		}
		reserveData[exUser.Phone].Data.ID = exUser.ID
	}
	if a.createdUser > 0 {
		var userToCreate = make([]dao.User, a.createdUser)
		i = 0
		for _, phone := range allPhone {
			for _, exUser := range existUsers {
				if phone == exUser.Phone {
					goto nextUser
				}
			}
			userToCreate[i] = reserveData[phone].Data
			i++
		nextUser:
		}
		if e = userSrv.CreateAll(userToCreate); e != nil {
			return e
		}
		for _, user := range userToCreate { // 回填 Uid
			reserveData[user.Phone].Data.ID = user.ID
		}
	}
	if a.unFrozenUser > 0 {
		var userToUnfroze = make([]uint, a.unFrozenUser)
		i = 0
		for _, exUser := range existUsers {
			if exUser.DeletedAt.Valid {
				userToUnfroze[i] = exUser.ID
				i++
			}
		}
		if e = userSrv.UnFrozeByIDSlice(userToUnfroze); e != nil {
			return e
		}
	}

	invalidUsers, e := userSrv.GetUserNotInPhoneSlice(allPhone)
	if e != nil {
		return e
	}
	if len(invalidUsers) > 0 {
		var invalidUID = make([]uint, len(invalidUsers))
		for i, user := range invalidUsers {
			delete(reserveData, user.Phone)
			invalidUID[i] = user.ID
		}
		if e = userSrv.FrozeByIDSlice(invalidUID); e != nil {
			return e
		}
		a.frozenUser = len(invalidUID)
	}
	return nil
}

// 数据库操作：同步用户部门关系
func (a *UserSyncProcessor) doSyncUserGroups(reserveData map[string]*RelatedUserInfo) error {
	userGroupSrv := service.UserGroupsSrv{DB: a.tx}
	existUserGroups, e := userGroupSrv.GetAll()
	if e != nil {
		return e
	}
	var userGroupsToAdd = list.New()
	var userGroupsToDelete = list.New() // uint
	var exUserGroupMap = make(map[uint][]uint, len(reserveData))
	for _, exUserGroup := range existUserGroups {
		exUserGroupMap[exUserGroup.UID] = append(exUserGroupMap[exUserGroup.UID], exUserGroup.GID)
	}
	for _, user := range reserveData {
		for _, userDepartment := range user.Departments {
			exGroups, ok := exUserGroupMap[user.Data.ID]
			if ok {
				for _, exGroup := range exGroups {
					if userDepartment == exGroup {
						goto nextUserDepartment
					}
				}
			}
			userGroupsToAdd.PushBack(dao.UserGroups{
				UID: user.Data.ID,
				GID: userDepartment,
			})
			_ = redis.UserJwt.Clear(user.Data.ID)
		nextUserDepartment:
		}
		for _, exUserDepartment := range exUserGroupMap[user.Data.ID] {
			for _, userDepartment := range user.Departments {
				if userDepartment == exUserDepartment {
					goto nextExUserDepartment
				}
			}
			userGroupsToDelete.PushBack(exUserDepartment)
			e = redis.UserJwt.Clear(user.Data.ID)
			if e != nil && e != redis.Nil {
				return e
			}
		nextExUserDepartment:
		}
	}
	if userGroupsToAdd.Len() != 0 {
		a.createdUserGroup = userGroupsToAdd.Len()
		userGroupsToAddModels := make([]dao.UserGroups, userGroupsToAdd.Len())
		el := userGroupsToAdd.Front()
		i := 0
		for el != nil {
			userGroupsToAddModels[i] = el.Value.(dao.UserGroups)
			i++
			el = el.Next()
		}
		if e = userGroupSrv.CreateAll(userGroupsToAddModels); e != nil {
			return e
		}
	}
	if userGroupsToDelete.Len() != 0 {
		a.deletedUserGroup = userGroupsToDelete.Len()
		userGroupsToDeleteSlice := make([]uint, userGroupsToDelete.Len())
		el := userGroupsToDelete.Front()
		i := 0
		for el != nil {
			userGroupsToDeleteSlice[i] = el.Value.(uint)
			i++
			el = el.Next()
		}
		if e = userGroupSrv.DeleteByIDSlice(userGroupsToDeleteSlice); e != nil {
			return e
		}
	}
	return nil
}

func AddUserSyncCron(spec string) error {
	worker := backoff.New(backoff.Conf{
		Content: func() error {
			sync := UserSyncProcessor{}
			e := sync.Run()
			if e != nil {
				log.Errorf("同步飞书用户列表失败: %v", e)
			} else {
				log.Infof("飞书用户列表同步成功，差异处理耗时 %dms", sync.Cost.Milliseconds())
				sync.PrintSyncResult()
			}

			return e
		},
		MaxRetryDelay: time.Minute * 60,
	})

	_, e := agent.AddRegular(&agent.Event{
		T: spec,
		E: worker.Start,
	})
	return e
}
