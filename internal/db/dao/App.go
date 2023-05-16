package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

type AppWithForeignKey struct {
	App  `gorm:"embedded"`
	User User `gorm:"-;foreignKey:UID;constraint:OnDelete:SET NULL"`
}

type App struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	gorm.DeletedAt
	// User.ID 拥有者
	UID            uint   `gorm:"column:uid;index"`
	Name           string `gorm:"not null;uniqueIndex;type:varchar(30)"`
	AppCode        string `gorm:"not null;uniqueIndex;type:varchar(36)"`
	AppSecret      string `gorm:"not null"`
	Callback       string `gorm:"not null"`
	PermitAllGroup bool
}

func (a *App) sqlJoinAppGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN app_groups ON app_groups.aid=apps.id")
}

// join AppGroups first
func (a *App) sqlJoinGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN base_groups ON base_groups.id=app_groups.gid")
}

// join Groups first
func (a *App) sqlJoinUserGroups(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN user_groups ON user_groups.gid=base_groups.id")
}

func (a *App) sqlGetForActionByUID(tx *gorm.DB) *gorm.DB {
	return tx.Model(a).Omit("app_secret").Where("uid=?", a.UID)
}

func (a *App) sqlGetByUIDForShow(tx *gorm.DB) *gorm.DB {
	return a.sqlGetForActionByUID(tx).Order("id DESC")
}

func (a *App) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *App) NameExist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("1").Where("name=?", a.Name).Limit(1).Find(&t).Error
}

func (a *App) FirstDetailedByIdAndUID(tx *gorm.DB) (*dto.AppShowDetail, error) {
	var t dto.AppShowDetail
	return &t, tx.Model(a).Where("id=? AND uid=?", a.ID, a.UID).First(&t).Error
}

func (a *App) FirstForLogin(tx *gorm.DB) error {
	return tx.Model(a).Omit("app_secret").Where("app_code=?", a.AppCode).First(a).Error
}

func (a *App) FirstAppKeyPairByID(tx *gorm.DB) error {
	return tx.Model(a).Select("app_code", "app_secret").Where("id=?", a.ID).First(a).Error
}

func (a *App) GetAppCode(tx *gorm.DB) ([]string, error) {
	var t []string
	return t, tx.Model(a).Select("app_code").Find(&t).Error
}

func (a *App) GetByUIDForAction(tx *gorm.DB) ([]App, error) {
	var t []App
	return t, a.sqlGetForActionByUID(tx).Find(&t).Error
}

func (a *App) GetByUIDForShow(tx *gorm.DB) ([]dto.AppShowOwner, error) {
	var t = make([]dto.AppShowOwner, 0)
	return t, a.sqlGetByUIDForShow(tx).Find(&t).Error
}

func (a *App) GetByUIDForShowDetailed(tx *gorm.DB) ([]dto.AppShowDetail, error) {
	var t = make([]dto.AppShowDetail, 0)
	return t, a.sqlGetByUIDForShow(tx).Find(&t).Error
}

func (a *App) GetPermitAll(tx *gorm.DB) ([]dto.AppShow, error) {
	var t = make([]dto.AppShow, 0)
	return t, tx.Model(a).Where("permit_all_group=?", true).Find(&t).Error
}

func (a *App) GetAccessible(tx *gorm.DB) ([]dto.AppShowWithGroup, error) {
	// 后续需要额外归类，故不 make
	var t []dto.AppShowWithGroup
	tx = tx.Model(a).Select("apps.*", "base_groups.id AS group_id", "base_groups.name as group_name")
	tx = a.sqlJoinAppGroups(tx)
	tx = a.sqlJoinGroups(tx)
	tx = a.sqlJoinUserGroups(tx)
	return t, tx.Where("user_groups.uid=?", a.UID).Find(&t).Error
}

func (a *App) DeleteByIdForUID(tx *gorm.DB) *gorm.DB {
	return tx.Model(a).Where("id=? AND uid=?", a.ID, a.UID).Delete(a)
}

func (a *App) UpdatesByID(tx *gorm.DB) error {
	return tx.Model(a).Select("name", "callback", "permit_all_group").Where("id=?", a.ID).Updates(a).Error
}
