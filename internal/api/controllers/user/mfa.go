package controllers

import (
	"bytes"
	"encoding/base64"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"image/jpeg"
	"time"
)

func MfaAdd(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer userSrv.Rollback()

	exist, e := userSrv.MfaExist(uid, daoUtil.LockForShare)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if exist {
		callback.Error(c, callback.ErrMfaAlreadyExist, e)
		return
	}

	mfaKey, e := tools.NewMfa(uid)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	qrcodeImage, e := mfaKey.Image(300, 300)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	var qrcodeBuffer = bytes.Buffer{}
	if e = jpeg.Encode(&qrcodeBuffer, qrcodeImage, &jpeg.Options{
		Quality: 100,
	}); e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	if e = redis.MfaEnable.Set(uid, mfaKey.Secret(), time.Minute*15); e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	callback.Success(c, gin.H{
		"url":    mfaKey.URL(),
		"qrcode": "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(qrcodeBuffer.Bytes()),
	})
}

func MfaAddCheck(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	uid := tools.GetUserInfo(c).ID

	mfaSecret, e := redis.MfaEnable.Get(uid)
	if e != nil {
		if e == redis.Nil {
			callback.Error(c, callback.ErrMfaAddExpired)
			return
		}
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	valid, e := tools.VerifyMfa(f.Code, mfaSecret)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	} else if !valid {
		callback.Error(c, callback.ErrMfaCode)
		return
	}

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer userSrv.Rollback()

	if e = userSrv.SetMfaSecret(uid, mfaSecret); e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	if e = redis.MfaEnable.Del(uid); e != nil && e != redis.Nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	if e = userSrv.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Default(c)
}

func MfaDel(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer userSrv.Rollback()

	uid := tools.GetUserInfo(c).ID

	mfaSecret, e := userSrv.FindMfa(uid, daoUtil.LockForUpdate)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if mfaSecret == "" {
		callback.Error(c, callback.ErrMfaNotExist)
		return
	}

	valid, e := tools.VerifyMfa(f.Code, mfaSecret)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	} else if !valid {
		callback.Error(c, callback.ErrMfaCode)
		return
	}

	if e = userSrv.DelMfa(uid); e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	if e = userSrv.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Default(c)
}
