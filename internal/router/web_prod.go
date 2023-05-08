//go:build !dev

package router

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/web"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

func calcEtag(d []byte) string {
	hash := md5.New()
	hash.Write(d)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

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
	indexEtag := calcEtag(fileContentBytes)

	fileServer := http.StripPrefix("/", http.FileServer(http.FS(fe)))

	return func(c *gin.Context) {
		if !frontendRouterCheck(c) {
			return
		}

		f, e := fe.Open(strings.TrimPrefix(c.Request.URL.Path, "/"))
		if e != nil {
			if _, ok := e.(*fs.PathError); ok {
				if c.GetHeader("If-None-Match") == indexEtag {
					c.AbortWithStatus(304)
					return
				}
				c.Header("Content-Type", "text/html")
				c.Header("Cache-Control", "no-cache")
				c.Header("Etag", indexEtag)
				c.String(200, index)
				c.Abort()
				return
			} else {
				log.Errorln("加载 embed fs 失败:", e)
				c.AbortWithStatus(500)
				return
			}
		}
		_ = f.Close()

		c.Header("Cache-Control", "public, max-age=2592000, immutable")
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
