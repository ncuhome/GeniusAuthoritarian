package dto

type UserThirdPartyIdentity struct {
	Phone string
}

type UserProfile struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	AvatarUrl  string  `json:"avatar_url"`
	Mfa        string  `json:"-"`
	MfaEnabled bool    `json:"mfa" gorm:"-"`
	Groups     []Group `json:"groups" gorm:"-"`
}

type UserU2fStatus struct {
	Prefer  string `json:"prefer"`
	Phone   bool   `json:"phone"`
	Mfa     bool   `json:"mfa"`
	Passkey bool   `json:"passkey"`
}

type UserInfoPublic struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	AvatarUrl string  `json:"avatar_url"`
	Groups    []Group `json:"groups" gorm:"-"`
}
