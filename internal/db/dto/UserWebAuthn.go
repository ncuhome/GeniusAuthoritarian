package dto

type UserCredential struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	CredID string `json:"cred_id"`
}
