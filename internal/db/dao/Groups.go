package dao

import "gorm.io/gorm"

type Group struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;uniqueIndex;type:varchar(10)"`
}

func (a *Group) GetAll(tx *gorm.DB) ([]Group, error) {
	var t []Group
	return t, tx.Model(a).Find(&t).Error
}
