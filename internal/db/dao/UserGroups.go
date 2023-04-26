package dao

import "gorm.io/gorm"

type UserGroups struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID  uint `gorm:"index;index:user_group_idx,unique;not null;column:uid;"`
	User User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	// Group.ID
	GID   uint  `gorm:"index;index:user_group_idx,unique;not null;column:gid"`
	Group Group `gorm:"foreignKey:GID;constraint:OnDelete:RESTRICT"`
}

func (a *UserGroups) InsertAll(tx *gorm.DB, data []UserGroups) error {
	return tx.Create(data).Error
}

// GetUserGroupsLimited 根据指定组 id 范围获取用户所在组
func (a *UserGroups) GetUserGroupsLimited(tx *gorm.DB, gid []uint) ([]Group, error) {
	var groups []Group
	return groups, tx.Model(&Group{}).
		Joins("INNER JOIN user_groups ug ON ug.gid=groups.id").
		Where("ug.uid=? AND groups.id IN ?", a.UID, gid).
		Find(&groups).Error
}

func (a *UserGroups) GetAll(tx *gorm.DB) ([]UserGroups, error) {
	var t []UserGroups
	return t, tx.Find(&t).Error
}

func (a *UserGroups) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "id IN ?", ids).Error
}
