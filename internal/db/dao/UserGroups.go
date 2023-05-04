package dao

import (
	"gorm.io/gorm"
)

type UserGroupsWithForeignKey struct {
	UserGroups `gorm:"embedded"`
	User       User  `gorm:"-;foreignKey:UID;constraint:OnDelete:CASCADE"`
	Group      Group `gorm:"-;foreignKey:GID;constraint:OnDelete:RESTRICT"`
}

func (a *UserGroupsWithForeignKey) TableName() string {
	return "user_groups"
}

type UserGroups struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID uint `gorm:"index;index:user_group_idx,unique;not null;column:uid;"`
	// Group.ID
	GID uint `gorm:"index;index:user_group_idx,unique;not null;column:gid"`
}

func (a *UserGroups) InsertAll(tx *gorm.DB, data []UserGroups) error {
	return tx.Model(a).Create(data).Error
}

func (a *UserGroups) GetUserGroupsByUID(tx *gorm.DB) ([]Group, error) {
	var t []Group
	return t, tx.Model(&Group{}).
		Joins("INNER JOIN user_groups ug ON ug.gid=groups.id").
		Where("ug.uid=?", a.UID).Find(&t).Error
}

// GetUserGroupsLimited 根据指定组范围获取用户所在组
func (a *UserGroups) GetUserGroupsLimited(tx *gorm.DB, groups []string) ([]Group, error) {
	var t []Group
	return t, tx.Model(&Group{}).
		Joins("INNER JOIN user_groups ug ON ug.gid=groups.id").
		Where("ug.uid=? AND groups.name IN ?", a.UID, groups).
		Find(&t).Error
}

func (a *UserGroups) GetAllUnfrozen(tx *gorm.DB) ([]UserGroups, error) {
	var t []UserGroups
	return t, tx.Model(a).
		Joins("INNER JOIN users u ON u.id=user_groups.uid AND u.deleted_at IS NULL").
		Find(&t).Error
}

func (a *UserGroups) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Delete(a, "id IN ?", ids).Error
}
