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

func (a UserGroupsSrv) GetAll() ([]dao.UserGroupModel, error) {
	return (&dao.UserGroups{}).GetAllUnfrozen(a.DB)
}

func (a UserGroupsSrv) CreateAll(data []dao.UserGroupModel) error {
	return (&dao.UserGroups{}).InsertAll(a.DB, data)
}

func (a UserGroupsSrv) DeleteByIDSlice(id []uint) error {
	return (&dao.UserGroups{}).DeleteByIDSlice(a.DB, id)
}
