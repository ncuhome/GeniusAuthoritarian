package dao

import "gorm.io/gorm"

type User struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null;uniqueIndex;type:varchar(15)"`
}

func (a *User) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *User) InsertAll(tx *gorm.DB, users []User) error {
	return tx.Create(users).Error
}

func (a *User) GetByPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Where("phone IN ?", phone).Find(&t).Error
}

func (a *User) GetNotInPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Where("NOT phone IN ?", phone).Find(&t).Error
}

func (a *User) DeleteByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "id IN ?", ids).Error
}
