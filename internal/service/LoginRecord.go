package service

import (
	"context"
	"errors"
	"github.com/Mmx233/daoUtil"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
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

func (a LoginRecordSrv) Add(uid, appID uint, ip, useragent, method string, tokenValid time.Duration) (uint, error) {
	var validBefore uint64
	if tokenValid != 0 {
		validBefore = uint64(time.Now().Add(tokenValid).Unix())
	}
	model := dao.LoginRecord{
		UID:         uid,
		ValidBefore: validBefore,
		Useragent:   useragent,
		IP:          ip,
		Method:      method,
	}
	if appID != 0 {
		model.AID = &appID
	}
	return model.ID, model.Insert(a.DB)
}

func (a LoginRecordSrv) SetDestroyed(id uint) error {
	return (&dao.LoginRecord{ID: id}).UpdateDestroyedByID(a.DB)
}

func (a LoginRecordSrv) SetDestroyedByIDS(ids []uint) error {
	return (&dao.LoginRecord{}).UpdateDestroyedByIDSlice(a.DB, ids)
}

func (a LoginRecordSrv) GetValidForApp(aid uint, opt ...daoUtil.ServiceOpt) ([]uint, error) {
	return (&dao.LoginRecord{AID: &aid}).GetIdByAID(daoUtil.TxOpts(a.DB, opt...))
}

func (a LoginRecordSrv) UserHistory(uid uint, limit int) ([]dto.LoginRecord, error) {
	return (&dao.LoginRecord{
		UID: uid,
	}).GetByUID(a.DB, limit)
}

func (a LoginRecordSrv) UserOnline(uid uint, currentLoginID uint) ([]dto.LoginRecordOnline, error) {
	validRecords, err := (&dao.LoginRecord{UID: uid}).GetValidForUser(a.DB)
	if err != nil {
		return nil, err
	}

	var _redis = redis.NewRecordedToken()
	var validCount int
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for i, record := range validRecords {
		err = _redis.NewStorePoint(uint64(record.ID)).Get(ctx, nil)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				validRecords[i].ID = 0
				continue
			}
			return nil, err
		}
		validCount++
	}

	var result = make([]dto.LoginRecordOnline, validCount)
	if validCount != 0 {
		for i, record := range validRecords {
			if record.ID != 0 {
				if record.ID == currentLoginID {
					record.IsMe = true
				}
				result[i] = record
			}
		}
	}
	return result, nil
}

func (a LoginRecordSrv) OnlineRecordExist(uid, id uint, opts ...daoUtil.ServiceOpt) (bool, error) {
	return (&dao.LoginRecord{
		ID:  id,
		UID: uid,
	}).ValidExist(daoUtil.TxOpts(a.DB, opts...))
}

func (a LoginRecordSrv) GetViewIDs(aid, startAt uint) ([]uint, error) {
	return (&dao.LoginRecord{
		AID: &aid,
	}).GetViewIds(a.DB, startAt)
}

func (a LoginRecordSrv) GetMultipleViewsIDs(apps []dao.App) ([]dto.ViewID, error) {
	return (&dao.LoginRecord{}).GetMultipleViewsIds(a.DB, apps)
}

func (a LoginRecordSrv) GetForAdminView(startTime time.Time) (*dto.AdminLoginDataView, error) {
	records, err := (&dao.LoginRecord{}).GetAdminViews(a.DB, startTime.Unix())
	if err != nil {
		return nil, err
	}

	appIdMap := make(map[uint]struct{}, 4)
	for _, record := range records {
		if record.AID == 0 {
			continue
		}
		appIdMap[record.AID] = struct{}{}
	}
	appIds := make([]uint, len(appIdMap))
	i := 0
	for id := range appIdMap {
		appIds[i] = id
		i++
	}

	apps, err := (&dao.App{}).GetDataViewByIds(a.DB, appIds...)
	if err != nil {
		return nil, err
	}
	apps = append(apps, dto.AppDataView{
		ID:   0,
		Name: global.ThisAppName,
	})

	return &dto.AdminLoginDataView{
		Apps:    apps,
		Records: records,
	}, nil
}
