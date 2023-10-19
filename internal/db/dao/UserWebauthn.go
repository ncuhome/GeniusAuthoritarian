package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

type UserWebauthn struct {
	ID   uint  `gorm:"primarykey"`
	UID  uint  `gorm:"index;not null;column:uid"`
	User *User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`

	Name   string `gorm:"type:varchar(30)"` // 设备名
	CredID string `gorm:"not null"`         // Credential ID

	Credential string `gorm:"not null"` // json marshaled
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
	return t, tx.Model(a).Select("id", "name", "cred_id").
		Where(a, "uid").Order("id DESC").Find(&t).Error
}
