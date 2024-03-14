package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerUser(G *gin.RouterGroup) {
	G.Use(middlewares.RequireAppSignature)

	G.GET("info", controllers.GetUserPublicInfo)
}
