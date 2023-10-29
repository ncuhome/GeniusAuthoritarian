package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

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

type UserToken struct {
	TypedClaims
	// dao.User.ID
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Groups []string `json:"groups,omitempty"`
}

type LoginToken struct {
	TypedClaims
	// 无意义 ID
	ID uint64 `json:"id"`
}

type LoginRedisClaims struct {
	UID       uint     `json:"uid"`
	AvatarUrl string   `json:"avatarUrl"`
	Name      string   `json:"name"`
	IP        string   `json:"ip"`
	Groups    []string `json:"groups"`

	AppID uint `json:"appID"`
}

type MfaToken struct {
	TypedClaims
	// 无意义 ID
	ID  uint64 `json:"id"`
	UID uint   `json:"uid"`
}

type MfaRedisClaims struct {
	LoginRedisClaims
	Mfa         string `json:"mfa"`
	AppCallback string `json:"appCallback"`
}

type U2fToken struct {
	TypedClaims
	UID uint   `json:"uid"`
	IP  string `json:"ip"`
}
