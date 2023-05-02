package dto

type UserProfile struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	AvatarID  uint   `json:"avatarID"`
	AvatarUrl string `json:"avatarUrl"`
}
