package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
)

func routerIdentity(G *gin.RouterGroup) {
	G.POST("sms", controllers.SendVerifySms)
}
