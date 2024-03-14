package router

import (
	"github.com/gin-gonic/gin"
)

func Api(G *gin.RouterGroup) {
	ApiV1(G.Group("v1"))
}
