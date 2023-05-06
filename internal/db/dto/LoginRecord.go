package dto

type LoginRecord struct {
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	Target    string `json:"target"`
	IP        string `json:"ip"`
}
