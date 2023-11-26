package response

type VerifyTokenSuccess struct {
	UserID    uint     `json:"userID"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	AvatarUrl string   `json:"avatarUrl"`
}

type ThirdPartyLogin struct {
	Token    string `json:"token"`
	Mfa      bool   `json:"mfa"`
	Callback string `json:"callback,omitempty"`
}
