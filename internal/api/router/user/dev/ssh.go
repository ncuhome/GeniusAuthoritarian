package dev

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/ssh"
)

func routerSsh(G *gin.RouterGroup) {
	G.GET("/", controllers.ShowSshKey)
}
