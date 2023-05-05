package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"gorm.io/gorm"
)

var UserAvatar = UserAvatarSrv{dao.DB}

type UserAvatarSrv struct {
	*gorm.DB
}

func (a UserAvatarSrv) Begin() (UserAvatarSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

func (a UserAvatarSrv) DelOldAvatar(uid uint) (*dao.UserAvatar, error) {
	t := dao.UserAvatar{
		UID: uid,
	}
	result := t.DelForUser(a.DB)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &t, nil
	}
	return nil, nil
}
