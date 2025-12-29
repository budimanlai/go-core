package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}
type LoginResponse struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Handphone string `json:"handphone"`
	Fullname  string `json:"fullname"`
	Token     Token  `json:"token"`
}
