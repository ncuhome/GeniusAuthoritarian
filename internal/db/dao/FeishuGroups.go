package dao

type FeishuGroups struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null;unique"`
	OpenDepartmentId string `gorm:"not null;uniqueInde;type:varchar(255)"`
	// Group.ID
	GID   uint  `gorm:"uniqueIndex;not null;column:gid"`
	Group Group `gorm:"foreignKey:GID,constraint:RESTRICT"`
}
