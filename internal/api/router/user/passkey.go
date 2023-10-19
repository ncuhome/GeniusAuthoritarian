package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerPasskey(G *gin.RouterGroup) {
	G.Use(middlewares.EnableSession("passkey-user"))
}
