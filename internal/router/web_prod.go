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
	"sync"
)

func calcEtag(f fs.File) (string, error) {
	defer f.Close()

	d, e := io.ReadAll(f)
	if e != nil {
		return "", e
	}

	hash := md5.New()
	hash.Write(d)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func frontendHandler() gin.HandlerFunc {
	fe, e := fs.Sub(web.FS, "dist")
	if e != nil {
		log.Fatalln(e)
	}
	fileServer := http.StripPrefix("/", http.FileServer(http.FS(fe)))

	const IndexPath = "/index.html"
	var etags = &sync.Map{}

	return func(c *gin.Context) {
		if !frontendRouterCheck(c) {
			return
		}

	checkFile:
		f, e := fe.Open(strings.TrimPrefix(c.Request.URL.Path, "/"))
		if e != nil {
			if _, ok := e.(*fs.PathError); ok && c.Request.URL.Path != IndexPath {
				c.Request.URL.Path = IndexPath
				goto checkFile
			} else {
				log.Errorln("加载 embed fs 失败:", e)
				c.AbortWithStatus(500)
				return
			}
		}

		etag, ok := etags.Load(c.Request.URL.Path)
		if ok {
			_ = f.Close()
		} else {
			etag, e = calcEtag(f)
			if e != nil {
				log.Errorln("计算 etag 失败:", e)
				c.AbortWithStatus(500)
				return
			}
			etags.Store(c.Request.URL.Path, etag)
		}

		c.Header("Etag", etag.(string))

		if c.GetHeader("If-none-match") == etag.(string) {
			c.AbortWithStatus(304)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
