package tokenStorePoint

type StorePointData struct {
	Iat    int64       `json:"iat"`
	Claims interface{} `json:"claims,omitempty"`
}
