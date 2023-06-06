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
	"github.com/pquerna/otp/totp"
	"image/jpeg"
	"time"
)

func MfaAdd(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	defer userSrv.Rollback()

	exist, e := userSrv.MfaExist(uid, daoUtil.LockForShare)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if exist {
		callback.Error(c, e, callback.ErrMfaAlreadyExist)
		return
	}

	mfaKey, e := totp.Generate(totp.GenerateOpts{})
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	qr, e := mfaKey.Image(300, 300)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	var qrBuffer = bytes.Buffer{}
	if e = jpeg.Encode(&qrBuffer, qr, &jpeg.Options{
		Quality: 100,
	}); e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	if e = redis.MfaEnable.Set(uid, mfaKey.Secret(), time.Minute*15); e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callback.Success(c, gin.H{
		"secret": mfaKey.Secret(),
		"url":    mfaKey.URL(),
		"qr":     "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(qrBuffer.Bytes()),
	})
}

func MfaCheck(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	uid := tools.GetUserInfo(c).ID

	mfaSecret, e := redis.MfaEnable.Get(uid)
	if e != nil {
		if e == redis.Nil {
			callback.Error(c, nil, callback.ErrMfaAddExpired)
			return
		}
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if !totp.Validate(f.Code, mfaSecret) {
		callback.Error(c, nil, callback.ErrMfaCode)
		return
	}

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	defer userSrv.Rollback()

	if e = userSrv.SetMfaSecret(uid, mfaSecret); e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if e = redis.MfaEnable.Del(uid); e != nil && e != redis.Nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	if e = userSrv.Commit().Error; e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callback.Default(c)
}

func MfaDel(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	defer userSrv.Rollback()

	uid := tools.GetUserInfo(c).ID

	mfaSecret, e := userSrv.FindMfa(uid, daoUtil.LockForUpdate)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	} else if mfaSecret == "" {
		callback.Error(c, nil, callback.ErrMfaNotExist)
		return
	}

	if !totp.Validate(f.Code, mfaSecret) {
		callback.Error(c, nil, callback.ErrMfaCode)
		return
	}

	if e = userSrv.DelMfa(uid); e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if e = userSrv.Commit().Error; e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callback.Default(c)
}
