package main

import (
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/linuxUser"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
)

var Address string
var Token string

func init() {
	Address = os.Getenv("Addr")
	if Address == "" {
		log.Fatalln("连接地址不能为空，请配置环境变量 Addr")
	}
	Token = os.Getenv("Token")
}

func main() {
	log.Infoln("Sys Boost")

	creds := credentials.NewClientTLSFromCert(nil, "")

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&GrpcAuth{Token: Token}))
	if err != nil {
		log.Fatalln("连接 grpc 服务失败:", err)
	}
	defer conn.Close()

	client := proto.NewSshAccountsClient(conn)
	watch, err := client.Watch(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalln("启动 SSH 账号 Watch 流失败:", err)
	}

	for {
		account, err := watch.Recv()
		if err != nil {
			log.Fatalln("SSH 账号流异常:", err)
		}
		if account.Username == "" {
			// 心跳
			continue
		}
		if account.IsDel {
			err = linuxUser.Delete(account.Username)
			if err != nil {
				log.Fatalf("删除账号 %s 失败: %s", account.Username, err)
			}
		} else {
			exist, err := linuxUser.Exist(account.Username)
			if err != nil {
				log.Fatalf("检查账号 %s 存在失败: %s", account.Username, err)
			}
			if !exist {
				err = linuxUser.Create(account.Username)
				if err != nil {
					log.Fatalf("创建账号 %s 失败: %s", account.Username, err)
				}
				log.Infof("用户 %s 已创建", account.Username)
			}
			err = linuxUser.PrepareSshDir(account.Username)
			if err != nil {
				log.Fatalf("创建账号 %s .ssh 失败: %s", account.Username, err)
			}
			err = linuxUser.WriteAuthorizedKeys(account.Username, account.PublicKey)
			if err != nil {
				log.Fatalf("写入账号 %s authorized_keys 失败: %s", account.Username, err)
			}
			log.Infof("用户 %s SSH key 已配置", account.Username)
		}
	}
}

type GrpcAuth struct {
	Token string
}

func (a *GrpcAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": a.Token}, nil
}

func (a *GrpcAuth) RequireTransportSecurity() bool {
	return true
}
