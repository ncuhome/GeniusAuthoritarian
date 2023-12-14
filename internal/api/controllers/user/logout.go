package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func Logout(c *gin.Context) {
	loginID := tools.GetUserInfo(c).ID
	err := redis.NewRecordedToken().NewStorePoint(loginID).Destroy(context.Background())
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
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

	uid := tools.GetUserInfo(c).UID
	exist, err := service.LoginRecord.OnlineRecordExist(uid, f.ID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if !exist {
		callback.Error(c, callback.ErrTargetDeviceOffline)
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

	callback.Default(c)
}
