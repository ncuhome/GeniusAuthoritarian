package main

import (
	"context"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/linuxUser"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Token string `yaml:"token"`
	Addr  string `yaml:"addr"`
}

func readConfig() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalln("读取配置文件失败:", err)
	}
	defer f.Close()

	var conf Config
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Fatalln("解析配置文件失败:", err)
	}
	return conf
}

func main() {
	log.Infoln("Sys Boost")

	conf := readConfig()

	creds := credentials.NewClientTLSFromCert(nil, "")

	conn, err := grpc.Dial(conf.Addr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&GrpcAuth{Token: conf.Token}))
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

				// 使用 -D 参数创建账号后 shadow 的密码值为 !，将无法使用 ssh 登录
				if err = linuxUser.DelPasswd(account.Username); err != nil {
					log.Fatalf("清除账号 %s 密码失败: %s", account.Username, err)
				}
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
