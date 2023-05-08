package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

type LoginRecordWithForeignKey struct {
	LoginRecord `gorm:"embedded"`
	User        User `gorm:"-;foreignKey:UID;constraint:OnDelete:CASCADE"`
}

func (a *LoginRecordWithForeignKey) TableName() string {
	return "login_records"
}

type LoginRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID
	UID    uint   `gorm:"not null;index;column:uid"`
	Target string `gorm:"not null"`
	IP     string
}

func (a *LoginRecord) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *LoginRecord) GetByUID(tx *gorm.DB, limit int) ([]dto.LoginRecord, error) {
	var t = make([]dto.LoginRecord, 0)
	return t, tx.Model(a).Where(a, "UID").Order("id DESC").Limit(limit).Find(&t).Error
}
