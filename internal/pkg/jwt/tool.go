package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func HeaderToken(c *gin.Context, Type string) (string, error) {
	tokenStr := c.GetHeader("Authorization")
	for _, token := range strings.Split(tokenStr, ", ") {
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
