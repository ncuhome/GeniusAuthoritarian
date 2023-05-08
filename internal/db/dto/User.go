package dto

type UserProfile struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Phone  string  `json:"phone"`
	Groups []Group `json:"groups" gorm:"-"`
}
