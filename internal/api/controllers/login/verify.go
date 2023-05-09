package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/signature"
	"gorm.io/gorm"
	"time"
)

func doVerifyToken(c *gin.Context, tx *gorm.DB, token string) *jwt.LoginTokenClaims {
	claims, valid, e := jwt.ParseLoginToken(token)
	if e != nil || !valid {
		callback.Error(c, e, callback.ErrUnauthorized)
		return nil
	}

	loginRecordSrv := service.LoginRecordSrv{DB: tx}
	if e = loginRecordSrv.Add(claims.UID, claims.AppID, claims.IP); e != nil {
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

	if time.Now().Sub(time.Unix(f.TimeStamp, 0)) > time.Minute*5 {
		callback.Error(c, nil, callback.ErrSignatureExpired)
		return
	}

	appSrv, e := service.App.Begin()
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}
	defer appSrv.Rollback()

	claims := doVerifyToken(c, appSrv.DB, f.Token)
	if c.IsAborted() {
		return
	}

	secret, e := appSrv.FirstAppSecret(claims.AppID)
	if e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
		return
	}

	if !signature.Check(f.Signature, &signature.VerifyClaims{
		Token:     f.Token,
		AppCode:   f.AppCode,
		TimeStamp: f.TimeStamp,
		AppSecret: secret,
	}) {
		callback.Error(c, nil, callback.ErrUnauthorized)
		return
	}

	if e = appSrv.Commit().Error; e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
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

	tx := dao.DB.Begin()
	if tx.Error != nil {
		callback.Error(c, tx.Error, callback.ErrDBOperation)
		return
	}
	defer tx.Rollback()

	claims := doVerifyToken(c, tx, f.Token)
	if c.IsAborted() {
		return
	}

	if e := tx.Commit().Error; e != nil {
		callback.Error(c, e, callback.ErrDBOperation)
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
