package dto

type OtpRequest struct {
	Channel    string `json:"channel" validate:"required,oneof=phone email"`
	Identifier string `json:"identifier" validate:"required"`
	TrxID      string `json:"trx_id" validate:"required"`
}

type OtpResponse struct {
	Identifier string `json:"identifier"`
	TrxID      string `json:"trx_id"`
	WaUrl      string `json:"wa_url"`
}

type OtpVerifyRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	TrxID      string `json:"trx_id" validate:"required"`
	PinCode    string `json:"pin_code" validate:"required"`
}

type OtpVerifyResponse struct {
	Identifier string `json:"identifier" validate:"required"`
	TrxID      string `json:"trx_id" validate:"required"`
	Valid      bool   `json:"valid"`
}

type OtpStatusRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	TrxID      string `json:"trx_id" validate:"required"`
}

type OtpStatusResponse struct {
	Identifier string `json:"identifier" validate:"required"`
	TrxID      string `json:"trx_id" validate:"required"`
	Valid      bool   `json:"valid"`
}
