package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"strings"
)

func VerifyAccessToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	parts := strings.Split(f.Token, " ")
	if len(parts) != 2 {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	appCodeValue := strings.Split(parts[0], ":")
	if len(appCodeValue) != 2 {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}
	appCode := appCodeValue[1]
	accessToken := parts[1]

	claims, valid, err := jwt.ParseAccessToken(accessToken)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrUnauthorized)
		return
	} else if claims.AppCode != appCode {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	valid, err = redis.NewAccessJwt(claims.ID).Pair(claims.IssuedAt.Time)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
	} else if !valid {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	callback.Success(c, gin.H{
		"uid":   claims.ID,
		"name":  claims.Name,
		"group": claims.Groups,
	})
}
