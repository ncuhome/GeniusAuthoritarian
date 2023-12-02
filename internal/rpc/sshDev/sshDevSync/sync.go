package sshDevSync

import (
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/redis"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/cronAgent"
	sshTool2 "github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevClient/sshTool"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevModel"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/backoff"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// 研发哥容器内 ssh 账号管理器

func AddSshAccountCron(c *cron.Cron, stat redis.SyncStat) {
	schedule, err := cronAgent.Parser.Parse("0 6 * * *")
	if err != nil {
		log.Fatalln("规划定时同步研发 SSH 账号失败:", err)
	}

	sshAccountBackoff := backoff.New(backoff.Conf{
		Content: stat.Inject(schedule, func() error {
			err := DoSync()
			if err != nil {
				log.Errorln("同步 SSH 账号失败:", err)
			} else {
				log.Infoln("同步 SSH 账号成功")
			}
			return err
		}),
		MaxRetryDelay: 0,
	})

	c.Schedule(schedule, cron.FuncJob(sshAccountBackoff.Start))
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
		userSshToCreate[i], err = sshTool2.NewSshDevModel(randRand, uid)
		if err != nil {
			return err
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

	// 通知 sshDev client
	length := len(userSshToCreate) + len(userToDelete)
	if length != 0 {
		sshRpcMessages := make([]sshDevModel.SshAccountMsg, length)
		i := 0
		for _, userSsh := range userSshToCreate {
			sshRpcMessages[i] = sshDevModel.SshAccountMsg{
				Username:  sshTool2.LinuxAccountName(userSsh.UID),
				PublicKey: userSsh.PublicSsh,
			}
			i++
		}
		for _, userSsh := range userToDelete {
			sshRpcMessages[i] = sshDevModel.SshAccountMsg{
				IsDel:    true,
				Username: sshTool2.LinuxAccountName(userSsh.UID),
			}
			i++
		}

		if err = redis.PublishSshDev(sshRpcMessages); err != nil {
			return err
		}
	}

	if err = userSshSrv.Commit().Error; err != nil {
		return err
	}
	return nil
}
