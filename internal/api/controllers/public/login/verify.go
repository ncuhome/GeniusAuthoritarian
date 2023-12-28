package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/models/response"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"gorm.io/gorm"
	"io"
	"time"
)

// 验证 token 并添加登录记录
func doVerifyToken(c *gin.Context, tx *gorm.DB, token string, tokenValid time.Duration) (uint, *jwtClaims.LoginRedis) {
	claims, valid, err := jwt.ParseLoginToken(token)
	if err != nil || !valid {
		callback.Error(c, callback.ErrUnauthorized, err)
		return 0, nil
	}

	loginRecordSrv := service.LoginRecordSrv{DB: tx}
	lid, err := loginRecordSrv.Add(claims.UID, claims.AppID, claims.IP, claims.Useragent, tokenValid)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return 0, nil
	}
	return lid, claims
}

// CompleteLogin 第三方应用后端调用校验认证权威性
func CompleteLogin(c *gin.Context) {
	var f struct {
		Token    string `json:"token"  binding:"required"`
		ClientIp string `json:"clientIp"`

		GrantType string `json:"grantType" binding:"eq=|eq=refresh_token|eq=once"`

		Payload string `json:"payload" binding:"max=32"`
		// refreshToken 有效期，秒，最长 30 天，最短不在此处处理
		Valid int64 `json:"valid" binding:"min=0,max=2592000"`
	}
	cb, _ := c.Get(gin.BodyBytesKey)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(cb.([]byte)))
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	var tokenValid time.Duration
	var refreshTokenMode bool
	if f.GrantType == "refresh_token" {
		refreshTokenMode = true
		if f.Valid == 0 {
			f.Valid = 604800
		} else if f.Valid < 604800 {
			callback.ErrorWithTip(c, callback.ErrForm, "valid time too short, min 604800 seconds (7 days)")
			return
		}
		tokenValid = time.Duration(f.Valid) * time.Second
	}

	appSrv, err := service.App.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer appSrv.Rollback()

	loginRecordID, claims := doVerifyToken(c, appSrv.DB, f.Token, tokenValid)
	if c.IsAborted() {
		return
	}

	if f.ClientIp != "" && claims.IP != "" &&
		// 跳过 192.168.*.* ip 变动
		!tools.IsIntranet(f.ClientIp) && !tools.IsIntranet(claims.IP) &&
		claims.IP != f.ClientIp {
		callback.Error(c, callback.ErrNetContextChanged, "context="+claims.IP, "got="+f.ClientIp)
		return
	}

	res := &response.VerifyTokenSuccess{
		UserID:    claims.UID,
		Name:      claims.Name,
		Groups:    claims.Groups,
		AvatarUrl: claims.AvatarUrl,
	}

	if refreshTokenMode {
		appCode := tools.GetAppCode(c)
		var refreshClaims *jwtClaims.RefreshToken
		res.RefreshToken, refreshClaims, err = jwt.GenerateRefreshToken(claims.UID, uint64(loginRecordID), appCode, f.Payload, tokenValid)
		if err != nil {
			callback.Error(c, callback.ErrUnexpected, err)
			return
		}

		res.AccessToken, err = jwt.GenerateAccessToken(refreshClaims.ID, claims.UID, appCode, f.Payload)
	}

	if err = appSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, res)
}

// DashboardLogin 用户后台登录
func DashboardLogin(c *gin.Context) {
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

	var tokenValid = time.Hour * 24 * 3

	loginRecordID, claims := doVerifyToken(c, tx, f.Token, tokenValid)
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

	token, err := jwt.GenerateUserToken(claims.UID, uint64(loginRecordID), claims.Name, groups, tokenValid)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"token":  jwt.TokenWithType(jwt.User, token),
		"groups": groups,
	})
}
