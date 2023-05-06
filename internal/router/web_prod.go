//go:build !dev

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/web"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

func frontendHandler() gin.HandlerFunc {
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
