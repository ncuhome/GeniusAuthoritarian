package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"gorm.io/gorm"
	"time"
)

// 验证 token 并添加登录记录
func doVerifyToken(c *gin.Context, tx *gorm.DB, token string) *jwt.LoginRedisClaims {
	claims, valid, err := jwt.ParseLoginToken(token)
	if err != nil || !valid {
		callback.Error(c, callback.ErrUnauthorized, err)
		return nil
	}

	loginRecordSrv := service.LoginRecordSrv{DB: tx}
	if err = loginRecordSrv.Add(claims.UID, claims.AppID, claims.IP); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return nil
	}

	return claims
}

// VerifyToken 第三方应用后端调用校验认证权威性
func VerifyToken(c *gin.Context) {
	var f struct {
		Token    string `json:"token" form:"token" binding:"required"`
		ClientIp string `json:"clientIp" form:"clientIp"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	appSrv, err := service.App.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer appSrv.Rollback()

	claims := doVerifyToken(c, appSrv.DB, f.Token)
	if c.IsAborted() {
		return
	}

	if f.ClientIp != "" && claims.IP != f.ClientIp {
		callback.Error(c, callback.ErrNetContextChanged, "context="+claims.IP, "got="+f.ClientIp)
		return
	}

	if err = appSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
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
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
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
	groups, err := userGroupSrv.GetNamesForUser(claims.UID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	var tokenValid = time.Hour * 24 * 3

	token, userTokenClaims, err := jwt.GenerateUserToken(claims.UID, claims.Name, groups, tokenValid)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	if err = redis.NewUserJwt(claims.UID).Set(userTokenClaims.IssuedAt.Time, tokenValid); err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"token":  token,
		"groups": groups,
	})
}
