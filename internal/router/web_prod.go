//go:build !dev

package router

import (
	webServe "github.com/Mmx233/GinWebServe"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/web"
	log "github.com/sirupsen/logrus"
)

func frontendHandler() gin.HandlerFunc {
	fs, err := web.Fs()
	if err != nil {
		log.Fatalln(err)
	}

	handler, err := webServe.New(fs)
	if err != nil {
		log.Fatalln(err)
	}

	return handler
}
