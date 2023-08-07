package client

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/linux"
	log "github.com/sirupsen/logrus"
)

// 本地账号字典，值为账号是否完成可登录配置
var accounts = make(map[string]bool)

func DoAccountDelete(logger *log.Entry, username string) error {
	if err := linux.DeleteUser(username); err != nil {
		logger.Errorln("删除账号失败:", err)
		return err
	}
	delete(accounts, username)
	return nil
}

func DoUserProcessKill(logger *log.Entry, username string) error {
	if err := linux.UserKillAll(username); err != nil {
		logger.Errorln("结束用户进程失败:", err)
		return err
	}
	logger.Infoln("用户进程已清理")
	return nil
}

func LinuxUserPreset(logger *log.Entry, username string) error {
	// 使用 -D 参数创建账号后 shadow 的密码值为 !，将无法使用 ssh 登录
	err := linux.DelUserPasswd(username)
	if err != nil {
		logger.Errorln("清除密码失败:", err)
		return err
	}

	accounts[username] = true
	return nil
}

// SshAccountSet 返回 error 代表流需要重启
func SshAccountSet(account *proto.SshAccount) error {
	logger := log.WithField("username", account.Username)

	var err error
	if account.IsDel {
		if err = DoUserProcessKill(logger, account.Username); err != nil {
			return err
		}
		if err = DoAccountDelete(logger, account.Username); err != nil {
			return err
		}
		logger.Infoln("用户已删除")
	} else if account.IsKill {
		_ = DoUserProcessKill(logger, account.Username)
	} else {
		ready, exist := accounts[account.Username]
		if !exist {
			err = linux.CreateUser(account.Username)
			if err != nil {
				logger.Errorln("创建账号失败:", err)
				return err
			}
			accounts[account.Username] = false
			logger.Infoln("用户已创建")

			if err = LinuxUserPreset(logger, account.Username); err != nil {
				return err
			}
		} else if !ready {
			if err = LinuxUserPreset(logger, account.Username); err != nil {
				return err
			}
		}
		err = linux.PrepareSshDir(account.Username)
		if err != nil {
			logger.Errorln("创建 .ssh 失败:", err)
			return err
		}
		err = linux.WriteAuthorizedKeys(account.Username, account.PublicKey)
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
		accountsShouldDelete := make(map[string]bool, len(accounts))
		for username := range accounts {
			accountsShouldDelete[username] = true
		}
		for _, account := range msg.Accounts {
			delete(accountsShouldDelete, account.Username)
		}
		if len(accountsShouldDelete) != 0 {
			newAccountsArray := make([]*proto.SshAccount, len(msg.Accounts)+len(accountsShouldDelete))
			for i, account := range msg.Accounts {
				newAccountsArray[i] = account
			}
			i := len(msg.Accounts)
			for username := range accountsShouldDelete {
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
