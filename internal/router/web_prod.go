//go:build !dev

package router

import (
	webServe "github.com/Mmx233/GinWebServe"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/web"
	log "github.com/sirupsen/logrus"
)

func frontendHandler() gin.HandlerFunc {
	fs, e := web.Fs()
	if e != nil {
		log.Fatalln(e)
	}

	handler, e := webServe.New(fs)
	if e != nil {
		log.Fatalln(e)
	}

	return handler
}
