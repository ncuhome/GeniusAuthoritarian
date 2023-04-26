package dao

import "gorm.io/gorm"

type LoginRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID
	UID    uint   `gorm:"not null;index;column:uid"`
	User   User   `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	Target string `gorm:"not null"`
}

func (a *LoginRecord) Insert(db *gorm.DB) error {
	return db.Create(a).Error
}
