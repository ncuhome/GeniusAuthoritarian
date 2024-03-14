package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerApp(G *gin.RouterGroup) {
	G.Use(middlewares.Secure())

	G.GET("/", controllers.AppInfo)
}
