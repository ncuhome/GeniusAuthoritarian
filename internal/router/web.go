//go:build web

package router

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func frontendRouterCheck(c *gin.Context) bool {
	return !strings.HasPrefix(c.Request.URL.Path, "/api")
}

func serveFrontend(E *gin.Engine) {
	//前端静态文件 serve
	E.Use(frontendHandler())
}
