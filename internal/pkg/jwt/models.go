package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims interface {
	jwt.Claims
	GetType() string
}

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

type LoginTokenClaims struct {
	UID       uint     `json:"uid"`
	AvatarUrl string   `json:"avatarUrl"`
	Name      string   `json:"name"`
	IP        string   `json:"ip"`
	Groups    []string `json:"groups"`

	AppID uint `json:"appID"`
}

type MfaToken struct {
	TypedClaims
	UID uint `json:"uid"`
}

type MfaLoginClaims struct {
	LoginTokenClaims
	Mfa         string `json:"mfa"`
	AppCallback string `json:"appCallback"`
}

type U2fClaims struct {
	TypedClaims
	UID uint   `json:"uid"`
	IP  string `json:"ip"`
}
