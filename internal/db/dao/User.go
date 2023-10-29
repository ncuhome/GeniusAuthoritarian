package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"not null"`
	Phone     string         `gorm:"not null;uniqueIndex;type:varchar(15)"`
	MFA       string
	PreferU2F string `gorm:"column:prefer_u2f"`
}

func (a *User) sqlJoinUserGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN user_groups ON user_groups.uid=users.id")
}

func (a *User) sqlJoinUserSshs(tx *gorm.DB) *gorm.DB {
	return tx.Joins("LEFT JOIN user_sshes ON user_sshes.uid=users.id")
}

func (a *User) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *User) InsertAll(tx *gorm.DB, users []User) error {
	return tx.Create(users).Error
}

func (a *User) FirstByID(tx *gorm.DB) error {
	return tx.Where(a, "id").First(a).Error
}

func (a *User) FirstByPhone(tx *gorm.DB) error {
	return tx.First(a, "phone=?", a.Phone).Error
}

func (a *User) FirstProfileByID(tx *gorm.DB) (*dto.UserProfile, error) {
	var t dto.UserProfile
	return &t, tx.Model(a).First(&t, "id=?", a.ID).Error
}

func (a *User) FirstForPasskey(tx *gorm.DB) error {
	return tx.Model(a).Select("name").Where(a, "id").First(a).Error
}

func (a *User) FirstMfa(tx *gorm.DB) error {
	return tx.Model(a).Select("mfa").Where(a, "id").First(a).Error
}

func (a *User) FirstPhoneByID(tx *gorm.DB) error {
	return tx.Model(a).Select("phone").Where(a, "id").First(a).Error
}

func (a *User) GetUnscopedByPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Unscoped().Where("phone IN ?", phone).Find(&t).Error
}

func (a *User) GetNotInPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Where("NOT phone IN ?", phone).Find(&t).Error
}

// GetNoSshDevIds 获取没有分发 ssh 账号的研发部门用户
func (a *User) GetNoSshDevIds(tx *gorm.DB) ([]uint, error) {
	tx = tx.Model(a)
	tx = a.sqlJoinUserGroups(tx)
	tx = a.sqlJoinUserSshs(tx)
	tx = (&UserGroups{}).sqlJoinBaseGroups(tx)

	var t []uint
	return t, tx.Select("users.id").Where("base_groups.name=? AND user_sshes.id IS NULL", departments.UDev).Find(&t).Error
}

func (a *User) U2fStatus(tx *gorm.DB) (*dto.UserU2fStatus, error) {
	var t dto.UserU2fStatus
	return &t, tx.Model(a).Select("users.prefer_u2f AS prefer", "1 AS phone", "(users.mfa IS NOT NULL AND users.mfa!='') AS mfa", "user_webauthns.id AS passkey").
		Joins("LEFT JOIN user_webauthns ON user_webauthns.uid=users.id").
		Where(a, "id").
		Limit(1).Find(&t).Error
}

func (a *User) FrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, "id IN ?", ids).Error
}

func (a *User) UpdateMfa(tx *gorm.DB) error {
	return tx.Model(a).Select("mfa").Where(a, "id").Updates(a).Error
}

func (a *User) UnfrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Unscoped().Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NULL")).Error
}

func (a *User) DelMfa(tx *gorm.DB) error {
	return tx.Model(a).Model(a).Where(a, "id").Update("mfa", "").Error
}
