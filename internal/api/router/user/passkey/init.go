package passkey

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/passkey"
)

func Router(G *gin.RouterGroup) {
	G.GET("/", controllers.ListPasskey)

	routerRegister(G.Group("register"))
}
