package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"not null"`
	Phone     string         `gorm:"not null;uniqueIndex;type:varchar(15)"`
	MFA       string
}

func (a *User) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *User) InsertAll(tx *gorm.DB, users []User) error {
	return tx.Create(users).Error
}

func (a *User) FirstByPhone(tx *gorm.DB) error {
	return tx.First(a, "phone=?", a.Phone).Error
}

func (a *User) FirstProfileByID(tx *gorm.DB) (*dto.UserProfile, error) {
	var t dto.UserProfile
	return &t, tx.Model(a).First(&t, "id=?", a.ID).Error
}

func (a *User) GetUnscopedByPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Unscoped().Where("phone IN ?", phone).Find(&t).Error
}

func (a *User) GetNotInPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Where("NOT phone IN ?", phone).Find(&t).Error
}

func (a *User) UpdateMfa(tx *gorm.DB) error {
	return tx.Model(a).Select("mfa").Where(a, "id").Updates(a).Error
}

func (a *User) MfaExist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("mfa").Where(a, "id").First(&t).Error
}

func (a *User) FirstMfa(tx *gorm.DB) error {
	return tx.Model(a).Select("mfa").Where(a, "id").First(a).Error
}

func (a *User) FrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "id IN ?", ids).Error
}

func (a *User) UnfrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Unscoped().Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NULL")).Error
}

func (a *User) DelMfa(tx *gorm.DB) error {
	return tx.Model(a).Model(a).Where(a, "id").Update("mfa", gorm.Expr("NULL")).Error
}
