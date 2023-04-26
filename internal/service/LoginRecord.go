package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"gorm.io/gorm"
)

var LoginRecord = LoginRecordSrv{dao.DB}

type LoginRecordSrv struct {
	*gorm.DB
}

func (a LoginRecordSrv) Begin() (*LoginRecordSrv, error) {
	a.DB = a.DB.Begin()
	return &a, a.Error
}

func (a LoginRecordSrv) Add(uid uint, target string) error {
	return (&dao.LoginRecord{
		UID:    uid,
		Target: target,
	}).Insert(a.DB)
}