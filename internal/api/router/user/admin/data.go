package admin

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/ncuhome/GeniusAuthoritarian/internal/api/controllers/user/admin"
)

func routerData(G *gin.RouterGroup) {
	G.GET("login", controllers.LoginData)
}
