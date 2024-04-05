package token

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func Router(G *gin.RouterGroup) {
	refresh := G.Group("refresh")
	refresh.POST("/", controllers.RefreshToken)
	refresh.PATCH("/", controllers.ModifyRefreshPayload)
	refresh.DELETE("/", controllers.DestroyRefreshToken)

	access := G.Group("access", middlewares.RequireAccessToken)
	access.POST("verify", controllers.VerifyAccessToken)

	user := access.Group("user")
	user.POST("info", controllers.GetUserInfo)
}
