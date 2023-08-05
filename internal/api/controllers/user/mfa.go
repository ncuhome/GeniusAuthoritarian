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
	var f struct {
		Code string `json:"code" form:"code" binding:"required"` // 身份校验码（短信）
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	uid := tools.GetUserInfo(c).ID

	ok, err := redis.UserIdentityCode.VerifyAndDestroy(uid, f.Code)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !ok {
		callback.Error(c, callback.ErrIdentityCodeNotCorrect)
		return
	}

	userSrv, err := service.User.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer userSrv.Rollback()

	exist, err := userSrv.MfaExist(uid, daoUtil.LockForShare)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if exist {
		callback.Error(c, callback.ErrMfaAlreadyExist, err)
		return
	}

	mfaKey, err := tools.NewMfa(uid)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	qrcodeImage, err := mfaKey.Image(300, 300)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	var qrcodeBuffer = bytes.Buffer{}
	if err = jpeg.Encode(&qrcodeBuffer, qrcodeImage, &jpeg.Options{
		Quality: 100,
	}); err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	if err = redis.MfaEnable.Set(uid, mfaKey.Secret(), time.Minute*15); err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
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

	mfaSecret, e := userSrv.FirstMfa(uid, daoUtil.LockForUpdate)
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
