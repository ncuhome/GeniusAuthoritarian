package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
	"gorm.io/gorm"
)

func NewUser(data *feishuApi.User) *User {
	return &User{
		Data: data,
	}
}

type User struct {
	ID   uint // 辅助字段，初始为空
	Data *feishuApi.User
}

func (u User) Model() dao.User {
	return dao.User{
		Name:      u.Data.Name,
		Phone:     u.Data.Mobile,
		AvatarUrl: u.Data.Avatar.AvatarOrigin,
	}
}

func (u User) IsInvalid() bool {
	return !u.Data.Status.IsActivated || u.Data.Status.IsFrozen || u.Data.Status.IsResigned || u.Data.Mobile == "" || // 账号未异常
		u.Data.EmployeeType != 1 // 仅允许正式员工状态账号
}

func (u User) Departments(groupMap map[string]uint) []uint {
	var departments = make([]uint, len(u.Data.DepartmentIds))
	var validLength int
	for _, groupOpenID := range u.Data.DepartmentIds {
		id, ok := groupMap[groupOpenID]
		if !ok {
			continue
		}
		departments[validLength] = id
		validLength++
	}
	return departments[:validLength]
}

func (u User) genDepartmentModels(uid uint, departments []uint) []dao.UserGroups {
	models := make([]dao.UserGroups, len(departments))
	for i, gid := range departments {
		models[i].UID = uid
		models[i].GID = gid
	}
	return models
}

func (u User) DepartmentModels(uid uint, groupMap map[string]uint) []dao.UserGroups {
	return u.genDepartmentModels(uid, u.Departments(groupMap))
}

func (u User) SyncDepartments(tx *gorm.DB, uid uint, groupMap map[string]uint) (changed bool, err error) {
	departments := u.Departments(groupMap)

	userGroupSrv := service.UserGroupsSrv{DB: tx}
	result := userGroupSrv.DeleteNotInGidSliceByUID(uid, departments)
	err = result.Error
	if err != nil {
		return
	}
	if result.RowsAffected != 0 {
		changed = true
	}

	existDepartments, err := userGroupSrv.GetForUser(uid)
	if err != nil {
		return
	}

	gidToAddLength := len(departments) - len(existDepartments)
	if gidToAddLength > 0 {
		changed = true
		var gidToAdd = make([]uint, gidToAddLength)
		err = userGroupSrv.CreateAll(u.genDepartmentModels(uid, gidToAdd))
		if err != nil {
			return
		}
	}
	return
}
