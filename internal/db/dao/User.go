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

func (a *User) limitedUpdateByID(tx *gorm.DB, selected ...interface{}) *gorm.DB {
	return tx.Model(a).Select(selected[0], selected[1:]...).Where(a, "ID").Updates(a)
}

func (a *User) Rename(tx *gorm.DB) *gorm.DB {
	return a.limitedUpdateByID(tx, "Name")
}

func (a *User) UpdatePhone(tx *gorm.DB) *gorm.DB {
	return a.limitedUpdateByID(tx, "Phone")
}
