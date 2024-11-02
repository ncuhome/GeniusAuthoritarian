package service

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"gorm.io/gorm"
)

var UserGroups = UserGroupsSrv{dao.DB}

type UserGroupsSrv struct {
	*gorm.DB
}

func (a UserGroupsSrv) Begin() (UserGroupsSrv, error) {
	a.DB = a.DB.Begin()
	return a, a.Error
}

// IsUnitMember 是否是指定组组员
func (a UserGroupsSrv) IsUnitMember(uid uint, Unit string) (bool, error) {
	return (&dao.User2Groups{
		UID: uid,
	}).ExistByName(a.DB, Unit)
}

func (a UserGroupsSrv) GetIdsForUser(uid uint) ([]uint, error) {
	return (&dao.User2Groups{UID: uid}).GetGetUserGroupIdsByUID(a.DB)
}

func (a UserGroupsSrv) GetNamesForUser(uid uint) ([]string, error) {
	return (&dao.User2Groups{UID: uid}).GetUserGroupNamesByUID(a.DB)
}

func (a UserGroupsSrv) GetForAppCode(uid uint, appCode string) ([]string, error) {
	var groups []string
	return groups, (&dao.User2Groups{
		UID: uid,
	}).GetUserGroupsForAppCodeByUID(a.DB, appCode).Find(&groups).Error
}

func (a UserGroupsSrv) GetAll() ([]dao.User2Groups, error) {
	return (&dao.User2Groups{}).GetAllNotFrozen(a.DB)
}

func (a UserGroupsSrv) CreateAll(data []dao.User2Groups) error {
	return (&dao.User2Groups{}).InsertAll(a.DB, data)
}

func (a UserGroupsSrv) DeleteByIDSlice(id []uint) error {
	return (&dao.User2Groups{}).DeleteByIDSlice(a.DB, id)
}

func (a UserGroupsSrv) DeleteNotInGidSliceByUID(uid uint, id []uint) *gorm.DB {
	return (&dao.User2Groups{UID: uid}).DeleteNotInGidSliceByUID(a.DB, id)
}
