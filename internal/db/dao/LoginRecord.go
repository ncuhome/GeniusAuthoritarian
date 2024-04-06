package dao

import (
	"github.com/ncuhome/GeniusAuthoritarian/internal/db/dto"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	"gorm.io/gorm"
	"time"
)

type LoginRecord struct {
	ID        uint  `gorm:"primarykey"`
	CreatedAt int64 `gorm:"index"`

	Destroyed   bool   `gorm:"index;comment:用于加速查询，不能用于销毁登录状态"`
	ValidBefore uint64 `gorm:"index"`

	IP        string
	Useragent string
	Method    string `gorm:"comment:登录方法"`

	// User.ID
	UID  uint  `gorm:"not null;index;column:uid"`
	User *User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	// App.ID 为 null 或 0 表示登录的是后台
	AID *uint `gorm:"column:aid;index"`
	App *App  `gorm:"foreignKey:AID;constraint:OnDelete:CASCADE"`
}

func (a *LoginRecord) sqlJoinApps(tx *gorm.DB) *gorm.DB {
	return tx.Joins("LEFT JOIN apps ON apps.id=login_records.aid")
}

func (a *LoginRecord) sqlGetByUID(tx *gorm.DB) *gorm.DB {
	tx = tx.Model(a).Where(a, "uid").
		Select("login_records.*,IFNULL(apps.name,?) as target", global.ThisAppName)
	tx = a.sqlJoinApps(tx)
	tx = tx.Order("login_records.id DESC")
	return tx
}

func (a *LoginRecord) sqlLoginValid(tx *gorm.DB) *gorm.DB {
	return tx.Where("login_records.valid_before>? AND NOT login_records.destroyed=1", time.Now().Unix())
}

func (a *LoginRecord) sqlGetForCancel(tx *gorm.DB) *gorm.DB {
	tx = tx.Model(&LoginRecord{})
	tx = a.sqlJoinApps(tx)
	tx = a.sqlLoginValid(tx)
	tx = tx.Select("login_records.id", "login_records.valid_before", "apps.app_code")
	return tx
}

func (a *LoginRecord) Insert(tx *gorm.DB) error {
	return tx.Create(a).Error
}

func (a *LoginRecord) UpdateDestroyedByID(tx *gorm.DB) error {
	return tx.Model(a).Update("destroyed", true).Error
}

func (a *LoginRecord) UpdateDestroyedByIDSlice(tx *gorm.DB, ids []uint) error {
	tx = tx.Model(a)
	tx = a.sqlLoginValid(tx)
	return tx.Where("id IN ?", ids).Update("destroyed", true).Error
}

func (a *LoginRecord) GetByUID(tx *gorm.DB, limit int) ([]dto.LoginRecord, error) {
	var t = make([]dto.LoginRecord, 0)
	tx = a.sqlGetByUID(tx).Limit(limit)
	return t, tx.Find(&t).Error
}

func (a *LoginRecord) GetForCancelByAID(tx *gorm.DB) ([]dto.LoginRecordForCancel, error) {
	var t []dto.LoginRecordForCancel
	tx = a.sqlGetForCancel(tx)
	return t, tx.Where(a, "aid").Find(&t).Error
}

func (a *LoginRecord) GetValidForUser(tx *gorm.DB) ([]dto.LoginRecordOnline, error) {
	var t = make([]dto.LoginRecordOnline, 0)
	tx = a.sqlGetByUID(tx)
	tx = a.sqlLoginValid(tx)
	return t, tx.Find(&t).Error
}

func (a *LoginRecord) GetLastMonth(tx *gorm.DB) ([]LoginRecord, error) {
	var t []LoginRecord
	return t, tx.Model(a).Where("created_at<=?", 604800).Order("id DESC").Find(&t).Error
}

func (a *LoginRecord) GetMultipleViewsIds(tx *gorm.DB, apps []App) ([]dto.ViewID, error) {
	var t []dto.ViewID
	tx = tx.Model(a).Select("login_records.id", "login_records.aid")
	tx = a.sqlJoinApps(tx)
	for _, app := range apps {
		tx = tx.Or("apps.id=? AND login_records.id>?", app.ID, app.ViewsID)
	}
	return t, tx.Order("apps.id,login_records.id DESC").Find(&t).Error
}

func (a *LoginRecord) GetViewIds(tx *gorm.DB, startAt uint) ([]uint, error) {
	var t []uint
	tx = tx.Model(a).Select("login_records.id")
	tx = a.sqlJoinApps(tx)
	return t, tx.Where("apps.id=? AND login_records.id>?", a.AID, startAt).Order("login_records.id DESC").Find(&t).Error
}

func (a *LoginRecord) GetAdminViews(tx *gorm.DB, startAt int64) ([]dto.LoginRecordDataView, error) {
	var t = make([]dto.LoginRecordDataView, 0)
	return t, tx.Model(a).Where("created_at>=?", startAt).Order("id DESC").Find(&t).Error
}

func (a *LoginRecord) TakeValidForCancel(tx *gorm.DB) (*dto.LoginRecordForCancel, error) {
	var t dto.LoginRecordForCancel
	tx = a.sqlGetForCancel(tx)
	tx = tx.Where(a, "id", "uid")
	return &t, tx.Take(&t).Error
}
