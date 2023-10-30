package user

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user"
)

func routerProfile(G *gin.RouterGroup) {
	G.GET("/", controllers.ProfileData)
}
