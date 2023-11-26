package login

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerToken(G *gin.RouterGroup) {
	G.Use(middlewares.RequireAppSignature)

	G.POST("/refresh", controllers.RefreshToken)
	G.POST("/access", controllers.VerifyAccessToken)
}
