package service

import (
	"github.com/Mmx233/daoUtil"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

var User = UserSrv{dao.DB}

type UserSrv struct {
	*gorm.DB
}

func (a UserSrv) Begin() (UserSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a UserSrv) GetUnscopedUserByPhoneSlice(phone []string) ([]dao.User, error) {
	return (&dao.User{}).GetUnscopedByPhoneSlice(a.DB, phone)
}

func (a UserSrv) GetUserNotInPhoneSlice(phone []string) ([]dao.User, error) {
	return (&dao.User{}).GetNotInPhoneSlice(a.DB, phone)
}

func (a UserSrv) CreateAll(users []dao.User) error {
	return (&dao.User{}).InsertAll(a.DB, users)
}

func (a UserSrv) FrozeByIDSlice(id []uint) error {
	return (&dao.User{}).FrozeByIDSlice(a.DB, id)
}

func (a UserSrv) UnFrozeByIDSlice(id []uint) error {
	return (&dao.User{}).UnfrozeByIDSlice(a.DB, id)
}

func (a UserSrv) UserInfo(phone string) (*dao.User, error) {
	var user = dao.User{
		Phone: phone,
	}
	return &user, user.FirstByPhone(a.DB)
}

func (a UserSrv) U2fStatus(id uint) (*dto.UserU2fStatus, error) {
	return (&dao.User{ID: id}).U2fStatus(a.DB)
}

func (a UserSrv) UserInfoByID(id uint) (*dao.User, error) {
	var user = dao.User{
		ID: id,
	}
	return &user, user.FirstByID(a.DB)
}

func (a UserSrv) UserProfile(id uint) (*dto.UserProfile, error) {
	profile, err := (&dao.User{ID: id}).FirstProfileByID(a.DB)
	if err != nil {
		return nil, err
	}

	profile.MfaEnabled = profile.Mfa != ""

	profile.Groups, err = (&dao.UserGroups{
		UID: id,
	}).GetUserGroupsForShowByUID(a.DB)
	return profile, err
}

func (a UserSrv) MfaExist(id uint, opts ...daoUtil.ServiceOpt) (bool, error) {
	mfa, err := a.FirstMfa(id, opts...)
	return mfa != "", err
}

func (a UserSrv) FirstMfa(id uint, opts ...daoUtil.ServiceOpt) (string, error) {
	var t = dao.User{
		ID: id,
	}
	return t.MFA, t.FirstMfa(daoUtil.TxOpts(a.DB, opts...))
}

func (a UserSrv) SetMfaSecret(id uint, secret string) error {
	return (&dao.User{ID: id, MFA: secret}).UpdateMfa(a.DB)
}

func (a UserSrv) DelMfa(id uint) error {
	return (&dao.User{ID: id}).DelMfa(a.DB)
}

func (a UserSrv) FirstPhoneByID(id uint) (string, error) {
	model := dao.User{ID: id}
	return model.Phone, model.FirstPhoneByID(a.DB)
}

func (a UserSrv) UpdateUserPreferU2F(id uint, prefer string) error {
	return (&dao.User{ID: id, PreferU2F: prefer}).UpdateU2fPreferByID(a.DB)
}
