package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public"
)

func CoreEngine() *gin.Engine {
	E := gin.Default()

	router.Api(E.Group("api"))

	serveFrontend(E)

	return E
}

func GateEngine() *gin.Engine {
	E := gin.Default()

	public.Router(E.Group("/api/v1/public"))

	return E
}
