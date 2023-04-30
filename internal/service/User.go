package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
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

func (a UserSrv) GetUserByPhoneSlice(phone []string) ([]dao.User, error) {
	return (&dao.User{}).GetByPhoneSlice(a.DB, phone)
}

func (a UserSrv) GetUserNotInPhoneSlice(phone []string) ([]dao.User, error) {
	return (&dao.User{}).GetNotInPhoneSlice(a.DB, phone)
}

func (a UserSrv) CreateAll(users []dao.User) error {
	return (&dao.User{}).InsertAll(a.DB, users)
}

func (a UserSrv) DeleteByIDSlice(id []uint) error {
	return (&dao.User{}).FrozeByIDSlice(a.DB, id)
}

func (a UserSrv) UserInfo(phone string) (*dao.User, []dao.Group, error) {
	var user = dao.User{
		Phone: phone,
	}
	e := user.First(a.DB)
	if e != nil {
		return nil, nil, e
	}

	userGroups, e := (&dao.UserGroups{
		UID: user.ID,
	}).GetUserGroupsByUID(a.DB)
	return &user, userGroups, e
}
