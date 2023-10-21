package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

type UserWebauthn struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  int
	LastUsedAt int
	UID        uint  `gorm:"index;not null;column:uid"`
	User       *User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`

	Name       string `gorm:"type:varchar(30)"` // 设备名
	Credential string `gorm:"not null"`         // json marshaled
}

func (a *UserWebauthn) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *UserWebauthn) GetByUID(tx *gorm.DB) ([]string, error) {
	var t []string
	return t, tx.Model(a).Select("credential").Where(a, "uid").Find(&t).Error
}

func (a *UserWebauthn) GetByUidForShow(tx *gorm.DB) ([]dto.UserCredential, error) {
	var t = make([]dto.UserCredential, 0)
	return t, tx.Model(a).Select("id", "name", "").
		Where(a, "uid").Order("id DESC").Find(&t).Error
}

func (a *UserWebauthn) Updates(tx *gorm.DB) error {
	return tx.Where(a, "id", "uid").Updates(a).Error
}

func (a *UserWebauthn) Exist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(a).Where(a, "id", "uid").Limit(1).Find(&t).Find(&t).Error
}
