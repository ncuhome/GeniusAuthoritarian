package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/app"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/accessible", controllers.ListAccessibleApp)
	G.GET("/landing", controllers.LandingApp)
}
