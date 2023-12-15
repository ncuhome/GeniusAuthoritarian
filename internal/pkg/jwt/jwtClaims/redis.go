package jwtClaims

type LoginRedis struct {
	UID       uint     `json:"uid"`
	AvatarUrl string   `json:"avatarUrl"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`

	AppID uint `json:"appID"`

	IP        string `json:"ip"`
	Useragent string `json:"useragent"`
}

type MfaRedis struct {
	LoginRedis
	Mfa         string `json:"mfa"`
	AppCallback string `json:"appCallback"`
}
