package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/departments"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID        uint           `gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name      string `gorm:"not null"`
	Phone     string `gorm:"not null;uniqueIndex;type:varchar(15)"`
	AvatarUrl string

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

func (a *User) Exist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("1").Where(a, "id").
		Limit(1).Find(&t).Error
}

func (a *User) FirstByID(tx *gorm.DB) error {
	return tx.Take(a, a.ID).Error
}

func (a *User) FirstByPhone(tx *gorm.DB) error {
	return tx.Take(a, "phone=?", a.Phone).Error
}

func (a *User) FirstProfileByID(tx *gorm.DB) (*dto.UserProfile, error) {
	var t dto.UserProfile
	return &t, tx.Model(a).Take(&t, a.ID).Error
}

func (a *User) FirstForPasskey(tx *gorm.DB) error {
	return tx.Select("name").Take(a, a.ID).Error
}

func (a *User) FirstMfa(tx *gorm.DB) error {
	return tx.Select("mfa").Take(a, a.ID).Error
}

func (a *User) FirstPhoneByID(tx *gorm.DB) error {
	return tx.Select("phone").Take(a, a.ID).Error
}

func (a *User) GetUnscopedByPhoneSlice(tx *gorm.DB, phone []string) ([]User, error) {
	var t []User
	return t, tx.Model(a).Unscoped().Where("phone IN ?", phone).
		Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "FIELD(phone,?)", Vars: []interface{}{phone}, WithoutParentheses: true}}).
		Find(&t).Error
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

func (a *User) GetUserInfoPubByIds(tx *gorm.DB, ids ...uint) ([]dto.UserInfoPublic, error) {
	var t = make([]dto.UserInfoPublic, 0)
	return t, tx.Model(a).Where("id IN ?", ids).
		Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{ids}, WithoutParentheses: true}}).
		Find(&t).Error
}

func (a *User) U2fStatus(tx *gorm.DB) (*dto.UserU2fStatus, error) {
	var t dto.UserU2fStatus
	return &t, tx.Model(&User{}).Select("users.prefer_u2f AS prefer", "1 AS phone", "(users.mfa IS NOT NULL AND users.mfa!='') AS mfa", "user_webauthns.id IS NOT NULL AS passkey").
		Joins("LEFT JOIN user_webauthns ON user_webauthns.uid=users.id").
		Limit(1).Find(&t, a.ID).Error
}

func (a *User) FrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Delete(a, ids).Error
}

func (a *User) FrozeByPhone(tx *gorm.DB) *gorm.DB {
	return tx.Delete(a, "phone=?", a.Phone)
}

func (a *User) UpdateMfa(tx *gorm.DB) error {
	return tx.Model(a).Select("mfa").Updates(a).Error
}

func (a *User) UnfrozeByIDSlice(tx *gorm.DB, ids []uint) error {
	return tx.Model(a).Unscoped().Where("id IN ?", ids).Update("deleted_at", gorm.Expr("NULL")).Error
}

func (a *User) UpdateU2fPreferByID(tx *gorm.DB) error {
	return tx.Model(a).Update("prefer_u2f", a.PreferU2F).Error
}

func (a *User) UpdateAllInfoByID(tx *gorm.DB) *gorm.DB {
	return tx.Model(a).Select("name", "phone", "avatar_url").Updates(a)
}

func (a *User) DelMfa(tx *gorm.DB) error {
	return tx.Model(a).Update("mfa", "").Error
}
