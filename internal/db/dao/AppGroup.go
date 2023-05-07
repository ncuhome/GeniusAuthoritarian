package dao

type AppGroupWithForeignKey struct {
	AppGroup `gorm:"embedded"`
	App      App   `gorm:"-;foreignKey:AID;constraint:OnDelete:CASCADE"`
	Group    Group `gorm:"-;foreignKey:GID;constraint:OnDelete:CASCADE"`
}

func (a *AppGroupWithForeignKey) TableName() string {
	return "app_groups"
}

type AppGroup struct {
	ID uint `gorm:"primarykey"`
	// App.ID
	AID uint `gorm:"column:aid;not null;index;index:app_group_idx,unique"`
	// Group.ID
	GID uint `gorm:"column:gid;not null;index;index:app_group_idx,unique"`
}
