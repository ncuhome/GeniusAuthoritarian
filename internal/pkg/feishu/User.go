package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
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

// Model 注意需要同时修改下方的 IsModelEmpty
func (u User) Model() dao.User {
	return dao.User{
		Name:      u.Data.Name,
		Phone:     u.Data.Mobile,
		AvatarUrl: u.Data.Avatar.AvatarOrigin,
	}
}
func (u User) IsModelEmpty() bool {
	return u.Data.Name == "" &&
		u.Data.Mobile == "" &&
		u.Data.Avatar.AvatarOrigin == ""
}

func (u User) IsInvalid() bool {
	return !u.Data.Status.IsActivated || u.Data.Status.IsFrozen || u.Data.Status.IsResigned || u.Data.Mobile == "" || // 账号未异常
		u.Data.EmployeeType != 1 // 仅允许正式员工状态账号
}

func (u User) Departments(groupMap map[string]uint) UserDepartment {
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
	return UserDepartment{
		GidList:    departments[:validLength],
		OpenIdList: u.Data.DepartmentIds,
	}
}
