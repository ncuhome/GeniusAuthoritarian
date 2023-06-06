package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/pquerna/otp/totp"
)

func VerifyMfa(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
		Code  string `json:"code" form:"code" binding:"required,len=6,numeric"`
	}
	if e := c.ShouldBind(&f); e != nil {
		callback.Error(c, e, callback.ErrForm)
		return
	}

	claims, e := jwt.ParseMfaToken(f.Token)
	if e != nil {
		callback.Error(c, nil, callback.ErrUnauthorized)
		return
	}

	if !totp.Validate(f.Code, claims.Mfa) {
		callback.Error(c, nil, callback.ErrMfaCode)
		return
	}

	token, e := jwt.GenerateLoginToken(claims.UID, claims.AppID, claims.Name, claims.IP, claims.Groups)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callbackUrl, e := tools.GenCallback(claims.AppCallback, token)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	}

	callback.Success(c, gin.H{
		"token":    token,
		"callback": callbackUrl,
	})
}
