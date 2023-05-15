package dto

type Group struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GroupRelateApp struct {
	Group
	AppID uint `json:"-"`
}
