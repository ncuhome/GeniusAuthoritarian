package response

import "encoding/json"

type RefreshToken struct {
	AccessToken string          `json:"access_token"`
	Payload     json.RawMessage `json:"payload,omitempty"`
}

type VerifyAccessToken struct {
	UID     uint            `json:"uid"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
