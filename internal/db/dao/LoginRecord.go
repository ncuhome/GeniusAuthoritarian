package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"gorm.io/gorm"
)

type LoginRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID
	UID  uint  `gorm:"not null;index;column:uid"`
	User *User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	IP   string
	// App.ID 为 null 或 0 表示登录的是后台
	AID *uint `gorm:"column:aid;index"`
	App *App  `gorm:"foreignKey:AID;constraint:OnDelete:CASCADE"`
}

func (a *LoginRecord) sqlJoinApps(tx *gorm.DB) *gorm.DB {
	return tx.Joins("LEFT JOIN apps ON apps.id=login_records.aid AND apps.deleted_at IS NULL")
}

func (a *LoginRecord) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *LoginRecord) GetByUID(tx *gorm.DB, limit int) ([]dto.LoginRecord, error) {
	var t = make([]dto.LoginRecord, 0)
	tx = tx.Model(a).Where(a, "UID").Select("login_records.*,IFNULL(apps.name,?) as target", global.ThisAppName)
	tx = a.sqlJoinApps(tx)
	tx = tx.Order("login_records.id DESC").Limit(limit)
	return t, tx.Find(&t).Error
}

func (a *LoginRecord) GetLastMonth(tx *gorm.DB) ([]LoginRecord, error) {
	var t []LoginRecord
	return t, tx.Model(a).Where("created_at<=?", 604800).Order("id DESC").Find(&t).Error
}

/*func (a *LoginRecord) GetViewCount(tx *gorm.DB) ([]dto.ViewCount, error) {
	var t []dto.ViewCount
	tx = tx.Model(a).Select("apps.id", "COUNT(login_records.id) AS views")
	tx = a.sqlJoinApps(tx)
	return t, tx.Group("apps.id").Find(&t).Error
}*/

func (a *LoginRecord) GetViewIds(tx *gorm.DB, startAt uint) ([]uint, error) {
	var t []uint
	tx = tx.Model(a).Select("login_record.id")
	tx = a.sqlJoinApps(tx)
	return t, tx.Where("apps.id=? AND login_record.id>?", a.AID, startAt).Order("login_record.id DESC").Find(&t).Error
}
