package feishu

import (
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	log "github.com/sirupsen/logrus"
)

func UserSync() error {
	userData, e := Api.LoadUserList()
	if e != nil {
		return e
	}

	feishuGroupsSrv, e := service.FeishuGroups.Begin()
	if e != nil {
		return e
	}
	defer feishuGroupsSrv.Rollback()

	// 过滤无效组
	var openID = make([]string, len(userData))
	i := 0
	for k := range userData {
		openID[i] = k
		i++
	}
	validGroups, e := feishuGroupsSrv.GetByOpenID(openID, daoUtil.LockForShare)
	if e != nil {
		return e
	}
	var validGroupsMap = make(map[string]uint, len(validGroups))
	for _, group := range validGroups {
		validGroupsMap[group.OpenDepartmentId] = group.GID
	}
	var invalidOpenID []string
	for k := range userData {
		if _, ok := validGroupsMap[k]; ok {
			goto nextGroup
		}
		invalidOpenID = append(invalidOpenID, k)
	nextGroup:
	}
	for _, key := range invalidOpenID {
		delete(userData, key)
	}

	// 过滤无效用户
	for k, users := range userData {
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
		userData[k] = filteredList
	}

	// 转换数据为 dao.Group.ID ==> []dao.User
	var filteredData = make(map[uint][]dao.User, len(userData))
	for openID, users := range userData {
		dbUserList := make([]dao.User, len(users))
		for i, user := range users {
			dbUserList[i] = dao.User{
				Name:  user.Name,
				Phone: user.Mobile,
			}
		}
		filteredData[validGroupsMap[openID]] = dbUserList
	}

	// 反转映射关系，以 dao.User.Phone 为 key
	type UserInfo struct {
		Data        dao.User
		Departments []uint
	}
	var lens int
	for _, users := range filteredData {
		lens += len(users)
	}
	var reserveData = make(map[string]*UserInfo, lens)
	for gid, users := range filteredData {
		for _, user := range users {
			if _, ok := reserveData[user.Phone]; ok {
				reserveData[user.Phone].Departments = append(reserveData[user.Phone].Departments, gid)
			} else {
				reserveData[user.Phone] = &UserInfo{
					Data:        user,
					Departments: []uint{gid},
				}
			}
		}
	}

	// 数据库操作：创建不存在的用户，解冻已冻结用户，冻结不在列表中的用户
	var allPhone = make([]string, len(reserveData))
	i = 0
	for phone := range reserveData {
		allPhone[i] = phone
		i++
	}
	userSrv := service.UserSrv{DB: feishuGroupsSrv.DB}
	existUsers, e := userSrv.GetUnscopedUserByPhoneSlice(allPhone)
	if e != nil {
		return e
	}
	lensToCreate := len(allPhone) - len(existUsers)
	lensToUnfroze := 0
	for _, exUser := range existUsers {
		if exUser.DeletedAt.Valid {
			lensToUnfroze++
		}
		reserveData[exUser.Phone].Data.ID = exUser.ID
	}
	if lensToCreate > 0 {
		var userToCreate = make([]dao.User, lensToCreate)
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
		for _, user := range userToCreate { // 回填 UID
			reserveData[user.Phone].Data.ID = user.ID
		}
	}
	if lensToUnfroze > 0 {
		var userToUnfroze = make([]uint, lensToUnfroze)
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
	}

	// 数据库操作：同步用户部门关系
	userGroupSrv := service.UserGroupsSrv{DB: feishuGroupsSrv.DB}
	existUserGroups, e := userGroupSrv.GetAll()
	if e != nil {
		return e
	}
	var userGroupsToAdd []dao.UserGroups
	var userGroupsToDelete []uint
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
			userGroupsToAdd = append(userGroupsToAdd, dao.UserGroups{
				UID: user.Data.ID,
				GID: userDepartment,
			})
		nextUserDepartment:
		}
		for _, exUserDepartment := range exUserGroupMap[user.Data.ID] {
			for _, userDepartment := range user.Departments {
				if userDepartment == exUserDepartment {
					goto nextExUserDepartment
				}
			}
			userGroupsToDelete = append(userGroupsToDelete, exUserDepartment)
		nextExUserDepartment:
		}
	}
	if len(userGroupsToAdd) > 0 {
		if e = userGroupSrv.CreateAll(userGroupsToAdd); e != nil {
			return e
		}
	}
	if len(userGroupsToDelete) > 0 {
		if e = userGroupSrv.DeleteByIDSlice(userGroupsToDelete); e != nil {
			return e
		}
	}

	return feishuGroupsSrv.Commit().Error
}

func AddUserSyncCron(spec string) error {
	_, e := agent.AddRegular(&agent.Event{
		T: spec,
		E: func() {
			defer tool.Recover()
			if e := UserSync(); e != nil {
				log.Errorf("同步飞书用户列表失败: %v", e)
			} else {
				log.Infoln("飞书用户列表同步成功")
			}
		},
	})
	return e
}
