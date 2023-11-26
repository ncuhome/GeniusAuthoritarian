package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func RefreshToken(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	claims, valid, err := jwt.ParseRefreshToken(f.Token)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	valid, err = redis.NewRefreshJwt(claims.ID).Pair(claims.IssuedAt.Time)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	accessToken, accessClaims, err := jwt.GenerateAccessToken(claims.ID, claims.Name, claims.AppCode, claims.Groups)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	err = redis.NewAccessJwt(claims.ID).Set(accessClaims.IssuedAt.Time, accessClaims.ExpiresAt.Sub(accessClaims.IssuedAt.Time))
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"access_token": accessToken,
	})
}
