package dao

import "gorm.io/gorm"

type AppGroup struct {
	ID uint `gorm:"primarykey"`
	// App.ID
	AID uint `gorm:"column:aid;not null;index;index:app_group_idx,unique"`
	App *App `gorm:"foreignKey:AID;constraint:OnDelete:CASCADE"`
	// BaseGroup.ID
	GID       uint       `gorm:"column:gid;not null;index;index:app_group_idx,unique"`
	BaseGroup *BaseGroup `gorm:"foreignKey:GID;constraint:OnDelete:CASCADE"`
}

func (a *AppGroup) sqlJoinApps(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN apps ON apps.id=app_groups.aid AND apps.deleted_at IS NULL")
}

func (a *AppGroup) sqlGetGroupsJoined(tx *gorm.DB) *gorm.DB {
	groupModel := &BaseGroup{}
	tx = tx.Model(groupModel)
	tx = groupModel.sqlJoinAppGroups(tx)
	tx = a.sqlJoinApps(tx)
	return tx
}

func (a *AppGroup) sqlGetGroupsByAppCode(tx *gorm.DB, appCode string) *gorm.DB {
	return a.sqlGetGroupsJoined(tx).Where("apps.app_code=?", appCode)
}

func (a *AppGroup) DeleteByAID(tx *gorm.DB) error {
	return tx.Model(a).Where("aid=?", a.AID).Delete(a).Error
}

func (a *AppGroup) DeleteByGidForApp(tx *gorm.DB, gids ...uint) error {
	return tx.Model(a).Where("aid=? AND gid IN ?", a.AID, gids).Delete(a).Error
}
