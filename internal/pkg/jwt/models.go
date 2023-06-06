package jwt

import "github.com/golang-jwt/jwt/v4"

type UserToken struct {
	jwt.RegisteredClaims
	// dao.User.ID
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Groups []string `json:"groups,omitempty"`
}

type LoginToken struct {
	jwt.RegisteredClaims
	// 无意义 ID
	ID uint64 `json:"id"`
}

type LoginTokenClaims struct {
	UID    uint     `json:"uid"`
	Name   string   `json:"name"`
	IP     string   `json:"ip"`
	Groups []string `json:"groups"`

	AppID uint `json:"appID"`
}

type MfaToken struct {
	jwt.RegisteredClaims
	UID uint `json:"uid"`
}

type MfaLoginClaims struct {
	LoginTokenClaims
	Mfa         string `json:"mfa"`
	AppCallback string `json:"appCallback"`
}
