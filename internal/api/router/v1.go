package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/router/user"
)

func ApiV1(G *gin.RouterGroup) {
	user.Router(G.Group("user"))
}
