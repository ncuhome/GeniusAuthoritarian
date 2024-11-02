package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

// BaseGroup 权限组
type BaseGroup struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;uniqueIndex;type:varchar(10)"`
}

func (a *BaseGroup) sqlJoinAppGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN app_groups ON app_groups.gid=base_groups.id")
}

func (a *BaseGroup) sqlJoinUserGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN user2groups ON user2groups.gid=base_groups.id")
}

func (a *BaseGroup) sqlJoinFeishuGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN feishu_groups ON feishu_groups.gid=base_groups.id")
}

func (a *BaseGroup) sqlGetByIds(tx *gorm.DB, ids ...uint) *gorm.DB {
	return tx.Model(a).Where("id IN ?", ids)
}

func (a *BaseGroup) GetAll(tx *gorm.DB) ([]BaseGroup, error) {
	var t []BaseGroup
	return t, tx.Model(a).Find(&t).Error
}

func (a *BaseGroup) GetAllForShow(tx *gorm.DB) ([]dto.Group, error) {
	var t = make([]dto.Group, 0)
	return t, tx.Model(a).Find(&t).Error
}

func (a *BaseGroup) GetByAppIdsRelatedForShow(tx *gorm.DB, apps ...uint) ([]dto.GroupRelateApp, error) {
	var t []dto.GroupRelateApp
	tx = tx.Model(a)
	return t, a.sqlJoinAppGroups(tx).Select("base_groups.*", "app_groups.aid AS app_id").Where("app_groups.aid IN ?", apps).Order("app_groups.id").Find(&t).Error
}

func (a *BaseGroup) GetByIdsForShow(tx *gorm.DB, ids ...uint) ([]dto.Group, error) {
	var t = make([]dto.Group, 0)
	return t, a.sqlGetByIds(tx, ids...).Find(&t).Error
}

func (a *BaseGroup) CreateGroups(tx *gorm.DB, groups []BaseGroup) error {
	return tx.Create(&groups).Error
}
