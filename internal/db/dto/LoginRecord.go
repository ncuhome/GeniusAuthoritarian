package dto

type LoginRecord struct {
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	IP        string `json:"ip"`
	Target    string `json:"target"`
}

type ViewCount struct {
	// App.ID
	ID    uint   `json:"id"`
	Views uint64 `json:"views"`
}
