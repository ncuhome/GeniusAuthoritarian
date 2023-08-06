package dev

import (
	"github.com/gin-gonic/gin"
	controllers2 "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/dev/app"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/", controllers2.ListOwnedApp)
	G.POST("/", controllers2.ApplyApp)
	G.PUT("/", controllers2.ModifyApp)
	G.PUT("/linkOff", controllers2.UpdateLinkState)
	G.DELETE("/", controllers2.DeleteApp)
}
