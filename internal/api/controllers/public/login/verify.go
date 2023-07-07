package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarianClient/pkg/signature"
	"gorm.io/gorm"
	"time"
)

// 验证 token 并添加登录记录
func doVerifyToken(c *gin.Context, tx *gorm.DB, token string) *jwt.LoginTokenClaims {
	claims, valid, e := jwt.ParseLoginToken(token)
	if e != nil || !valid {
		callback.Error(c, callback.ErrUnauthorized, e)
		return nil
	}

	loginRecordSrv := service.LoginRecordSrv{DB: tx}
	if e = loginRecordSrv.Add(claims.UID, claims.AppID, claims.IP); e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return nil
	}

	return claims
}

// VerifyToken 第三方应用后端调用校验认证权威性
func VerifyToken(c *gin.Context) {
	var f struct {
		Token     string `json:"token" form:"token" binding:"required"`
		AppCode   string `json:"appCode" form:"appCode" binding:"required"`
		TimeStamp int64  `json:"timeStamp" form:"timeStamp" binding:"required"`
		Signature string `json:"signature" form:"signature" binding:"required"`

		ClientIp string `json:"clientIp" form:"clientIp"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	if time.Now().Sub(time.Unix(f.TimeStamp, 0)) > time.Minute*5 {
		callback.Error(c, callback.ErrSignatureExpired)
		return
	}

	appSrv, e := service.App.Begin()
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}
	defer appSrv.Rollback()

	claims := doVerifyToken(c, appSrv.DB, f.Token)
	if c.IsAborted() {
		return
	}

	if f.ClientIp != "" && claims.IP != f.ClientIp {
		callback.Error(c, callback.ErrNetContextChanged)
		return
	}

	appCode, appSecret, e := appSrv.FirstAppKeyPair(claims.AppID)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	} else if f.AppCode != appCode {
		callback.Error(c, callback.ErrOperationIllegal)
		return
	}

	if !signature.Check(f.Signature, &signature.VerifyClaims{
		Token:     f.Token,
		AppCode:   f.AppCode,
		TimeStamp: f.TimeStamp,
		AppSecret: appSecret,
	}) {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	if e = appSrv.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	callback.Success(c, response.VerifyTokenSuccess{
		UserID:    claims.UID,
		Name:      claims.Name,
		Groups:    claims.Groups,
		AvatarUrl: claims.AvatarUrl,
	})
}

// Login 用户后台登录
func Login(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, callback.ErrForm, e)
		return
	}

	tx := dao.DB.Begin()
	if tx.Error != nil {
		callback.Error(c, callback.ErrDBOperation, tx.Error)
		return
	}
	defer tx.Rollback()

	claims := doVerifyToken(c, tx, f.Token)
	if c.IsAborted() {
		return
	} else if claims.AppID != 0 {
		callback.Error(c, callback.ErrOperationIllegal)
		return
	} else if claims.IP != c.ClientIP() {
		callback.Error(c, callback.ErrNetContextChanged)
		return
	}

	userGroupSrv := service.UserGroupsSrv{DB: tx}
	groups, e := userGroupSrv.GetForUser(claims.UID)
	if e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	if e = tx.Commit().Error; e != nil {
		callback.Error(c, callback.ErrDBOperation, e)
		return
	}

	token, e := jwt.GenerateUserToken(claims.UID, claims.Name, groups)
	if e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	if e = redis.UserJwt.Set(claims.UID, token, time.Hour*24*15); e != nil {
		callback.Error(c, callback.ErrUnexpected, e)
		return
	}

	callback.Success(c, gin.H{
		"token":  token,
		"groups": groups,
	})
}
