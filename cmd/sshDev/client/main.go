package main

import (
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/client"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/linux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func init() {
	go linux.DaemonSshd()
}

func main() {
	log.Infoln("Sys Boost")

	// 读取配置
	conf := client.ReadConfig()

	// 连接 grpc
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(
		conf.Addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&client.GrpcAuth{Token: conf.Token}),
	)
	if err != nil {
		log.Fatalln("连接 grpc 服务失败:", err)
	}
	defer conn.Close()
	log.Infoln("GRPC 已连接")

beginSync:

	// 创建 ssh 账号流
	grpcClient := proto.NewSshAccountsClient(conn)
	ctx, cancelWatch := context.WithCancel(context.Background())
	watch, err := grpcClient.Watch(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalln("启动 SSH 账号 Watch 流失败:", err)
	}
	log.Infoln("已创建 SSH 账号 Watch 流")

	// 账号流处理流程
	for {
		msg, err := watch.Recv()
		if err != nil {
			log.Errorln("SSH 账号 Watch 流异常:", err)
			break
		}

		if msg.IsHeartBeat {
			continue
		}

		if err = client.SshAccountSync(msg); err != nil {
			break
		}
	}

	cancelWatch()
	log.Warnln("同步流程出错，即将开始重新同步……")
	time.Sleep(time.Second * 3) // 减小重试频率
	goto beginSync
}
