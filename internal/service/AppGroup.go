package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
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

func (a AppGroupSrv) BindForApp(aid uint, groups []uint) error {
	groupIds, e := (&dao.Group{}).GetIdsByIds(a.DB, groups...)
	if e != nil {
		return e
	} else if len(groupIds) != len(groups) {
		return gorm.ErrRecordNotFound
	}

	var toCreate = make([]dao.AppGroup, len(groupIds))
	for i, gid := range groupIds {
		toCreate[i].GID = gid
		toCreate[i].AID = aid
	}
	return a.DB.Create(&toCreate).Error
}
