package jwtClaims

import "github.com/golang-jwt/jwt/v5"

type Claims interface {
	jwt.Claims
	GetType() string
}

// TypedClaims type 字段用于区分不同类型的 token，防止类型窜用导致的安全漏洞
type TypedClaims struct {
	jwt.RegisteredClaims
	Type string `json:"type"`
}

func (a *TypedClaims) GetType() string {
	return a.Type
}
