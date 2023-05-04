package dao

import "gorm.io/gorm"

type UserAvatarWithForeignKey struct {
	UserAvatar `gorm:"embedded"`
	User       User `gorm:"-;foreignKey:UID;constraint:OnDelete:RESTRICT"`
}

func (a *UserAvatarWithForeignKey) TableName() string {
	return "user_avatars"
}

type UserAvatar struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID uint `gorm:"column:uid;not null;uniqueIndex"`
}

func (a *UserAvatar) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}
