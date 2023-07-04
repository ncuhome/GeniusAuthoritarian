//go:build !dev

package global

import (
	"github.com/Mmx233/EnvConfig"
	"github.com/gin-gonic/gin"
)

func initConfig() {
	EnvConfig.Load("", &Config)
	gin.SetMode(gin.ReleaseMode)
}
