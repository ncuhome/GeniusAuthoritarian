//go:build fe

package router

import (
	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	E := gin.Default()
	serveFrontend(E)
	return E
}
