package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
)

func GetUserInfo(c *gin.Context) *jwt.UserToken {
	v, _ := c.Get("user")
	return v.(*jwt.UserToken)
}
