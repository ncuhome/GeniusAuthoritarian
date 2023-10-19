package public

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public/login"
)

func Router(G *gin.RouterGroup) {
	routerApp(G.Group("app"))
	login.Router(G.Group("login"))
}
