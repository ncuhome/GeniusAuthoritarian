package dao

import "gorm.io/gorm"

// Group 该模型仅用于添加数据库约束，请勿用于创建含写入操作的 CRUD 接口
type Group struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;uniqueIndex;type:varchar(10)"`
}

func (a *Group) GetAll(tx *gorm.DB) ([]Group, error) {
	var t []Group
	return t, tx.Model(a).Find(&t).Error
}

func (a *Group) GetByNames(tx *gorm.DB, groups ...string) *gorm.DB {
	return tx.Model(a).Where("name IN ?", groups)
}

func (a *Group) CreateGroups(tx *gorm.DB, groups []string) ([]Group, error) {
	var targetGroups = make([]Group, len(groups))
	for i, groupName := range groups {
		targetGroups[i].Name = groupName
	}
	return targetGroups, tx.Create(&targetGroups).Error
}
