package tools

import (
	"errors"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
	"time"
)

var Http *tool.Http

func init() {
	defaultTimeout := time.Second * 30

	Http = tool.NewHttpTool(tool.GenHttpClient(&tool.HttpClientOptions{
		Transport: &http.Transport{
			TLSHandshakeTimeout: defaultTimeout,
			IdleConnTimeout:     time.Hour * 3,
		},
		Timeout: defaultTimeout,
	}))
}

func MustTcpListen(addr string) net.Listener {
	tcpListen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("tcp listen addr %s failed: %v", addr, err)
	}
	return tcpListen
}

func RunHttpSrv(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		log.Fatalln("run api server failed:", err)
	}
}

func RunGrpcSrv(tcpListen net.Listener, srv *grpc.Server) {
	err := srv.Serve(tcpListen)
	if err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return
		}
		log.Fatalln("run rpc server failed:", err)
	}
}

// IsIntranet 是否是同一局域网
func IsIntranet(ip string) bool {
	return strings.HasPrefix(ip, "192.168.")
}
