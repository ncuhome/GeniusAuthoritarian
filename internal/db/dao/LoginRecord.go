package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"gorm.io/gorm"
)

type LoginRecordWithForeignKey struct {
	LoginRecord `gorm:"embedded"`
	App         App  `gorm:"-;foreignKey:AID;constraint:OnDelete:CASCADE"`
	User        User `gorm:"-;foreignKey:UID;constraint:OnDelete:CASCADE"`
}

func (a *LoginRecordWithForeignKey) TableName() string {
	return "login_records"
}

type LoginRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID
	UID uint `gorm:"not null;index;column:uid"`
	IP  string
	// App.ID
	AID uint `gorm:"column:aid;not null;index"`
}

func (a *LoginRecord) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *LoginRecord) GetByUID(tx *gorm.DB, limit int) ([]dto.LoginRecord, error) {
	var t = make([]dto.LoginRecord, 0)
	return t, tx.Model(a).Where(a, "UID").Select("login_records.*,IFNULL(apps.name,?) as target", global.ThisAppName).
		Joins("LEFT JOIN apps ON apps.id=login_records.aid").
		Order("login_records.id DESC").Limit(limit).Find(&t).Error
}
