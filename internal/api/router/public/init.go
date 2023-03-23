package public

import "github.com/gin-gonic/gin"

func Router(G *gin.RouterGroup) {
	routerLogin(G.Group("login"))
}
