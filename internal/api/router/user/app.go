package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
)

func routerApp(G *gin.RouterGroup) {
	G.GET("/accessible", controllers.ListAccessibleApp)
}
