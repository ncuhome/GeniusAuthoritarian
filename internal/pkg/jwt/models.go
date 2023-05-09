package jwt

import "github.com/golang-jwt/jwt/v4"

type UserToken struct {
	jwt.RegisteredClaims
	// Uid
	ID uint `json:"id"`
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
