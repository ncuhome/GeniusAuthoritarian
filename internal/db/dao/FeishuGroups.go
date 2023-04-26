package dao

import "gorm.io/gorm"

type FeishuGroups struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null;unique"`
	OpenDepartmentId string `gorm:"not null;uniqueInde;type:varchar(255)"`
	// Group.ID
	GID   uint  `gorm:"uniqueIndex;not null;column:gid"`
	Group Group `gorm:"foreignKey:GID,constraint:RESTRICT"`
}

func (a *FeishuGroups) DeleteAll(tx *gorm.DB) error {
	return tx.Unscoped().Delete(&FeishuGroups{}).Error
}

func (a *FeishuGroups) CreateAll(tx *gorm.DB, data []FeishuGroups) error {
	return tx.Create(data).Error
}
