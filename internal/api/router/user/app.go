package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
)

func routerApp(G *gin.RouterGroup) {
	G.Use(middlewares.LimitGroup([]string{departments.UDev}))
	G.POST("/", controllers.ApplyApp)
}
