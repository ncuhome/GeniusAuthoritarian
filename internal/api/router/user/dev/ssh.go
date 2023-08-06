package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/dev/ssh"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerSsh(G *gin.RouterGroup) {
	G.Use(middlewares.RequireMfa)

	G.GET("/", controllers.ShowSshKey)
}
