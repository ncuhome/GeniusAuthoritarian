package dao

import (
	"gorm.io/gorm"
)

type FeishuGroupModel struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null;unique"`
	OpenDepartmentId string `gorm:"not null;uniqueInde;type:varchar(255)"`
	// Group.ID
	GID uint `gorm:"uniqueIndex;not null;column:gid"`
}

type FeishuGroups struct {
	FeishuGroupModel
	Group Group `gorm:"foreignKey:GID;constraint:RESTRICT"`
}

func (a *FeishuGroups) GetAll(tx *gorm.DB) ([]FeishuGroupModel, error) {
	var t []FeishuGroupModel
	return t, tx.Model(a).Find(&t).Error
}

func (a *FeishuGroups) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "ID IN ?", ids).Error
}

func (a *FeishuGroups) CreateAll(tx *gorm.DB, data []FeishuGroupModel) error {
	return tx.Model(a).Create(data).Error
}

func (a *FeishuGroups) GetGroupsByOpenIDSlice(tx *gorm.DB, openID []string) ([]Group, error) {
	var t []Group
	return t, tx.Model(&Group{}).
		Joins("INNER JOIN feishu_groups fg ON fg.gid=groups.id").
		Where("fg.open_department_id IN ?", openID).Find(&t).Error
}

func (a *FeishuGroups) GetByOpenIDSlice(tx *gorm.DB, openID []string) ([]FeishuGroupModel, error) {
	var t []FeishuGroupModel
	return t, tx.Model(a).Where("open_department_id IN ?", openID).Find(&t).Error
}
