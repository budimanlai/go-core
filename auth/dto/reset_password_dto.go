package dto

type ResetPasswordRequest struct {
	Channel         string `json:"channel" validate:"required,oneof=phone email"`
	Identifier      string `json:"identifier" validate:"required"`
	TrxID           string `json:"trx_id" validate:"required"`
	Password        string `json:"password" validate:"required,min=7,max=18"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}
