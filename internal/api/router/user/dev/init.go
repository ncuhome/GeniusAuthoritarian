package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
)

func Router(G *gin.RouterGroup) {
	G.Use(middlewares.LimitGroup(departments.UDev))

	routerApp(G.Group("app"))
	routerSsh(G.Group("ssh"))
}
