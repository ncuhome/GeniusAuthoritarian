package dto

type UserCredential struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CreatedAt  int    `json:"created_at"`
	LastUsedAt int    `json:"last_used_at"`
}
