package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
)

func routerGroups(G *gin.RouterGroup) {
	G.GET("/list", controllers.ListGroups)
}
