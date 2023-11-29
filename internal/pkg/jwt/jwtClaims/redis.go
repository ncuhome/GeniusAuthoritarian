package jwtClaims

type LoginRedis struct {
	UID       uint     `json:"uid"`
	AvatarUrl string   `json:"avatarUrl"`
	Name      string   `json:"name"`
	IP        string   `json:"ip"`
	Groups    []string `json:"groups"`

	AppID uint `json:"appID"`
}

type MfaRedis struct {
	LoginRedis
	Mfa         string `json:"mfa"`
	AppCallback string `json:"appCallback"`
}
