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

// GetUserFor 根据指定组 id 范围获取用户所在组
func (a *UserGroups) GetUserFor(tx *gorm.DB, gid []uint) ([]Group, error) {
	var groups []Group
	return groups, tx.Model(&Group{}).
		Joins("INNER JOIN user_groups ug ON ug.gid=groups.id").
		Where("ug.uid=? AND groups.id IN ?", a.UID, gid).
		Find(&groups).Error
}

// DelUser 删除某用户所有组关系
func (a *UserGroups) DelUser(tx *gorm.DB) error {
	return tx.Delete(&UserGroups{}, "uid=?", a.UID).Error
}

// AddUser 批量创建用户组关系
func (a *UserGroups) AddUser(tx *gorm.DB, gid []uint) error {
	var userGroups = make([]UserGroups, len(gid))
	for i, id := range gid {
		userGroups[i].UID = a.UID
		userGroups[i].GID = id
	}
	return tx.Create(userGroups).Error
}
