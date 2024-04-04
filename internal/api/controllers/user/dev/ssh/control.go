package controllers

import (
	"github.com/Mmx233/daoUtil"
	"github.com/gin-gonic/gin"
	"github.com/ncuhome/GeniusAuthoritarian/internal/api/callback"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevClient/sshTool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevModel"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/internal/tools"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/ed25519"
)

func ResetSshKeyPair(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID

	userSshSrv, err := service.UserSsh.Begin()
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}
	defer userSshSrv.Rollback()

	exist, err := userSshSrv.Exist(uid, daoUtil.LockForUpdate)
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	} else if !exist {
		callback.Error(c, callback.ErrMfaNotExist)
		return
	}

	keyPair, err := ed25519.Generate()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}
	publicPem, privatePem, err := keyPair.MarshalPem()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}
	publicSsh, privateSsh, err := keyPair.MarshalSSH()
	if err != nil {
		callback.Error(c, callback.ErrUnexpected, err)
		return
	}

	publicPemStr, privatePemStr := string(publicPem), string(privatePem)
	publicSshStr, privateSshStr := string(publicSsh), string(privateSsh)

	if err = userSshSrv.UpdateKeys(uid,
		publicPemStr, privatePemStr,
		publicSshStr, privateSshStr,
	); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = redis.PublishSshDev([]sshDevModel.SshAccountMsg{
		{
			Username:  sshTool.LinuxAccountName(uid),
			PublicKey: publicSshStr,
		},
	}); err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	if err = userSshSrv.Commit().Error; err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Success(c, &dto.SshSecrets{
		Username: sshTool.LinuxAccountName(uid),
		Pem: dto.SshKeyPair{
			Public:  publicPemStr,
			Private: privatePemStr,
		},
		Ssh: dto.SshKeyPair{
			Public:  publicSshStr,
			Private: privateSshStr,
		},
	})
}

func KillAllProcess(c *gin.Context) {
	uid := tools.GetUserInfo(c).UID

	err := redis.PublishSshDev([]sshDevModel.SshAccountMsg{
		{
			IsKill:   true,
			Username: sshTool.LinuxAccountName(uid),
		},
	})
	if err != nil {
		callback.Error(c, callback.ErrDBOperation, err)
		return
	}

	callback.Default(c)
}
