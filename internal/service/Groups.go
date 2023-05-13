package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

var Groups = GroupsSrv{dao.DB}

type GroupsSrv struct {
	*gorm.DB
}

func (a GroupsSrv) Begin() (GroupsSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a GroupsSrv) LoadGroups() ([]dao.Group, error) {
	var t []dao.Group
	return t, (&dao.Group{}).GetAll(a.DB).Find(&t).Error
}

func (a GroupsSrv) LoadGroupsForShow() ([]dto.Group, error) {
	var t = make([]dto.Group, 0)
	return t, (&dao.Group{}).GetAll(a.DB).Find(&t).Error
}

func (a GroupsSrv) CreateGroups(groups []string) ([]dao.Group, error) {
	return (&dao.Group{}).CreateGroups(a.DB, groups)
}
