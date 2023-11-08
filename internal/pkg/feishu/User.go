package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/feishuApi"
)

type UserSync interface {
	IsInvalid() bool
	Departments(groupMap map[string]uint) []uint
}

func NewUser(data *feishuApi.User) *User {
	return &User{
		data: data,
	}
}

type User struct {
	data        *feishuApi.User
	departments []uint
}

func (u User) IsInvalid() bool {
	return !u.data.Status.IsActivated || u.data.Status.IsFrozen || u.data.Status.IsResigned || u.data.Mobile == "" || // 账号未异常
		u.data.EmployeeType != 1 // 仅允许正式员工状态账号
}

func (u User) Departments(groupMap map[string]uint) []uint {
	var departments = make([]uint, len(u.data.DepartmentIds))
	var validLength int
	for _, groupOpenID := range u.data.DepartmentIds {
		id, ok := groupMap[groupOpenID]
		if !ok {
			continue
		}
		departments[validLength] = id
		validLength++
	}
	return departments[:validLength]
}
