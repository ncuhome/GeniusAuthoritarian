package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func UserAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	claims, valid, err := jwt.ParseUserToken(token)
	if err != nil || !valid {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	valid, err = redis.NewUserJwt(claims.ID).Pair(token)
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !valid {
		callback.Error(c, callback.ErrUnauthorized)
		return
	}

	c.Set("user", claims)
}
