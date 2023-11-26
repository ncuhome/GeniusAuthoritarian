package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
)

func VerifyMfa(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
		Code  string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	claims, err := jwt.ParseMfaToken(f.Token)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	if claims.IP != c.ClientIP() {
		callback.Error(c, callback.ErrNetContextChanged)
		return
	}

	valid, err := tools.VerifyMfa(f.Code, claims.Mfa)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrMfaCode)
		return
	}

	token, err := jwt.GenerateLoginToken(claims.LoginRedisClaims)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	app, err := service.App.FirstAppByID(claims.AppID)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	accessToken, accessClaims, err := jwt.GenerateAccessToken(claims.UID, claims.Name, app.AppCode, claims.Groups)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	refreshToken, refreshClaims, err := jwt.GenerateRefreshToken(claims.UID, claims.Name, app.AppCode, claims.Groups)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	err = redis.NewAccessJwt(claims.UID).Set(accessClaims.IssuedAt.Time, accessClaims.ExpiresAt.Time.Sub(accessClaims.IssuedAt.Time))
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	err = redis.NewRefreshJwt(claims.UID).Set(refreshClaims.IssuedAt.Time, refreshClaims.ExpiresAt.Time.Sub(refreshClaims.IssuedAt.Time))
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callbackUrl, err := tools.GenCallback(claims.AppCallback, token)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"token":        token,
		"callback":     callbackUrl,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
