package passkey

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/passkey"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerRegister(G *gin.RouterGroup) {
	G.GET("/", middlewares.RequireU2F, controllers.BeginPasskeyRegistration)
	G.POST("/", controllers.FinishPasskeyRegistration)
}
