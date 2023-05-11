package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"gorm.io/gorm"
	"net/url"
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
	list, e := (&dao.SiteWhiteList{}).Get(a.DB)
	if e != nil {
		return false, e
	}

	for _, v := range list {
		if strings.HasSuffix(domain, v) {
			return true, nil
		}
	}
	return false, nil
}

func (a SiteWhiteListSrv) CheckUrl(link string) (bool, error) {
	u, e := url.Parse(link)
	if e != nil {
		return false, e
	}

	return a.Exist(u.Hostname())
}
