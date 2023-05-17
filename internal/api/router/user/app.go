package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
)

func routerApp(G *gin.RouterGroup) {
	G.Use(middlewares.LimitGroup([]string{departments.UDev}))
	G.GET("/", controllers.ListOwnedApp)
	G.GET("/accessible", controllers.ListAccessibleApp)
	G.POST("/", controllers.ApplyApp)
	G.PUT("/", controllers.ModifyApp)
	G.DELETE("/", controllers.DeleteApp)
}
