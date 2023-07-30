package dao

import "gorm.io/gorm"

type UserSsh struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID  uint  `gorm:"uniqueIndex;not null;column:uid;"`
	User *User `gorm:"foreignKey:UID;constraint:OnDelete:RESTRICT"`

	PublicKey  string `gorm:"not null"`
	PrivateKey string `gorm:"not null"`
}

func (a *UserSsh) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}
