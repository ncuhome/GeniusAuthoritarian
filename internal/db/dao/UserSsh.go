package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	"gorm.io/gorm"
)

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

func (a *UserSsh) InsertAll(tx *gorm.DB, models []UserSsh) error {
	return tx.Create(models).Error
}

// GetInvalid 获取应该清除的 user ssh
func (a *UserSsh) GetInvalid(tx *gorm.DB) ([]UserSsh, error) {
	tx = tx.Model(a).Joins("LEFT JOIN users ON users.id=user_sshes.uid") // 无 deleted_at 限制，相当于 UnScoped
	tx = tx.Joins("LEFT JOIN user_groups ON user_groups.uid=user_sshes.uid")
	tx = tx.Joins("LEFT JOIN base_groups ON base_groups.id=user_groups.gid AND base_groups.name=?", departments.UDev)

	var t []UserSsh
	return t, tx.Group("user_sshes.id,users.deleted_at").Having("COUNT(base_groups.id)=0 OR users.deleted_at IS NOT NULL").Find(&t).Error
}

func (a *UserSsh) GetAll(tx *gorm.DB) ([]dto.SshDeploy, error) {
	var t []dto.SshDeploy
	return t, tx.Model(a).Find(&t).Error
}
