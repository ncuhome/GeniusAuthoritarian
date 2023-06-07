package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
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

	if claims.IP != c.ClientIP() {
		callback.Error(c, nil, callback.ErrNetContextChanged)
		return
	}

	valid, e := tools.VerifyMfa(f.Code, claims.Mfa)
	if e != nil {
		callback.Error(c, e, callback.ErrUnexpected)
		return
	} else if !valid {
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
