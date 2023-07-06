//go:build !web && !fe

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router"
)

func Engine() *gin.Engine {
	E := gin.Default()

	router.Api(E.Group("api"))

	return E
}
