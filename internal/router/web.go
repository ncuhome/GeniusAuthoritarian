package router

import (
	"github.com/Mmx233/tool"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/util"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	"github.com/ncuhome/GeniusAuthoritarian/web"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

func frontendRouterCheck(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/api")
}

func frontendLocalDevHandler(c *gin.Context) {
	if !frontendRouterCheck(c) {
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")

	if c.IsWebsocket() {
		connClient, e := tools.UpgradeDevWs(c)
		if e != nil {
			return
		}

		var headers = make(http.Header)
		headers.Set("Sec-Websocket-Protocol", c.Request.Header.Get("Sec-Websocket-Protocol"))
		for k, v := range c.Request.Header {
			if strings.HasPrefix(k, "Accept") {
				headers.Set(k, v[0])
			}
		}
		connServer, res, e := websocket.DefaultDialer.Dial("ws://127.0.0.1:5173"+c.Request.URL.Path, headers)
		if e != nil {
			connClient.Close()
			return
		}

		c.AbortWithStatus(res.StatusCode)

		go func() {
			defer connServer.Close()
			for {
				t, p, e := connClient.ReadMessage()
				if e != nil {
					return
				}
				_ = connServer.WriteMessage(t, p)
			}
		}()

		go func() {
			defer connClient.Close()
			for {
				t, p, e := connServer.ReadMessage()
				if e != nil {
					return
				}
				_ = connClient.WriteMessage(t, p)
			}
		}()
		return
	}

	query := make(map[string]interface{}, len(c.Request.URL.Query()))
	for k, v := range c.Request.URL.Query() {
		query[k] = v[0]
	}
	var header = make(map[string]interface{})
	for k, v := range c.Request.Header {
		if strings.HasPrefix(k, "Accept") {
			header[k] = v[0]
		}
	}
	res, e := util.Http.Request(c.Request.Method, &tool.DoHttpReq{
		Url:    "http://127.0.0.1:5173" + c.Request.URL.Path,
		Header: header,
		Query:  query,
		Body:   c.Request.Body,
	})
	if e != nil {
		c.AbortWithStatus(502)
		return
	}
	defer res.Body.Close()

	for k, v := range res.Header {
		if strings.HasPrefix(k, "Content-") {
			c.Header(k, v[0])
		}
	}
	c.Status(res.StatusCode)

	_, e = io.Copy(c.Writer, res.Body)
	if e != nil {
		c.AbortWithStatus(502)
		return
	}
	c.Abort()
}

func frontendProductionHandler() gin.HandlerFunc {
	fe, e := fs.Sub(web.FS, "dist")
	if e != nil {
		log.Fatalln(e)
	}
	file, err := fe.Open("index.html")
	if e != nil {
		log.Fatalln(e)
	}
	fileContentBytes, e := io.ReadAll(file)
	if err != nil {
		log.Fatalln(e)
	}
	_ = file.Close()
	index := string(fileContentBytes)

	fileServer := http.StripPrefix("/", http.FileServer(http.FS(fe)))

	return func(c *gin.Context) {
		if !frontendRouterCheck(c) {
			return
		}

		f, e := fe.Open(strings.TrimPrefix(c.Request.URL.Path, "/"))
		if e != nil {
			if _, ok := e.(*fs.PathError); ok {
				c.Header("Content-Type", "text/html")
				c.String(200, index)
				c.Abort()
				return
			}
			c.AbortWithStatus(500)
			return
		}
		_ = f.Close()
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func serveFrontend(E *gin.Engine) {
	//前端静态文件 serve
	if global.DevMode {
		E.Use(frontendLocalDevHandler)
	} else {
		E.Use(frontendProductionHandler())
	}
}
