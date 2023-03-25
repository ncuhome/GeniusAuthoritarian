package cookie

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"time"
)

const refreshTokenKey = "token"

func SetRefreshToken(c *gin.Context, token string) {
	c.SetCookie(refreshTokenKey, token, int((time.Hour * 24 * 15).Seconds()), "", "", !global.DevMode, true)
}

func ClearRefreshToken(c *gin.Context) {
	c.SetCookie(refreshTokenKey, "", 0, "", "", !global.DevMode, true)
}
