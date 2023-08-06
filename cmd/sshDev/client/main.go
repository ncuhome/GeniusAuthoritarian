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
	"time"
)

// 账号同步记录
var accounts = make(map[string]bool)

func init() {
	go linuxUser.DaemonSshd()
}

func main() {
	log.Infoln("Sys Boost")

	// 读取配置
	conf := readConfig()

	// 连接 grpc
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(conf.Addr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&GrpcAuth{Token: conf.Token}))
	if err != nil {
		log.Fatalln("连接 grpc 服务失败:", err)
	}
	defer conn.Close()
	log.Infoln("GRPC 已连接")

beginSync:

	// 创建 ssh 账号流
	client := proto.NewSshAccountsClient(conn)
	ctx, cancelWatch := context.WithCancel(context.Background())
	watch, err := client.Watch(ctx, &emptypb.Empty{})
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

		if err = SshAccountSync(msg); err != nil {
			log.Warnln("同步流程出错，即将开始重新同步……")
			break
		}
	}

	cancelWatch()
	time.Sleep(time.Second * 3) // 减小重试频率
	goto beginSync
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

func DoAccountDelete(username string, logger *log.Entry) error {
	if err := linuxUser.Delete(username); err != nil {
		logger.Errorln("删除账号失败:", err)
		return err
	}
	delete(accounts, username)
	return nil
}

func SshAccountSet(account *proto.SshAccount) error {
	logger := log.WithField("username", account.Username)

	if account.IsDel {
		err := DoAccountDelete(account.Username, logger)
		if err != nil {
			return err
		}
	} else {
		exist, err := linuxUser.Exist(account.Username)
		if err != nil {
			logger.Errorln("检查账号存在失败:", err)
			return err
		}
		if !exist {
			err = linuxUser.Create(account.Username)
			if err != nil {
				logger.Errorln("创建账号失败:", err)
				return err
			} else {
				accounts[account.Username] = true
			}
			logger.Infoln("用户已创建")

			// 使用 -D 参数创建账号后 shadow 的密码值为 !，将无法使用 ssh 登录
			if err = linuxUser.DelPasswd(account.Username); err != nil {
				logger.Errorln("清除密码失败:", err)
				return err
			}
		} else {
			if err = linuxUser.DelPasswd(account.Username); err != nil {
				logger.Errorln("清除密码失败:", err)
				return err
			}
		}
		err = linuxUser.PrepareSshDir(account.Username)
		if err != nil {
			logger.Errorln("创建 .ssh 失败:", err)
			return err
		}
		err = linuxUser.WriteAuthorizedKeys(account.Username, account.PublicKey)
		if err != nil {
			logger.Errorln("写入 authorized_keys 失败:", err)
			return err
		}
		logger.Infoln("SSH key 已配置")
	}

	return nil
}

func SshAccountSync(msg *proto.AccountStream) error {
	if msg.IsInit {
		// 查找本地多出来的账号
		accountsInit := make(map[string]bool, len(msg.Accounts))
		for _, account := range msg.Accounts {
			accountsInit[account.Username] = true
		}
		for username := range accounts {
			delete(accountsInit, username)
		}
		if len(accountsInit) != 0 {
			newAccountsArray := make([]*proto.SshAccount, len(msg.Accounts)+len(accountsInit))
			for i, account := range msg.Accounts {
				newAccountsArray[i] = account
			}
			i := len(msg.Accounts)
			for username := range accountsInit {
				newAccountsArray[i] = &proto.SshAccount{
					IsDel:    true,
					Username: username,
				}
				i++
			}
		}
	}

	for _, account := range msg.Accounts {
		err := SshAccountSet(account)
		if err != nil {
			return err
		}
	}
	return nil
}
