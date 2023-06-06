package controllers

import (
	"bytes"
	"encoding/base64"
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/pquerna/otp/totp"
	"image/jpeg"
)

func AddMfa(c *gin.Context) {
	uid := tools.GetUserInfo(c).ID

	userSrv, e := service.User.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	defer userSrv.Rollback()

	exist, e := userSrv.MfaExist(uid, daoUtil.LockForUpdate)
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

	if e = userSrv.SetMfaSecret(uid, mfaKey.Secret()); e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if e = userSrv.Commit().Error; e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	callback.Success(c, gin.H{
		"secret": mfaKey.Secret(),
		"url":    mfaKey.URL(),
		"qr":     "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(qrBuffer.Bytes()),
	})
}
