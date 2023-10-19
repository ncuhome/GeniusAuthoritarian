package login

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
)

func routerLoginPasskey(G *gin.RouterGroup) {
	G.GET("/", controllers.BeginPasskeyLogin)
	G.POST("/", controllers.FinishPasskeyLogin)
}
