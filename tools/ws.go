package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var DevUpper = websocket.Upgrader{
	HandshakeTimeout: time.Minute * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpgradeDevWs(c *gin.Context) (*websocket.Conn, error) {
	return DevUpper.Upgrade(c.Writer, c.Request, map[string][]string{
		"Sec-WebSocket-Protocol": {c.GetHeader("Sec-WebSocket-Protocol")},
	})
}
