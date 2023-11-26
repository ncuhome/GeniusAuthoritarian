package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt"
	"net/url"
)

func GetUserInfo(c *gin.Context) *jwt.UserToken {
	v, _ := c.Get("user")
	return v.(*jwt.UserToken)
}

func GenCallback(callback, token string) (string, error) {
	callbackUrl, err := url.Parse(callback)
	if err != nil {
		return "", err
	}
	callbackQuery := callbackUrl.Query()
	callbackQuery.Set("token", token)
	callbackUrl.RawQuery = callbackQuery.Encode()
	return callbackUrl.String(), nil
}

func GetAppCode(c *gin.Context) string {
	v, _ := c.Get("appCode")
	return v.(string)
}

func GetAccessClaims(c *gin.Context) *jwt.AccessToken {
	v, _ := c.Get("access")
	return v.(*jwt.AccessToken)
}
