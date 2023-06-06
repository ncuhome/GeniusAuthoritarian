package response

type VerifyTokenSuccess struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

type ThirdPartyLogin struct {
	Token    string `json:"token"`
	Mfa      bool   `json:"mfa"`
	Callback string `json:"callback,omitempty"`
}
