package dao

import (
	"gorm.io/gorm"
)

type AppWithForeignKey struct {
	App  `gorm:"embedded"`
	User User `gorm:"-;foreignKey:UID;constraint:OnDelete:SET NULL"`
}

type App struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	gorm.DeletedAt
	// User.ID 拥有者
	UID            uint   `gorm:"column:uid;index"`
	Name           string `gorm:"not null;uniqueIndex;type:varchar(30)"`
	AppCode        string `gorm:"not null;uniqueIndex;type:varchar(36)"`
	AppSecret      string `gorm:"not null"`
	Callback       string `gorm:"not null"`
	PermitAllGroup bool
}

func (a *App) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *App) NameExist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("1").Where("name=?", a.Name).Limit(1).Find(&t).Error
}

func (a *App) FirstForLogin(tx *gorm.DB) error {
	return tx.Model(a).Omit("app_secret").Where("app_code=?", a.AppCode).First(a).Error
}

func (a *App) FirstAppKeyPairByID(tx *gorm.DB) error {
	return tx.Model(a).Select("app_code,app_secret").Where("id=?", a.ID).First(a).Error
}

func (a *App) Get(tx *gorm.DB) ([]string, error) {
	var t []string
	return t, tx.Model(a).Select("app_code").Find(&t).Error
}
