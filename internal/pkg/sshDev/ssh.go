package sshDev

import (
	"fmt"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
)

// 研发哥容器内 ssh 账号管理器

func LinuxAccountName(uid uint) string {
	return "home" + fmt.Sprint(uid)
}

func DoSync() error {
	/*userSshSrv, err := service.UserSsh.Begin()
	if err != nil {
		return err
	}
	defer userSshSrv.Rollback()*/

	users, err := service.UserSsh.GetInvalid()
	if err != nil {
		return err
	}

}
