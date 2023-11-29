package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func UserAuth(c *gin.Context) {
	token, err := jwt.HeaderToken(c, jwt.User)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	claims, valid, err := jwt.ParseUserToken(token)
	if err != nil || !valid {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	c.Set("user", claims)
}
