package dao

import (
	"gorm.io/gorm"
)

type FeishuGroupsWithForeignKey struct {
	FeishuGroups `gorm:"embedded"`
	Group        Group `gorm:"-;foreignKey:GID;constraint:OnDelete:CASCADE"`
}

func (a *FeishuGroupsWithForeignKey) TableName() string {
	return "feishu_groups"
}

type FeishuGroups struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null;unique"`
	OpenDepartmentId string `gorm:"not null;uniqueInde;type:varchar(255)"`
	// Group.ID
	GID uint `gorm:"uniqueIndex;not null;column:gid"`
}

func (a *FeishuGroups) GetAll(tx *gorm.DB) ([]FeishuGroups, error) {
	var t []FeishuGroups
	return t, tx.Model(a).Find(&t).Error
}

func (a *FeishuGroups) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "ID IN ?", ids).Error
}

func (a *FeishuGroups) CreateAll(tx *gorm.DB, data []FeishuGroups) error {
	return tx.Model(a).Create(data).Error
}

func (a *FeishuGroups) GetGroupsByOpenIDSlice(tx *gorm.DB, openID []string) ([]Group, error) {
	var t []Group
	return t, tx.Model(&Group{}).
		Joins("INNER JOIN feishu_groups fg ON fg.gid=groups.id").
		Where("fg.open_department_id IN ?", openID).Find(&t).Error
}

func (a *FeishuGroups) GetByOpenIDSlice(tx *gorm.DB, openID []string) ([]FeishuGroups, error) {
	var t []FeishuGroups
	return t, tx.Model(a).Where("open_department_id IN ?", openID).Find(&t).Error
}
