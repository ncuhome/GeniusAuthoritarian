package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
)

func routerMfa(G *gin.RouterGroup) {
	G.GET("/", controllers.MfaAdd)
	G.POST("/", controllers.MfaAddCheck)
	G.DELETE("/", controllers.MfaDel)
}
