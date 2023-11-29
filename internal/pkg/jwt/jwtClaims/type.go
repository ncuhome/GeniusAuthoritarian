package jwtClaims

import "github.com/golang-jwt/jwt/v5"

type Claims interface {
	jwt.Claims
	GetType() string
}
type ClaimsUser interface {
	Claims
	GetUID() uint
	GetUserOperateID() uint64
}

// TypedClaims type 字段用于区分不同类型的 token，防止类型窜用导致的安全漏洞
type TypedClaims struct {
	jwt.RegisteredClaims
	Type string `json:"type"`
}

func (c TypedClaims) GetType() string {
	return c.Type
}

type UserClaims struct {
	TypedClaims
	// 用户 ID
	UID           uint   `json:"uid"`
	UserOperateID uint64 `json:"oid"`
}

func (u UserClaims) GetUID() uint {
	return u.UID
}

func (u UserClaims) GetUserOperateID() uint64 {
	return u.UserOperateID
}
