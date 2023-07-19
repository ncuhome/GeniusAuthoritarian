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

func (a BaseGroupsSrv) LoadGroupsRelation() (map[string]dao.BaseGroup, error) {
	groups, e := a.LoadGroups()
	if e != nil {
		return nil, e
	}

	var groupRelations = make(map[string]dao.BaseGroup, len(groups))
	for _, group := range groups {
		groupRelations[group.Name] = group
	}
	return groupRelations, nil
}

func (a BaseGroupsSrv) LoadGroupsForShow() ([]dto.Group, error) {
	return (&dao.BaseGroup{}).GetAllForShow(a.DB)
}

func (a BaseGroupsSrv) CreateGroups(groups []dao.BaseGroup) error {
	return (&dao.BaseGroup{}).CreateGroups(a.DB, groups)
}
