package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router"
)

func CoreEngine() *gin.Engine {
	E := gin.Default()

	router.Api(E.Group("api"))

	serveFrontend(E)

	return E
}
