package client

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/sshDev/proto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/linuxUser"
	log "github.com/sirupsen/logrus"
)

// 本地账号字典
var accounts = make(map[string]bool)

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
