package main

import (
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cronAgent"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/department2BaseGroup"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/feishu"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/views"
	"github.com/ncuhome/GeniusAuthoritarian/internal/router"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/app"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/refreshToken"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev"
	sshDevServer "github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevSync"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cronInstance *cron.Cron

func init() {
	department2BaseGroup.Init(redis.NewSyncStat("init-base-groups"))

	cronInstance = cronAgent.New()
	feishu.InitSync(cronInstance)
	views.InitRenewAgent(cronInstance, redis.NewSyncStat("renew-views"))
	// 建议放在用户同步的时间之后
	sshDevServer.AddSshAccountCron(cronInstance, redis.NewSyncStat("dev-ssh"))

	cronInstance.Start()
}

func main() {
	log.Infoln("Sys Boost")

	httpSrv := &http.Server{
		Addr:    ":80",
		Handler: router.Engine(),
	}
	sshDevRpc := sshDev.NewRpc(global.Config.SshDev.Token)
	refreshTokenRpc := refreshToken.NewRpc()
	appRpc := app.NewRpc()

	go tools.RunHttpSrv(httpSrv)
	go tools.RunGrpcSrv(tools.MustTcpListen(":81"), sshDevRpc)
	go tools.RunGrpcSrv(tools.MustTcpListen(":82"), refreshTokenRpc)
	go tools.RunGrpcSrv(tools.MustTcpListen(":83"), appRpc)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	log.Infoln("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := httpSrv.Shutdown(ctx)
	if err != nil {
		log.Errorln("Http Server Shutdown:", err)
	}

	sshDevRpc.GracefulStop()
	refreshTokenRpc.GracefulStop()

	cronInstance.Stop()
}
