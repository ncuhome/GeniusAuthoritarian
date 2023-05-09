package dao

import "gorm.io/gorm"

type AppGroupWithForeignKey struct {
	AppGroup `gorm:"embedded"`
	App      App   `gorm:"-;foreignKey:AID;constraint:OnDelete:CASCADE"`
	Group    Group `gorm:"-;foreignKey:GID;constraint:OnDelete:CASCADE"`
}

func (a *AppGroupWithForeignKey) TableName() string {
	return "app_groups"
}

type AppGroup struct {
	ID uint `gorm:"primarykey"`
	// App.ID
	AID uint `gorm:"column:aid;not null;index;index:app_group_idx,unique"`
	// Group.ID
	GID uint `gorm:"column:gid;not null;index;index:app_group_idx,unique"`
}

func (a *AppGroup) GetGroups(tx *gorm.DB, appCode string) *gorm.DB {
	return tx.Model(&Group{}).Select("groups.name").
		Joins("INNER JOIN app_groups ag ON ag.gid=groups.id").
		Where("ag.app_code = ?", appCode)
}
