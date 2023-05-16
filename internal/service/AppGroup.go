package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

var AppGroup = AppGroupSrv{dao.DB}

type AppGroupSrv struct {
	*gorm.DB
}

func (a AppGroupSrv) Begin() (AppGroupSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a AppGroupSrv) BindForApp(aid uint, groupIds []uint) ([]dto.Group, error) {
	groups, e := (&dao.BaseGroup{}).GetByIdsForShow(a.DB, groupIds...)
	if e != nil {
		return nil, e
	} else if len(groupIds) != len(groups) {
		return nil, gorm.ErrRecordNotFound
	}

	var toCreate = make([]dao.AppGroup, len(groupIds))
	for i, gid := range groupIds {
		toCreate[i].GID = gid
		toCreate[i].AID = aid
	}
	return groups, a.DB.Create(&toCreate).Error
}

func (a AppGroupSrv) UnBindForApp(aid uint, groupIds []uint) error {
	return (&dao.AppGroup{AID: aid}).DeleteByGidForApp(a.DB, groupIds...)
}

func (a AppGroupSrv) DeleteAllForApp(aid uint) error {
	return (&dao.AppGroup{AID: aid}).DeleteByAID(a.DB)
}
