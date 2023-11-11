package feishu

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dao"
	"github.com/ncuhome/GeniusAuthoritarian/internal/service"
	"gorm.io/gorm"
)

type UserDepartment struct {
	list []uint
}

func (d UserDepartment) Ids() []uint {
	return d.list
}

func (d UserDepartment) Models(uid uint) []dao.UserGroups {
	models := make([]dao.UserGroups, len(d.list))
	for i, gid := range d.list {
		models[i].UID = uid
		models[i].GID = gid
	}
	return models
}

func (d UserDepartment) Sync(tx *gorm.DB, uid uint) error {
	userGroupSrv := service.UserGroupsSrv{DB: tx}
	err := userGroupSrv.DeleteNotInGidSliceByUID(uid, d.list).Error
	if err != nil {
		return err
	}

	existDepartments, err := userGroupSrv.GetIdsForUser(uid)
	if err != nil {
		return err
	}

	gidToAddLength := len(d.list) - len(existDepartments)
	if gidToAddLength > 0 {
		var gidToAdd = make([]uint, gidToAddLength)
		var gidToAddIndex int
		for _, gid := range d.list {
			for _, existGid := range existDepartments {
				if existGid == gid {
					goto next
				}
			}
			gidToAdd[gidToAddIndex] = gid
			gidToAddIndex++
			if gidToAddIndex == gidToAddLength {
				break
			}
		next:
		}
		err = userGroupSrv.CreateAll(UserDepartment{list: gidToAdd}.Models(uid))
		if err != nil {
			return err
		}
	}
	return nil
}
