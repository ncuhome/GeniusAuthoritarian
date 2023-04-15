package jwt

import "github.com/golang-jwt/jwt/v4"

type RefreshToken struct {
	jwt.RegisteredClaims
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

type AuthToken struct {
	jwt.RegisteredClaims
	ID     uint64   `json:"id"`
	Groups []string `json:"groups"`
}
