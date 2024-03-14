package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerU2F(G *gin.RouterGroup) {
	G.GET("/", controllers.AvailableU2fMethod)
	G.POST("/:method", middlewares.Secure(), controllers.BeginU2F)

	G.PUT("prefer", controllers.UpdateU2fPrefer)
}
