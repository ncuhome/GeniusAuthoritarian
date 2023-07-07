package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/GroupOperator"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	"gorm.io/gorm"
)

var UserGroups = UserGroupsSrv{dao.DB}

type UserGroupsSrv struct {
	*gorm.DB
}

func (a UserGroupsSrv) Begin() (UserGroupsSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

// IsCenterMember 是否是中心组组员
func (a UserGroupsSrv) IsCenterMember(uid uint) (bool, error) {
	return (&dao.UserGroups{
		ID:  GroupOperator.GroupRelation[departments.UCe],
		UID: uid,
	}).Exist(a.DB)
}

func (a UserGroupsSrv) GetForUser(uid uint) ([]string, error) {
	return (&dao.UserGroups{UID: uid}).GetUserGroupNamesByUID(a.DB)
}

func (a UserGroupsSrv) GetAll() ([]dao.UserGroups, error) {
	return (&dao.UserGroups{}).GetAllUnfrozen(a.DB)
}

func (a UserGroupsSrv) CreateAll(data []dao.UserGroups) error {
	return (&dao.UserGroups{}).InsertAll(a.DB, data)
}

func (a UserGroupsSrv) DeleteByIDSlice(id []uint) error {
	return (&dao.UserGroups{}).DeleteByIDSlice(a.DB, id)
}
