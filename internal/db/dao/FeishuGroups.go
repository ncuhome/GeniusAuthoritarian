package dao

import (
	"gorm.io/gorm"
)

type FeishuGroups struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null;unique"`
	OpenDepartmentId string `gorm:"not null;uniqueInde;type:varchar(255)"`
	// BaseGroup.ID
	GID   uint       `gorm:"uniqueIndex;not null;column:gid"`
	Group *BaseGroup `gorm:"foreignKey:GID;constraint:OnDelete:CASCADE"`
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

func (a *FeishuGroups) GetGroupsByOpenIDSlice(tx *gorm.DB, openID []string) ([]BaseGroup, error) {
	var t []BaseGroup
	groupModel := &BaseGroup{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinFeishuGroups(tx)
	return t, tx.Where("feishu_groups.open_department_id IN ?", openID).Find(&t).Error
}

func (a *FeishuGroups) GetByOpenIDSlice(tx *gorm.DB, openID []string) ([]FeishuGroups, error) {
	var t []FeishuGroups
	return t, tx.Model(a).Where("open_department_id IN ?", openID).Find(&t).Error
}
