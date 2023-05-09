package dao

import (
	"gorm.io/gorm"
)

type App struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      int64
	Name           string `gorm:"not null"`
	AppCode        string `gorm:"not null;uniqueIndex;type:varchar(36)"`
	AppSecret      string `gorm:"not null"`
	Callback       string `gorm:"not null"`
	PermitAllGroup bool
}

func (a *App) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *App) FirstForLogin(tx *gorm.DB) error {
	return tx.Model(a).Omit("app_secret").Where("app_code=?", a.AppCode).First(a).Error
}

func (a *App) FindSecret(tx *gorm.DB, appCode string) (string, error) {
	var t string
	return t, tx.Model(a).Select("app_secret").Where("app_code = ?", appCode).Find(&t).Error
}

func (a *App) Get(tx *gorm.DB) ([]string, error) {
	var t []string
	return t, tx.Model(a).Select("app_code").Find(&t).Error
}
