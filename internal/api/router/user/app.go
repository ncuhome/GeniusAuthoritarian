package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
)

func routerApp(G *gin.RouterGroup) {
	G.Use(middlewares.LimitGroup([]string{departments.UDev}))
}
