package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerToken(G *gin.RouterGroup) {
	G.Use(middlewares.RequireAppSignature)

	G.POST("refresh", controllers.RefreshToken)

	access := G.Group("access", middlewares.RequireAccessToken)
	access.POST("verify", controllers.VerifyAccessToken)
}
