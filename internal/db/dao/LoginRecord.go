package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoginRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID
	UID    uint   `gorm:"not null;index;column:uid"`
	User   *User  `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	Target string `gorm:"not null"`
	IP     string
}

func (a *LoginRecord) Insert(tx *gorm.DB) error {
	return tx.Omit(clause.Associations).Create(a).Error
}
