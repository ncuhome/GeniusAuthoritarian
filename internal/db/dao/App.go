package dao

import (
	"gorm.io/gorm"
)

type App struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	Name      string `gorm:"not null"`
	AppCode   string `gorm:"not null;uniqueIndex;type:varchar(36)"`
	AppSecret string `gorm:"not null"`
	Callback  string `gorm:"not null"`
}

func (a *App) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}
