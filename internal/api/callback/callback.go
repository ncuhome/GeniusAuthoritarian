package callback

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Msg struct {
	Code       uint8       `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Msg        string      `json:"msg,omitempty"`
	HttpStatus int         `json:"-"`
}

func Error(c *gin.Context, msg *Msg) {
	c.JSON(msg.HttpStatus, msg)
	c.AbortWithStatus(msg.HttpStatus)
}

func ErrorWithTip(c *gin.Context, msg *Msg, tip any) {
	msg.Msg = fmt.Sprint(tip)
	Error(c, msg)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, &Msg{
		Data: data,
	})
}

func Default(c *gin.Context) {
	Success(c, nil)
}
