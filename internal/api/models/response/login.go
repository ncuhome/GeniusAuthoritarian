package response

type VerifyTokenSuccess struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}
