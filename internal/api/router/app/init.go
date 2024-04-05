package app

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/app/token"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/app/user"
)

func Router(G *gin.RouterGroup) {
	G.Use(middlewares.RequireAppSignature)

	keypair := G.Group("keypair")
	keypair.GET("server", controllers.ServerPublicKeys)
	keypair.POST("rpc", controllers.RpcClientCredential)

	user.Router(G.Group("user"))
	token.Router(G.Group("token"))
}
