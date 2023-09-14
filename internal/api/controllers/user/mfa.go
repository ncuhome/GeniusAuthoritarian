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
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	uid := tools.GetUserInfo(c).ID

	mfaSecret, err := redis.MfaEnable.Get(uid)
	if err != nil {
		if err == redis.Nil {
			callback.Error(c, callback.ErrMfaAddExpired)
			return
		}
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	valid, err := tools.VerifyMfa(f.Code, mfaSecret)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrMfaCode)
		return
	}

	userSrv, err := service.User.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer userSrv.Rollback()

	if err = userSrv.SetMfaSecret(uid, mfaSecret); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = redis.MfaEnable.Del(uid); err != nil && err != redis.Nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	if err = userSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}

func MfaDel(c *gin.Context) {
	var f struct {
		Code string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	userSrv, err := service.User.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer userSrv.Rollback()

	uid := tools.GetUserInfo(c).ID

	mfaSecret, err := userSrv.FirstMfa(uid, daoUtil.LockForUpdate)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if mfaSecret == "" {
		callback.Error(c, callback.ErrMfaNotExist)
		return
	}

	valid, err := tools.VerifyMfa(f.Code, mfaSecret)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrMfaCode)
		return
	}

	if err = userSrv.DelMfa(uid); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = userSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}
