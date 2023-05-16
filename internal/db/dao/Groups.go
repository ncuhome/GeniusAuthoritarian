package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

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

func (a *Group) sqlGetByIds(tx *gorm.DB, ids ...uint) *gorm.DB {
	return tx.Model(a).Where("id IN ?", ids)
}

func (a *Group) GetAll(tx *gorm.DB) ([]Group, error) {
	var t []Group
	return t, tx.Model(a).Find(&t).Error
}

func (a *Group) GetAllForShow(tx *gorm.DB) ([]dto.Group, error) {
	var t = make([]dto.Group, 0)
	return t, tx.Model(a).Find(&t).Error
}

func (a *Group) GetByAppIdsRelatedForShow(tx *gorm.DB, apps ...uint) ([]dto.GroupRelateApp, error) {
	var t []dto.GroupRelateApp
	tx = tx.Model(a)
	return t, a.sqlJoinAppGroups(tx).Select("`groups`.*", "app_groups.aid AS app_id").Where("app_groups.aid IN ?", apps).Order("app_groups.id").Find(&t).Error
}

func (a *Group) GetByIdsForShow(tx *gorm.DB, ids ...uint) ([]dto.Group, error) {
	var t = make([]dto.Group, 0)
	return t, a.sqlGetByIds(tx, ids...).Find(&t).Error
}

func (a *Group) CreateGroups(tx *gorm.DB, groups []string) ([]Group, error) {
	var targetGroups = make([]Group, len(groups))
	for i, groupName := range groups {
		targetGroups[i].Name = groupName
	}
	return targetGroups, tx.Create(&targetGroups).Error
}
