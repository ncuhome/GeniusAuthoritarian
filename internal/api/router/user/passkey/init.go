package passkey

import "github.com/gin-gonic/gin"

func Router(G *gin.RouterGroup) {
	routerRegister(G.Group("register"))
}
