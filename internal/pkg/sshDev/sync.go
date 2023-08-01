package sshDev

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/rpc"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
	"math/rand"
	"time"
)

// 研发哥容器内 ssh 账号管理器

func DoSync() error {
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

	if err = userSshSrv.CreateAll(userSshToCreate); err != nil {
		return err
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
		sshRpcMessages := make([]rpc.SshAccountMsg, length)
		i := 0
		for _, userSsh := range userSshToCreate {
			sshRpcMessages[i] = rpc.SshAccountMsg{
				Username:  rpc.LinuxAccountName(userSsh.UID),
				PublicKey: userSsh.PublicSsh,
			}
			i++
		}
		for _, userSsh := range userToDelete {
			sshRpcMessages[i] = rpc.SshAccountMsg{
				IsDel:     true,
				Username:  rpc.LinuxAccountName(userSsh.UID),
				PublicKey: userSsh.PublicSsh,
			}
			i++
		}
		rpc.MsgChannel <- sshRpcMessages
	}

	return nil
}
