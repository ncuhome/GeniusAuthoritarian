package tools

import "github.com/gin-gonic/gin"

func GetUID(c *gin.Context) uint {
	v, _ := c.Get("Uid")
	return v.(uint)
}
