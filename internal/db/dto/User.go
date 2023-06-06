package dto

type UserProfile struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Phone  string  `json:"phone"`
	MFA    bool    `json:"mfa"`
	Groups []Group `json:"groups" gorm:"-"`
}
