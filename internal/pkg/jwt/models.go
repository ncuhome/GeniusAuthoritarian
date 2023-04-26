package jwt

import "github.com/golang-jwt/jwt/v4"

type RefreshToken struct {
	jwt.RegisteredClaims
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

type LoginToken struct {
	jwt.RegisteredClaims
	// 无意义 ID
	ID uint64 `json:"id"`
}

type LoginTokenClaims struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}
