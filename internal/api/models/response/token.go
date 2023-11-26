package response

type RefreshToken struct {
	AccessToken string `json:"access_token"`
	Payload     string `json:"payload,omitempty"`
}

type VerifyAccessToken struct {
	UID     uint   `json:"uid"`
	Payload string `json:"payload,omitempty"`
}
