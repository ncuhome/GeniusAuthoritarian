package service

import (
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var App = AppSrv{dao.DB}

type AppSrv struct {
	*gorm.DB
}

func (a AppSrv) Begin() (AppSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a AppSrv) This(host string) *dao.App {
	return &dao.App{
		Name:           global.ThisAppName,
		Callback:       fmt.Sprintf("https://%s/login", host),
		PermitAllGroup: true,
	}
}

func (a AppSrv) New(uid uint, name, callback string, permitAll bool) (*dao.App, error) {
	randSrc := rand.NewSource(time.Now().UnixNano())
	var t = dao.App{
		Name:           name,
		UID:            uid,
		AppCode:        tool.RandString(randSrc, 8),
		AppSecret:      tool.RandString(randSrc, 100),
		Callback:       callback,
		PermitAllGroup: permitAll,
	}
	return &t, t.Insert(a.DB)
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

func (a AppSrv) FirstAppKeyPair(id uint) (string, string, error) {
	var t = dao.App{
		ID: id,
	}
	return t.AppCode, t.AppSecret, t.FirstAppKeyPairByID(a.DB)
}

func (a AppSrv) NameExist(name string) (bool, error) {
	return (&dao.App{Name: name}).NameExist(a.DB)
}
