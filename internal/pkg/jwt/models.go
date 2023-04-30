package jwt

import "github.com/golang-jwt/jwt/v4"

type UserToken struct {
	jwt.RegisteredClaims
	// UID
	ID uint `json:"id"`
}

type LoginToken struct {
	jwt.RegisteredClaims
	// 无意义 ID
	ID uint64 `json:"id"`
}

type LoginTokenClaims struct {
	UID    uint     `json:"uid"`
	IP     string   `json:"ip"`
	Name   string   `json:"name"`
	Target string   `json:"target"`
	Groups []string `json:"groups"`
}
