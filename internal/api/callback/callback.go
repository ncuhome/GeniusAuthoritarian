package callback

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Msg struct {
	Code       uint8       `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Msg        string      `json:"msg,omitempty"`
	HttpStatus int         `json:"-"`
}

func Error(c *gin.Context, msg *Msg, errors ...error) {
	log.Debugln(msg.Msg, errors)
	c.JSON(msg.HttpStatus, msg)
	c.AbortWithStatus(msg.HttpStatus)
}

func ErrorWithTip(c *gin.Context, msg *Msg, tip any, errors ...error) {
	tipMsg := *msg
	tipMsg.Msg = fmt.Sprint(tip)
	Error(c, &tipMsg, errors...)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, &Msg{
		Data: data,
	})
}

func Default(c *gin.Context) {
	Success(c, nil)
}
