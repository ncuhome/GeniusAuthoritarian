package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func Router(G *gin.RouterGroup) {
	G.Use(middlewares.UserAuth)

	routerProfile(G.Group("profile"))
	routerApp(G.Group("app"))
	routerGroups(G.Group("group"))
}
