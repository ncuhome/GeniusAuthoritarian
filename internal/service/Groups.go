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

func (a GroupsSrv) LoadGroups() ([]dao.BaseGroup, error) {
	return (&dao.BaseGroup{}).GetAll(a.DB)
}

func (a GroupsSrv) LoadGroupsForShow() ([]dto.Group, error) {
	return (&dao.BaseGroup{}).GetAllForShow(a.DB)
}

func (a GroupsSrv) CreateGroups(groups []string) ([]dao.BaseGroup, error) {
	return (&dao.BaseGroup{}).CreateGroups(a.DB, groups)
}
