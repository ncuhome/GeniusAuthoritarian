package dao

import "gorm.io/gorm"

// Group 该模型仅用于添加数据库约束，请勿用于创建含写入操作的 CRUD 接口
type Group struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;uniqueIndex;type:varchar(10)"`
}

func (a *Group) sqlJoinAppGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN app_groups ON app_groups.gid=groups.id")
}

// sqlJoinApps join AppGroups first
func (a *Group) sqlJoinApps(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN apps ON apps.id=app_groups.aid AND apps.deleted_at IS NULL")
}

func (a *Group) sqlJoinUserGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN user_groups ON user_groups.gid=groups.id")
}

func (a *Group) sqlJoinFeishuGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN feishu_groups ON feishu_groups.gid=groups.id")
}

func (a *Group) GetAll(tx *gorm.DB) *gorm.DB {
	return tx.Model(a)
}

func (a *Group) GetByNames(tx *gorm.DB, groups ...uint) *gorm.DB {
	return tx.Model(a).Where("id IN ?", groups)
}

func (a *Group) CreateGroups(tx *gorm.DB, groups []string) ([]Group, error) {
	var targetGroups = make([]Group, len(groups))
	for i, groupName := range groups {
		targetGroups[i].Name = groupName
	}
	return targetGroups, tx.Create(&targetGroups).Error
}
