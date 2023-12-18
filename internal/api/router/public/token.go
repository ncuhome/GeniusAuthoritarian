package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerToken(G *gin.RouterGroup) {
	G.Use(middlewares.RequireAppSignature)

	refresh := G.Group("refresh")
	refresh.POST("/", controllers.RefreshToken)
	refresh.PATCH("/", controllers.ModifyRefreshPayload)
	refresh.DELETE("/", controllers.DestroyRefreshToken)

	access := G.Group("access", middlewares.RequireAccessToken)
	access.POST("verify", controllers.VerifyAccessToken)

	user := access.Group("user")
	user.POST("info", controllers.GetUserInfo)
}
