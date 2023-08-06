package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/dev/ssh"
)

func routerSsh(G *gin.RouterGroup) {
	G.GET("/", controllers.ShowSshKey)
}
