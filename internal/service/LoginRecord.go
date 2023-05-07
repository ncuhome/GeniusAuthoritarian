package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
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

func (a LoginRecordSrv) Add(uid uint, ip, target string) error {
	return (&dao.LoginRecord{
		UID:    uid,
		IP:     ip,
		Target: target,
	}).Insert(a.DB)
}

func (a LoginRecordSrv) UserHistory(uid uint, limit int) ([]dto.LoginRecord, error) {
	return (&dao.LoginRecord{
		UID: uid,
	}).GetByUID(a.DB, limit)
}
