package dao

type User struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null;uniqueIndex;type:varchar(15)"`
}
