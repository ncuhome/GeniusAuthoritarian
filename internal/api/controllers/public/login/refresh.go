package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func RefreshToken(c *gin.Context) {
	var f struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		callback.Error(c, callback.ErrForm, err)
		return
	}

	claims, valid, err := jwt.ParseRefreshToken(f.RefreshToken)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	accessToken, err := jwt.GenerateAccessToken(claims.ID, claims.Name, claims.AppCode, claims.Groups)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	callback.Success(c, gin.H{
		"access_token": accessToken,
	})
}
