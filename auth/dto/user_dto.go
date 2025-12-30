package dto

type User struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	Handphone string `json:"handphone"`
	AvatarURL string `json:"avatar_url"`
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Fullname        string `json:"fullname" validate:"required"`
	Handphone       string `json:"handphone" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FromIP          string
	UserAgent       string
}
