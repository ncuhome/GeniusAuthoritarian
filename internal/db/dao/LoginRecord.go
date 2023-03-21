package dao

import "gorm.io/gorm"

type LoginRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	Name      string `gorm:"not null"`
	Referer   string
	Target    string
}

func (a *LoginRecord) Insert(db *gorm.DB) error {
	return db.Create(a).Error
}
