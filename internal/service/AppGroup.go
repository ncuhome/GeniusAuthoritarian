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

func (a AppGroupSrv) GetGroupsByAppCode(appCode string) ([]string, error) {
	return (&dao.AppGroup{}).GetGroups(a.DB, appCode)
}
