package passkey

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/passkey"
)

func routerRegister(G *gin.RouterGroup) {
	G.GET("/", controllers.BeginPasskeyRegistration)
	G.POST("/", controllers.FinishPasskeyRegistration)
}
