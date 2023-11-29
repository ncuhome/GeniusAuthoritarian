package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/jwt/jwtClaims"
	"net/url"
)

func GetUserInfo(c *gin.Context) *jwtClaims.UserToken {
	v, _ := c.Get("user")
	return v.(*jwtClaims.UserToken)
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

func GetAccessClaims(c *gin.Context) *jwtClaims.AccessToken {
	v, _ := c.Get("access")
	return v.(*jwtClaims.AccessToken)
}

func ShouldBindReused(c *gin.Context, obj any) error {
	bodyData, ok := c.Get(gin.BodyBytesKey)
	if ok {
		bodyBytes, ok := bodyData.([]byte)
		if ok {
			return binding.JSON.BindBody(bodyBytes, obj)
		}
	}
	return c.ShouldBindQuery(obj)
}
