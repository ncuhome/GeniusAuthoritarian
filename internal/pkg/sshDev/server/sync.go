package server

import (
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/agent"
	rpc2 "github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/server/rpc"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/sshTool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// 研发哥容器内 ssh 账号管理器

func AddSshAccountCron(spec string) {
	_, err := agent.AddRegular(&agent.Event{
		T: spec,
		E: func() {
			err := DoSync()
			if err != nil {
				log.Errorln("同步 SSH 账号失败:", err)
			} else {
				log.Infoln("同步 SSH 账号成功")
			}
		},
	})
	if err != nil {
		log.Fatalln("添加 SSH 账号同步任务失败:", err)
	}
}

func DoSync() error {
	defer tool.Recover()

	users, err := service.UserSsh.GetToCreateUid()
	if err != nil {
		return err
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randRand := rand.New(randSource)

	// 生成密钥对
	userSshToCreate := make([]dao.UserSsh, len(users))
	for i, uid := range users {
		sshKey, err := ed25519.Generate(randRand)
		if err != nil {
			return err
		}

		publicPem, privatePem, err := sshKey.MarshalPem()
		if err != nil {
			return err
		}
		publicSsh, privateSsh, err := sshKey.MarshalSSH()
		if err != nil {
			return err
		}

		userSshToCreate[i] = dao.UserSsh{
			UID:        uid,
			PublicPem:  string(publicPem),
			PrivatePem: string(privatePem),
			PublicSsh:  string(publicSsh),
			PrivateSsh: string(privateSsh),
		}
	}

	userSshSrv, err := service.UserSsh.Begin()
	if err != nil {
		return err
	}
	defer userSshSrv.Rollback()

	if len(userSshToCreate) != 0 {
		if err = userSshSrv.CreateAll(userSshToCreate); err != nil {
			return err
		}
	}

	userToDelete, err := userSshSrv.DeleteInvalid()
	if err != nil {
		return err
	}

	if err = userSshSrv.Commit().Error; err != nil {
		return err
	}

	// 通知 sshDev client
	length := len(userSshToCreate) + len(userToDelete)
	if length != 0 {
		sshRpcMessages := make([]rpc2.SshAccountMsg, length)
		i := 0
		for _, userSsh := range userSshToCreate {
			sshRpcMessages[i] = rpc2.SshAccountMsg{
				Username:  sshTool.LinuxAccountName(userSsh.UID),
				PublicKey: userSsh.PublicSsh,
			}
			i++
		}
		for _, userSsh := range userToDelete {
			sshRpcMessages[i] = rpc2.SshAccountMsg{
				IsDel:     true,
				Username:  sshTool.LinuxAccountName(userSsh.UID),
				PublicKey: userSsh.PublicSsh,
			}
			i++
		}
		rpc2.MsgChannel <- sshRpcMessages
	}

	return nil
}
