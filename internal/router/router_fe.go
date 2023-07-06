//go:build fe

package router

import (
	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	E := gin.Default()
	serveFrontend(E)
	return E
}
