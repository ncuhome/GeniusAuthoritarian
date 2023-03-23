package jwt

import "github.com/golang-jwt/jwt/v4"

type JWT struct {
	jwt.RegisteredClaims
	Groups []string `json:"groups"`
}
