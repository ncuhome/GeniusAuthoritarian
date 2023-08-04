package dev

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/app"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/", controllers.ListOwnedApp)
	G.POST("/", controllers.ApplyApp)
	G.PUT("/", controllers.ModifyApp)
	G.PUT("/linkOff", controllers.UpdateLinkState)
	G.DELETE("/", controllers.DeleteApp)
}
