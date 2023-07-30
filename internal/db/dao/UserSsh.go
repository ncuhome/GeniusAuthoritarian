package dao

import "gorm.io/gorm"

type UserSsh struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID  uint  `gorm:"uniqueIndex;not null;column:uid;"`
	User *User `gorm:"foreignKey:UID;constraint:OnDelete:RESTRICT"`

	PublicPem  string `gorm:"not null"`
	PrivatePem string `gorm:"not null"`

	PublicSsh  string `gorm:"not null"`
	PrivateSsh string `gorm:"not null"`
}

func (a *UserSsh) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}
