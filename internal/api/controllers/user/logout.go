package controllers

import (
	"context"
	"errors"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func Logout(c *gin.Context) {
	loginID := tools.GetUserInfo(c).ID

	loginRecordSrv, err := service.LoginRecord.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer loginRecordSrv.Rollback()

	err = loginRecordSrv.SetDestroyed(uint(loginID))
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	err = redis.NewRecordedToken().NewStorePoint(loginID).Destroy(context.Background())
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	if err = loginRecordSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}

func LogoutDevice(c *gin.Context) {
	var f struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	loginRecordSrv, err := service.LoginRecord.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer loginRecordSrv.Rollback()

	uid := tools.GetUserInfo(c).UID
	exist, err := loginRecordSrv.OnlineRecordExist(uid, f.ID, daoUtil.LockForUpdate)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if !exist {
		callback.Error(c, callback.ErrTargetDeviceOffline)
		return
	}

	err = loginRecordSrv.SetDestroyed(f.ID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	err = redis.NewRecordedToken().NewStorePoint(uint64(f.ID)).Destroy(context.Background())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			callback.Error(c, callback.ErrTargetDeviceOffline)
			return
		}
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	if err = loginRecordSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}
