package public

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/middlewares"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/app/token"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/app/user"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/public/login"
)

func Router(G *gin.RouterGroup) {
	routerApp(G.Group("app"))
	login.Router(G.Group("login"))

	// Deprecated, keep for compatibility
	token.Router(G.Group("token", middlewares.RequireAppSignature))
	user.Router(G.Group("user", middlewares.RequireAppSignature))
}
