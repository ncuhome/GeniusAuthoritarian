package app

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func Router(G *gin.RouterGroup) {
	G.Use(middlewares.RequireAppSignature)

	keypair := G.Group("keypair")
	keypair.GET("server", controllers.ServerPublicKeys)
}
