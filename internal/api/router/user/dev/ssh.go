package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/dev/ssh"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
)

func routerSsh(G *gin.RouterGroup) {
	G.Use(middlewares.RequireU2F)

	G.GET("/", controllers.ShowSshKey)
	G.PUT("/", controllers.ResetSshKeyPair)
	G.POST("killall", controllers.KillAllProcess)
}
