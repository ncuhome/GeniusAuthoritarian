package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func UserAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		callback.Error(c, nil, callback.ErrUnauthorized)
		return
	}

	claims, valid, e := jwt.ParseUserToken(token)
	if e != nil || !valid {
		callback.Error(c, e, callback.ErrUnauthorized)
		return
	}

	c.Set("Uid", claims.ID)
}
