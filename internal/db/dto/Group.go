package dto

type Group struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GroupRelateApp struct {
	Group
	AppID uint `json:"-"`
}

type GroupWithOrder struct {
	ID   uint
	UID  uint `gorm:"uid"`
	Name string
}
