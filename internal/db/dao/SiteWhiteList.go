package dao

import "gorm.io/gorm"

type SiteWhiteList struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    int64
	DomainSuffix string `gorm:"not null;uniqueIndex"`
}

func (a *SiteWhiteList) Insert(db *gorm.DB) error {
	return db.Create(a).Error
}

func (a *SiteWhiteList) Get(db *gorm.DB) ([]SiteWhiteList, error) {
	var t []SiteWhiteList
	return t, db.Model(a).Find(&t).Error
}
