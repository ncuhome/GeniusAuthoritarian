package jwtClaims

type UserToken struct {
	UserClaims
	Name   string   `json:"name"`
	Groups []string `json:"groups,omitempty"`
}

type LoginToken struct {
	UserClaims
	ID uint64 `json:"id"`
}

type MfaToken struct {
	UserClaims
	ID uint64 `json:"id"`
}

type U2fToken struct {
	UserClaims
	ID uint64 `json:"id"`
	IP string `json:"ip"`
}

type RefreshToken struct {
	UserClaims
	ID      uint64 `json:"id"`
	AppCode string `json:"appCode"`
	Payload string `json:"payload,omitempty"`
}

type AccessToken struct {
	RefreshToken
}
