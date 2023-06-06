package dto

type UserProfile struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Mfa        string  `json:"-"`
	MfaEnabled bool    `json:"mfa" gorm:"-"`
	Groups     []Group `json:"groups" gorm:"-"`
}
