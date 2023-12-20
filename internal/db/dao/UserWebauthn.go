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

	CredID     string `gorm:"not null;index;type:varchar(255)"`
	Name       string `gorm:"type:varchar(30)"` // 设备名
	Credential string `gorm:"not null"`         // json marshaled
}

func (a *UserWebauthn) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *UserWebauthn) GetByUID(tx *gorm.DB) ([]string, error) {
	var t []string
	return t, tx.Model(a).Select("credential").
		Where(a, "uid").Find(&t).Error
}

func (a *UserWebauthn) GetByUidForShow(tx *gorm.DB) ([]dto.UserCredential, error) {
	var t = make([]dto.UserCredential, 0)
	return t, tx.Model(a).Select("id", "name", "created_at", "last_used_at").
		Where(a, "uid").Order("id DESC").Find(&t).Error
}

func (a *UserWebauthn) UpdatesByID(tx *gorm.DB) error {
	return tx.Model(a).Updates(a).Error
}

func (a *UserWebauthn) UpdateLastUsedAt(tx *gorm.DB) *gorm.DB {
	return tx.Model(a).Where(a, "cred_id", "uid").Update("last_used_at", a.LastUsedAt)
}

func (a *UserWebauthn) Exist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(&UserWebauthn{}).Select("1").
		Where(a, "id", "uid").Limit(1).Find(&t).Error
}

func (a *UserWebauthn) Delete(tx *gorm.DB) *gorm.DB {
	return tx.Model(&UserWebauthn{}).Where(a, "id", "uid").Delete(nil)
}
