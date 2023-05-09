package service

import (
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

func (a UserSrv) UserInfoForAppCode(phone, appCode string) (*dao.User, []dao.Group, error) {
	var user = dao.User{
		Phone: phone,
	}
	e := user.FirstByPhone(a.DB)
	if e != nil {
		return nil, nil, e
	}

	var groups []dao.Group
	return &user, groups, (&dao.UserGroups{
		UID: user.ID,
	}).GetUserGroupsForAppCodeByUID(a.DB, appCode).Find(&groups).Error
}

func (a UserSrv) UserProfile(uid uint) (*dto.UserProfile, error) {
	profile, e := (&dao.User{ID: uid}).FirstProfileByID(a.DB)
	if e != nil {
		return nil, e
	}

	profile.Groups = make([]dto.Group, 0)
	return profile, (&dao.UserGroups{
		UID: uid,
	}).GetUserGroupsByUID(a.DB).Find(&profile.Groups).Error
}
