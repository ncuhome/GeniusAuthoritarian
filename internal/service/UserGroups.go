package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
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
		UID: uid,
	}).ExistByName(a.DB, departments.UCe)
}

func (a UserGroupsSrv) GetForUser(uid uint) ([]string, error) {
	return (&dao.UserGroups{UID: uid}).GetUserGroupNamesByUID(a.DB)
}

func (a UserGroupsSrv) GetForAppCode(uid uint, appCode string) ([]string, error) {
	var groups []string
	return groups, (&dao.UserGroups{
		UID: uid,
	}).GetUserGroupsForAppCodeByUID(a.DB, appCode).Find(&groups).Error
}

func (a UserGroupsSrv) GetAll() ([]dao.UserGroups, error) {
	return (&dao.UserGroups{}).GetAllNotFrozen(a.DB)
}

func (a UserGroupsSrv) CreateAll(data []dao.UserGroups) error {
	return (&dao.UserGroups{}).InsertAll(a.DB, data)
}

func (a UserGroupsSrv) DeleteByIDSlice(id []uint) error {
	return (&dao.UserGroups{}).DeleteByIDSlice(a.DB, id)
}

func (a UserGroupsSrv) DeleteNotInGidSliceByUID(uid uint, id []uint) *gorm.DB {
	return (&dao.UserGroups{UID: uid}).DeleteNotInGidSliceByUID(a.DB, id)
}
