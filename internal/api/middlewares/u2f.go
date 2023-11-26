package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func RequireU2F(c *gin.Context) {
	token, err := jwt.HeaderToken(c, jwt.U2F)
	if err != nil {
		callback.Error(c, callback.ErrUnauthorized, err)
		return
	}

	ok, err := jwt.ParseU2fToken(token, c.ClientIP())
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	} else if !ok {
		callback.Error(c, callback.ErrU2fTokenExpired)
		return
	}
}
