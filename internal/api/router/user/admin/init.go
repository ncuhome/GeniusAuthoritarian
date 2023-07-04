package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
)

func Router(G *gin.RouterGroup) {
	G.Use(middlewares.LimitGroup(departments.UCe))
}
