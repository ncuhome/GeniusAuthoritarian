package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
)

type UserSync interface {
	IsInvalid() bool
	Departments(groupMap map[string]uint) []uint
	Model() *dao.User
}

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
		Name:  u.Data.Name,
		Phone: u.Data.Mobile,
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

func (u User) DepartmentModels(uid uint, groupMap map[string]uint) []dao.UserGroups {
	departments := u.Departments(groupMap)
	models := make([]dao.UserGroups, len(departments))
	for i, gid := range departments {
		models[i].UID = uid
		models[i].GID = gid
	}
	return models
}
