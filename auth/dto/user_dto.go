package dto

type User struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	Handphone string `json:"handphone"`
	AvatarURL string `json:"avatar_url"`
}
