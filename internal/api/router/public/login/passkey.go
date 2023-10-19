package login

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/public/login"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerLoginPasskey(G *gin.RouterGroup) {
	G.Use(middlewares.EnableSession("passkey-login"))

	G.GET("/", controllers.BeginPasskeyLogin)
	G.POST("/", controllers.FinishPasskeyLogin)
}
