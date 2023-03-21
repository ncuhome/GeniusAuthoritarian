package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"gorm.io/gorm"
	"strings"
)

var SiteWhiteList = SiteWhiteListSrv{dao.DB}

type SiteWhiteListSrv struct {
	*gorm.DB
}

func (a SiteWhiteListSrv) Begin() (*SiteWhiteListSrv, error) {
	a.DB = a.DB.Begin()
	return &a, a.Error
}

func (a SiteWhiteListSrv) Exist(domain string) (bool, error) {
	list, e := redis.SiteWhiteList.Load()
	if e != nil {
		if e == redis.Nil {
			list, e = (&dao.SiteWhiteList{}).Get(a.DB)
			if e != nil {
				return false, e
			}
			_ = redis.SiteWhiteList.Add(list...)
		} else {
			return false, e
		}
	}

	for _, v := range list {
		if strings.HasSuffix(domain, v) {
			return true, nil
		}
	}
	return false, nil
}
