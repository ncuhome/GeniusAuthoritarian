package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
	"time"
)

var LoginRecord = LoginRecordSrv{dao.DB}

type LoginRecordSrv struct {
	*gorm.DB
}

func (a LoginRecordSrv) Begin() (*LoginRecordSrv, error) {
	a.DB = a.DB.Begin()
	return &a, a.Error
}

func (a LoginRecordSrv) Add(uid, appID uint, ip, useragent string, tokenValid time.Duration) (uint, error) {
	var validBefore uint64
	if tokenValid != 0 {
		validBefore = uint64(time.Now().Add(tokenValid).Unix())
	}
	model := dao.LoginRecord{
		UID:         uid,
		Useragent:   useragent,
		IP:          ip,
		ValidBefore: validBefore,
	}
	if appID != 0 {
		model.AID = &appID
	}
	return model.ID, model.Insert(a.DB)
}

func (a LoginRecordSrv) UserHistory(uid uint, limit int) ([]dto.LoginRecord, error) {
	return (&dao.LoginRecord{
		UID: uid,
	}).GetByUID(a.DB, limit)
}

func (a LoginRecordSrv) GetViewIDs(aid, startAt uint) ([]uint, error) {
	return (&dao.LoginRecord{
		AID: &aid,
	}).GetViewIds(a.DB, startAt)
}

func (a LoginRecordSrv) GetMultipleViewsIDs(apps []dao.App) ([]dto.ViewID, error) {
	return (&dao.LoginRecord{}).GetMultipleViewsIds(a.DB, apps)
}
