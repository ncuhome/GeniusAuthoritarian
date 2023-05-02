package feishu

import (
	"github.com/Mmx233/daoUtil"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserSyncProcessor struct {
	tx *gorm.DB
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

	return nil
}

func (a *UserSyncProcessor) downloadUserList() (map[string][]feishuApi.ListUserContent, error) {
	return Api.LoadUserList()
}

// 过滤无效组，并返回有效组映射 飞书 OpenID ==> dao.Group.ID
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
	var invalidOpenID []string
	for k := range feishuUserList {
		if _, ok := validGroupsMap[k]; ok {
			goto nextGroup
		}
		invalidOpenID = append(invalidOpenID, k)
	nextGroup:
	}
	for _, key := range invalidOpenID {
		delete(feishuUserList, key)
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

// 转换数据为 dao.Group.ID ==> []dao.User
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
		for _, user := range userToCreate { // 回填 Uid
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
	return nil
}

// 数据库操作：同步用户部门关系
func (a *UserSyncProcessor) doSyncUserGroups(reserveData map[string]*RelatedUserInfo) error {
	userGroupSrv := service.UserGroupsSrv{DB: a.tx}
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
	return nil
}

func AddUserSyncCron(spec string) error {
	_, e := agent.AddRegular(&agent.Event{
		T: spec,
		E: func() {
			defer tool.Recover()
			sync := UserSyncProcessor{}
			if e := sync.Run(); e != nil {
				log.Errorf("同步飞书用户列表失败: %v", e)
			} else {
				log.Infoln("飞书用户列表同步成功")
			}
		},
	})
	return e
}
