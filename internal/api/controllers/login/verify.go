package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/signature"
	"time"
)

func doVerifyToken(c *gin.Context, token string) *jwt.LoginTokenClaims {
	claims, valid, e := jwt.ParseLoginToken(token)
	if e != nil || !valid {
		callback.Error(c, e, callback.ErrUnauthorized)
		return nil
	}

	loginRecordSrv, e := service.LoginRecord.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return nil
	}
	defer loginRecordSrv.Rollback()

	if e = loginRecordSrv.Add(claims.UID, claims.AppID, claims.IP); e != nil || loginRecordSrv.Commit().Error != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return nil
	}

	return claims
}

func VerifyToken(c *gin.Context) {
	var f struct {
		Token     string `json:"token" form:"token" binding:"required"`
		AppCode   string `json:"appCode" form:"appCode" binding:"required"`
		TimeStamp int64  `json:"timeStamp" form:"timeStamp" binding:"required"`
		Signature string `json:"signature" form:"signature" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	now := time.Now().Unix()
	diff := now - f.TimeStamp
	if diff > 300 {
		callback.Error(c, nil, callback.ErrUnauthorized)
		return
	}

	secret, e := service.App.GetSecretByAppCode(f.AppCode)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	ok, e := signature.CheckSignature(f.Signature, secret, struct {
		Token     string
		AppCode   string
		TimeStamp int64
	}{TimeStamp: f.TimeStamp, AppCode: f.AppCode, Token: f.Token})
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}
	if !ok {
		callback.Error(c, nil, callback.ErrUnauthorized)
		return
	}

	allowedGroups, e := service.AppGroup.GetGroupsByAppCode(f.AppCode)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callbackUrl, e := service.App.GetCallbackByAppCode(f.AppCode)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	claims := doVerifyToken(c, f.Token, callbackUrl, allowedGroups)
	if c.IsAborted() {
		return
	}

	callback.Success(c, response.VerifyTokenSuccess{
		Name:   claims.Name,
		Groups: claims.Groups,
	})
}

func Login(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	claims := doVerifyToken(c, f.Token, "", nil)
	if c.IsAborted() {
		return
	}

	token, e := jwt.GenerateUserToken(claims.UID)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callback.Success(c, gin.H{
		"token": token,
	})
}
