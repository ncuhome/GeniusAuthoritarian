package public

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/app"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/", controllers.AppInfo)
}
