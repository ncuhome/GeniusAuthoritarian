package tools

import (
	"context"
	"github.com/Mmx233/tool"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Http *tool.Http

func init() {
	defaultTimeout := time.Second * 30

	Http = tool.NewHttpTool(tool.GenHttpClient(&tool.HttpClientOptions{
		Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
			Timeout: defaultTimeout,
		}),
		Timeout: defaultTimeout,
	}))
}

func SoftHttpSrv(E *gin.Engine) error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: E,
	}

	shutdown := make(chan bool)
	go func(srv *http.Server) {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-quit
		log.Infoln("Shutdown Server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		e := srv.Shutdown(ctx)
		if e != nil {
			log.Errorln("Server Shutdown:", e)
		}
		close(shutdown)
	}(srv)

	e := srv.ListenAndServe()
	if e == http.ErrServerClosed {
		<-shutdown
		return nil
	}
	return e
}
