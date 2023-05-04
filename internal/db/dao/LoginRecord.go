package dao

import (
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
