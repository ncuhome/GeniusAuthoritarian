package dao

import (
	"gorm.io/gorm"
)

type UserGroupsWithForeignKey struct {
	UserGroups `gorm:"embedded"`
	User       User  `gorm:"-;foreignKey:UID;constraint:OnDelete:CASCADE"`
	Group      Group `gorm:"-;foreignKey:GID;constraint:OnDelete:CASCADE"`
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

func (a *UserGroups) sqlJoinUsers(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN users ON users.id=user_groups.uid AND users.deleted_at IS NULL")
}

func (a *UserGroups) sqlGetUserGroupsByUID(tx *gorm.DB) *gorm.DB {
	groupModel := &Group{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinUserGroups(tx)
	return tx.Where("user_groups.uid=?", a.UID)
}

func (a *UserGroups) InsertAll(tx *gorm.DB, data []UserGroups) error {
	return tx.Model(a).Create(data).Error
}

func (a *UserGroups) GetUserGroupsForAppCodeByUID(tx *gorm.DB, appCode string) *gorm.DB {
	appGroupsTx := (&AppGroup{}).GetGroups(tx, appCode).Select("groups.name")
	return a.sqlGetUserGroupsByUID(appGroupsTx)
}

// GetUserGroupsLimited 根据指定组范围获取用户所在组
func (a *UserGroups) GetUserGroupsLimited(tx *gorm.DB, groups []string) ([]Group, error) {
	var t []Group
	groupModel := &Group{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinUserGroups(tx)
	return t, tx.Where("user_groups.uid=? AND groups.name IN ?", a.UID, groups).
		Find(&t).Error
}

func (a *UserGroups) GetAllUnfrozen(tx *gorm.DB) ([]UserGroups, error) {
	var t []UserGroups
	tx = tx.Model(a)
	tx = a.sqlJoinUsers(tx)
	return t, tx.Find(&t).Error
}

func (a *UserGroups) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Delete(a, "id IN ?", ids).Error
}
