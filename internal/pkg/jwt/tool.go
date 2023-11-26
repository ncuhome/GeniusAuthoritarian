package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func HeaderToken(c *gin.Context, Type string) (string, error) {
	tokens := c.Request.Header.Values("Authorization")
	for _, token := range tokens {
		if strings.HasPrefix(token, Type) {
			s := strings.Split(token, " ")
			if len(s) != 2 {
				return "", fmt.Errorf("token format error")
			}
			return s[1], nil
		}
	}
	return "", fmt.Errorf("%s token not found in header", Type)
}

func TokenWithType(Type, token string) string {
	return fmt.Sprintf("%s %s", Type, token)
}
