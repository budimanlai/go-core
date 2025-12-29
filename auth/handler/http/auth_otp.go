package http

import (
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// RequestOtp godoc
// @Summary      Request OTP
// @Description  Generate and send OTP to user's phone number
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        otpRequest  body      dto.OtpRequest  true  "OTP Request"
// @Success      200         {object}  dto.OtpResponse
// @Failure      400         {object}  response.ErrorResponse
// @Router       /auth/otp/request [post]
func (h *AuthHandler) RequestOtp(ctx *fiber.Ctx) error {
	var req dto.OtpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorI18n(ctx, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// validate request
	if err := validator.ValidateStructWithContext(ctx, &req); err != nil {
		return response.ValidationErrorI18n(ctx, err)
	}

	resp, err := h.OtpUC.GenerateOTP(ctx.Context(), req)
	if err != nil {
		return response.BadRequestI18n(ctx, err.Error(), nil)
	}

	return response.SuccessI18n(ctx, "app.success", resp)
}

// RequestOtp godoc
// @Summary      Request OTP
// @Description  Generate and send OTP to user's phone number
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        otpRequest  body      dto.OtpRequest  true  "OTP Request"
// @Success      200         {object}  dto.OtpResponse
// @Failure      400         {object}  response.ErrorResponse
// @Router       /auth/otp/request [post]
func (h *AuthHandler) StatusOTP(ctx *fiber.Ctx) error {
	var req dto.OtpStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorI18n(ctx, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// validate request
	if err := validator.ValidateStructWithContext(ctx, &req); err != nil {
		return response.ValidationErrorI18n(ctx, err)
	}

	valid, err := h.OtpUC.Status(ctx.Context(), req.Identifier, req.TrxID)
	if err != nil {
		return response.BadRequestI18n(ctx, err.Error(), nil)
	}
	if valid == false {
		return response.ErrorI18n(ctx, fiber.StatusUnauthorized, "auth.error.invalid_otp", nil)
	}

	var out dto.OtpStatusResponse = dto.OtpStatusResponse{
		Identifier: req.Identifier,
		TrxID:      req.TrxID,
		Valid:      valid,
	}

	return response.SuccessI18n(ctx, "app.success", out)
}

// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verify the provided OTP for the user's phone number
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        otpVerifyRequest  body      dto.OtpVerifyRequest  true  "OTP Verify Request"
// @Success      200              {object}  dto.OtpVerifyResponse
// @Failure      400              {object}  response.ErrorResponse
// @Router       /auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(ctx *fiber.Ctx) error {
	var req dto.OtpVerifyRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorI18n(ctx, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// validate request
	if err := validator.ValidateStructWithContext(ctx, &req); err != nil {
		return response.ValidationErrorI18n(ctx, err)
	}

	err := h.OtpUC.VerifyOtp(ctx.Context(), req.Identifier, req.TrxID, req.PinCode)
	if err != nil {
		return response.BadRequestI18n(ctx, err.Error(), nil)
	}

	var out dto.OtpVerifyResponse = dto.OtpVerifyResponse{
		Identifier: req.Identifier,
		TrxID:      req.TrxID,
		Valid:      true,
	}

	return response.SuccessI18n(ctx, "app.success", out)
}
