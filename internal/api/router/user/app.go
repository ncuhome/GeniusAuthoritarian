package user

import (
	"github.com/gin-gonic/gin"
	controllers2 "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/dev/app"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/accessible", controllers2.ListAccessibleApp)
	G.GET("/landing", controllers2.LandingApp)
}
