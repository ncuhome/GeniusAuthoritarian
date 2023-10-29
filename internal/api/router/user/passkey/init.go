package passkey

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/passkey"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func Router(G *gin.RouterGroup) {
	G.GET("/", controllers.ListPasskey)
	G.DELETE("/", middlewares.RequireU2F, controllers.DeletePasskey)
	G.PATCH("/", controllers.RenamePasskey)

	routerRegister(G.Group("register"))
}
