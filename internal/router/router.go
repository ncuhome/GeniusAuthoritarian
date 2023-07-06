//go:build !fe

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
)

func init() {
	if !global.Config.TraceMode {
		gin.SetMode(gin.ReleaseMode)
	}
}
