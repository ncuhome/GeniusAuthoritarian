package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

var UserSsh = UserSshSrv{dao.DB}

type UserSshSrv struct {
	*gorm.DB
}

func (a UserSshSrv) Begin() (UserSshSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a UserSshSrv) GetInvalid() ([]uint, error) {
	return (&dao.User{}).GetNoSshDevIds(a.DB)
}

func (a UserSshSrv) GetAllExist() ([]dto.SshDeploy, error) {
	return (&dao.UserSsh{}).GetAll(a.DB)
}
