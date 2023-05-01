package dao

import (
	"gorm.io/gorm"
)

type LoginRecordModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID
	UID    uint   `gorm:"not null;index;column:uid"`
	Target string `gorm:"not null"`
	IP     string
}

type LoginRecord struct {
	LoginRecordModel
	User User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
}

func (a *LoginRecord) Insert(tx *gorm.DB) error {
	return tx.Model(a).Create(&a.LoginRecordModel).Error
}
