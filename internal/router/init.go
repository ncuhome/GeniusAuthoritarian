package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router"
)

var E *gin.Engine

func init() {
	E = gin.Default()

	router.Api(E.Group("/api"))
}
