package dao

type SiteWhiteList struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    int64
	DomainSuffix string `gorm:"not null;uniqueIndex"`
}
