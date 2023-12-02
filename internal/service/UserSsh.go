package service

import (
	"github.com/Mmx233/daoUtil"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/sshDevClient/sshTool"
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

func (a UserSshSrv) Exist(uid uint, opt ...daoUtil.ServiceOpt) (bool, error) {
	return (&dao.UserSsh{UID: uid}).Exist(daoUtil.TxOpts(a.DB, opt...))
}

// GetToCreateUid 获取没有生成 ssh 账号的 uid
func (a UserSshSrv) GetToCreateUid() ([]uint, error) {
	return (&dao.User{}).GetNoSshDevIds(a.DB)
}

func (a UserSshSrv) DeleteInvalid() ([]dao.UserSsh, error) {
	model := dao.UserSsh{}
	invalid, err := model.GetInvalid(a.DB)
	if err != nil {
		return nil, err
	}

	if len(invalid) != 0 {
		idSlice := make([]uint, len(invalid))
		for i, userSsh := range invalid {
			idSlice[i] = userSsh.ID
		}
		if err = model.DeleteByIds(a.DB, idSlice...); err != nil {
			return nil, err
		}
	}

	return invalid, nil
}

func (a UserSshSrv) GetAllExist() ([]dto.SshDeploy, error) {
	return (&dao.UserSsh{}).GetAll(a.DB)
}

func (a UserSshSrv) CreateAll(data []dao.UserSsh) error {
	return (&dao.UserSsh{}).InsertAll(a.DB, data)
}

func (a UserSshSrv) FirstSshSecretsForUserShow(uid uint) (*dto.SshSecrets, error) {
	model := dao.UserSsh{UID: uid}
	err := model.FirstForUserShow(a.DB)
	if err != nil {
		return nil, err
	}
	return &dto.SshSecrets{
		Username: sshTool.LinuxAccountName(model.UID),
		Pem: dto.SshKeyPair{
			Public:  model.PublicPem,
			Private: model.PrivatePem,
		},
		Ssh: dto.SshKeyPair{
			Public:  model.PublicSsh,
			Private: model.PrivateSsh,
		},
	}, nil
}

func (a UserSshSrv) UpdateKeys(uid uint, publicPem, privatePem, publicSsh, privateSsh string) error {
	return (&dao.UserSsh{
		UID:        uid,
		PublicPem:  publicPem,
		PrivatePem: privatePem,
		PublicSsh:  publicSsh,
		PrivateSsh: privateSsh,
	}).UpdateByUid(a.DB)
}
