package login

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerLoginPasskey(G *gin.RouterGroup) {
	G.Use(middlewares.EnableSession("passkey-login"))
}
