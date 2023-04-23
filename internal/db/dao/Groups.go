package dao

type Group struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;uniqueIndex;type:varchar(10)"`
}
