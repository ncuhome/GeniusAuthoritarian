package router

import (
	gateway "github.com/Mmx233/Gateway/v2"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/tools"
	"github.com/ncuhome/GeniusAuthoritarian/web"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

func frontendRouterCheck(c *gin.Context) bool {
	return !strings.HasPrefix(c.Request.URL.Path, "/api")
}

var frontendLocalDevHandler = gateway.Proxy(&gateway.ApiConf{
	Addr:      "localhost:5173",
	Transport: tools.Http.Client.Transport,
	ErrorHandler: func(_ http.ResponseWriter, _ *http.Request, err error) {
		log.Warnf("调试页面请求转发失败: %v", err)
	},
	AllowRequest: func(c *gin.Context) bool {
		return frontendRouterCheck(c)
	},
})

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
