package jwtClaims

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

type MfaToken struct {
	TypedClaims
	// 无意义 ID
	ID  uint64 `json:"id"`
	UID uint   `json:"uid"`
}

type U2fToken struct {
	TypedClaims
	// 无意义 ID
	ID  uint64 `json:"id"`
	UID uint   `json:"uid"`
	IP  string `json:"ip"`
}

type RefreshToken struct {
	TypedClaims
	UID     uint   `json:"uid"`
	AppCode string `json:"appCode"`
	Payload string `json:"payload,omitempty"`
}

type AccessToken struct {
	RefreshToken
}
