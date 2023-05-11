package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
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

func (a UserGroupsSrv) GetForUser(uid uint) ([]string, error) {
	var groups []string
	return groups, (&dao.UserGroups{UID: uid}).GetUserGroupsByUID(a.DB).Select("groups.name").Find(&groups).Error
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
