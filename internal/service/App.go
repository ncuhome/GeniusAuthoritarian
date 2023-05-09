package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"gorm.io/gorm"
)

var App = AppSrv{dao.DB}

type AppSrv struct {
	*gorm.DB
}

func (a AppSrv) Begin() (AppSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a AppSrv) Exist(appCode string) (bool, error) {
	list, e := redis.AppCode.Load()
	if e != nil {
		if e == redis.Nil {
			list, e = (&dao.App{}).Get(a.DB)
			if e != nil {
				return false, e
			}
			_ = redis.AppCode.Add(list...)
		} else {
			return false, e
		}
	}

	for _, v := range list {
		if v == appCode {
			return true, nil
		}
	}
	return false, nil
}

func (a AppSrv) CheckAppCode(appCode string) (bool, error) {
	return a.Exist(appCode)
}

func (a AppSrv) FistAppForLogin(appCode string) (*dao.App, error) {
	var t = dao.App{
		AppCode: appCode,
	}
	return &t, t.FirstForLogin(a.DB)
}

func (a AppSrv) GetSecretByAppCode(appCode string) (string, error) {
	return (&dao.App{}).FindSecret(a.DB, appCode)
}
