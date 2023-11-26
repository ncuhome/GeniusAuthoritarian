package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"gorm.io/gorm"
)

type App struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int64
	// User.ID 拥有者
	UID            uint   `gorm:"column:uid;index"`
	User           *User  `gorm:"foreignKey:UID;constraint:OnDelete:SET NULL"`
	Name           string `gorm:"not null;uniqueIndex;type:varchar(20)"`
	AppCode        string `gorm:"not null;uniqueIndex;type:varchar(36)"`
	AppSecret      string `gorm:"not null"`
	Callback       string `gorm:"not null"`
	PermitAllGroup bool
	// 以下仅用于导航标识
	LinkOff bool   `gorm:"index"`
	Views   uint64 `gorm:"index"`
	// LoginRecord.ID 统计用，无需约束
	ViewsID uint
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

func (a *App) sqlGetForWithGroup(tx *gorm.DB) *gorm.DB {
	tx = tx.Model(a).Select("apps.*", "base_groups.id AS group_id", "base_groups.name as group_name")
	tx = a.sqlJoinAppGroups(tx)
	tx = a.sqlJoinGroups(tx)
	return tx.Order("base_groups.id,apps.id")
}

func (a *App) sqlOrderForNav(tx *gorm.DB) *gorm.DB {
	return tx.Order("apps.link_off,apps.views DESC,apps.name")
}

func (a *App) Exist(tx *gorm.DB) (bool, error) {
	var t bool
	return t, tx.Model(a).Select("1").Where(a).Limit(1).Find(&t).Error
}

func (a *App) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *App) FirstDetailedByIdAndUID(tx *gorm.DB) (*dto.AppShowDetail, error) {
	var t dto.AppShowDetail
	return &t, tx.Model(a).Where(a, "id", "uid").First(&t).Error
}

func (a *App) UserAccessible(tx *gorm.DB) (bool, error) {
	var t bool
	tx = tx.Model(a).Select("1")
	tx = a.sqlJoinAppGroups(tx)
	tx = a.sqlJoinGroups(tx)
	tx = a.sqlJoinUserGroups(tx)
	return t, tx.Where("apps.id=? AND user_groups.uid=?", a.ID, a.UID).First(&t).Error
}

func (a *App) FirstByID(tx *gorm.DB) error {
	return tx.Model(a).Omit("app_secret").First(a, a.ID).Error
}

func (a *App) FirstAppCodeByID(tx *gorm.DB) error {
	return tx.Model(a).Where(a, "id", "uid").Select("app_code").First(a).Error
}

func (a *App) FirstCallbackByID(tx *gorm.DB) error {
	return tx.Model(a).Select("callback").First(a, a.ID).Error
}

func (a *App) FirstByAppCode(tx *gorm.DB) error {
	return tx.Model(a).Omit("app_secret").Where("app_code=?", a.AppCode).First(a).Error
}

func (a *App) FirstAppKeyPairByID(tx *gorm.DB) error {
	return tx.Model(a).Select("app_code", "app_secret").First(a, a.ID).Error
}

func (a *App) FirstAppKeyPairByAppCode(tx *gorm.DB) error {
	return tx.Model(a).Select("app_code", "app_secret").
		First(a, "app_code=?", a.AppCode).Error
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
	tx = tx.Model(a).Where("permit_all_group=?", true)
	tx = a.sqlOrderForNav(tx)
	return t, tx.Find(&t).Error
}

func (a *App) GetAccessible(tx *gorm.DB) ([]dto.AppShowWithGroup, error) {
	// 后续需要二次归类，故不 make
	var t []dto.AppShowWithGroup
	tx = a.sqlGetForWithGroup(tx)
	tx = a.sqlJoinUserGroups(tx)
	tx = a.sqlOrderForNav(tx)
	return t, tx.Where("user_groups.uid=?", a.UID).Find(&t).Error
}

func (a *App) GetAllWithGroup(tx *gorm.DB) ([]dto.AppShowWithGroup, error) {
	var t []dto.AppShowWithGroup
	tx = a.sqlGetForWithGroup(tx)
	return t, tx.Find(&t).Error
}

func (a *App) GetForUpdateView(tx *gorm.DB) ([]App, error) {
	var t []App
	return t, tx.Model(a).Select("id", "views_id", "views").Where("link_off IS NULL OR link_off=?", false).Order("id").Find(&t).Error
}

func (a *App) DeleteByIdForUID(tx *gorm.DB) error {
	return tx.Model(a).Where(a, "id", "uid").Delete(a).Error
}

func (a *App) UpdatesByID(tx *gorm.DB) error {
	return tx.Model(a).Select("name", "callback", "permit_all_group").Where(a, "id").Updates(a).Error
}

func (a *App) UpdateViewByID(tx *gorm.DB) error {
	return tx.Model(a).Select("views", "views_id").Updates(a).Error
}

func (a *App) UpdateLinkOff(tx *gorm.DB) error {
	return tx.Model(a).Where(a, "id", "uid").Update("link_off", a.LinkOff).Error
}
