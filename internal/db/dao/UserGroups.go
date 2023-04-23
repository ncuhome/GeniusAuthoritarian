package dao

type UserGroups struct {
	ID uint `gorm:"primarykey"`
	// User.ID
	UID  uint `gorm:"index;index:user_group_idx,unique;not null;column:uid;"`
	User User `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE"`
	// Group.ID
	GID   uint  `gorm:"index;index:user_group_idx,unique;not null;column:gid"`
	Group Group `gorm:"foreignKey:GID;constraint:RESTRICT"`
}
