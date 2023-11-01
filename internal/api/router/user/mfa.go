package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerMfa(G *gin.RouterGroup) {
	G.GET("/", middlewares.RequireU2F, controllers.MfaAdd)
	G.POST("/", controllers.MfaAddCheck)
	G.DELETE("/", middlewares.RequireU2F, controllers.MfaDel)
}
