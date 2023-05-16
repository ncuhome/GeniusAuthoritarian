package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

var BaseGroups = BaseGroupsSrv{dao.DB}

type BaseGroupsSrv struct {
	*gorm.DB
}

func (a BaseGroupsSrv) Begin() (BaseGroupsSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a BaseGroupsSrv) LoadGroups() ([]dao.BaseGroup, error) {
	return (&dao.BaseGroup{}).GetAll(a.DB)
}

func (a BaseGroupsSrv) LoadGroupsForShow() ([]dto.Group, error) {
	return (&dao.BaseGroup{}).GetAllForShow(a.DB)
}

func (a BaseGroupsSrv) CreateGroups(groups []string) ([]dao.BaseGroup, error) {
	return (&dao.BaseGroup{}).CreateGroups(a.DB, groups)
}
