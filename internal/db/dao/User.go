package dao

import "gorm.io/gorm"

type User struct {
	ID uint `gorm:"primarykey"`
	gorm.DeletedAt
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null;uniqueIndex;type:varchar(15)"`
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

func (a *User) GetUnscopedByPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Unscoped().Where("phone IN ?", phone).Find(&t).Error
}

func (a *User) GetNotInPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Where("NOT phone IN ?", phone).Find(&t).Error
}

func (a *User) FrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "id IN ?", ids).Error
}

func (a *User) UnfrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Unscoped().Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NULL")).Error
}
