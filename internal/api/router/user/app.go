package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/accessible", controllers.ListAccessibleApp)
	G.GET("/landing", controllers.LandingApp)

	routerOwnedApp(G.Group("owned"))
}

func routerOwnedApp(G *gin.RouterGroup) {
	G.Use(middlewares.LimitGroup(departments.UDev))
	G.GET("/", controllers.ListOwnedApp)
	G.POST("/", controllers.ApplyApp)
	G.PUT("/", controllers.ModifyApp)
	G.PUT("/linkOff", controllers.UpdateLinkState)
	G.DELETE("/", controllers.DeleteApp)
}
