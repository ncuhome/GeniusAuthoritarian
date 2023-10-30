package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
)

func routerU2F(G *gin.RouterGroup) {
	G.GET("/", controllers.AvailableU2fMethod)
	G.POST("/:method", controllers.BeginU2F)

	G.PUT("prefer", controllers.UpdateU2fPrefer)
}
