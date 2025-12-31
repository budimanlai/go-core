package http

import (
	"github.com/budimanlai/go-core/auth/domain/usecase"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	UserUC        usecase.UserUsecase
	UserSessionUC usecase.UserSessionUsecase
	OtpUC         usecase.OtpUsecase
}

func NewAuthHandler(userUsecase usecase.UserUsecase, userSessionUsecase usecase.UserSessionUsecase, otpUsecase usecase.OtpUsecase) *AuthHandler {
	return &AuthHandler{
		UserSessionUC: userSessionUsecase,
		OtpUC:         otpUsecase,
		UserUC:        userUsecase,
	}
}

// Login godoc
// @Summary      User Login
// @Description  Authenticate user and return JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body      dto.LoginRequest  true  "Login Request"
// @Success      200           {object}  dto.LoginResponse
// @Failure      400           {object}  response.ErrorResponse
// @Failure      401           {object}  response.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorI18n(ctx, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// validate request
	if err := validator.ValidateStructWithContext(ctx, &req); err != nil {
		return response.ValidationErrorI18n(ctx, err)
	}

	loginResponse, err := h.UserSessionUC.Login(ctx.Context(), req.Username, req.Password, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		return response.ErrorI18n(ctx, fiber.StatusUnauthorized, "auth.error.invalid_credentials", nil)
	}

	return response.SuccessI18n(ctx, "app.success", loginResponse)
}

// Logout godoc
// @Summary      User Logout
// @Description  Invalidate user session
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.SuccessResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	// get user session from context
	userSession := ctx.Locals("user_token")
	if userSession == nil {
		return response.ErrorI18n(ctx, fiber.StatusUnauthorized, "auth.error.unauthorized", nil)
	}

	// revoke session
	err := h.UserSessionUC.Logout(ctx.Context(), userSession.(string))
	if err != nil {
		return response.ErrorI18n(ctx, fiber.StatusInternalServerError, "auth.error.logout_failed", nil)
	}

	return response.SuccessI18n(ctx, "auth.success", nil)
}

// ResetPassword godoc
// @Summary      Reset User Password
// @Description  Reset user password using OTP verification
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        resetPasswordRequest  body      dto.ResetPasswordRequest  true  "Reset Password Request"
// @Success      200                   {object}  dto.ResetPasswordResponse
// @Failure      400                   {object}  response.ErrorResponse
// @Router       /auth/password/reset [post]
func (h *AuthHandler) ResetPassword(ctx *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorI18n(ctx, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// validate request
	if err := validator.ValidateStructWithContext(ctx, &req); err != nil {
		return response.ValidationErrorI18n(ctx, err)
	}

	// reset user password
	err := h.UserUC.ResetPassword(ctx.Context(), req)
	if err != nil {
		return response.BadRequestI18n(ctx, err.Error(), nil)
	}

	var out dto.ResetPasswordResponse = dto.ResetPasswordResponse{
		Message: "auth.password_reset_success",
	}

	return response.SuccessI18n(ctx, "auth.success", out)
}

// Register godoc
// @Summary      User Registration
// @Description  Register a new user with OTP verification
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        registerRequest  body      dto.RegisterRequest  true  "Register Request"
// @Success      200              {object}  dto.LoginResponse
// @Failure      400              {object}  response.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorI18n(ctx, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// validate request
	if err := validator.ValidateStructWithContext(ctx, &req); err != nil {
		return response.ValidationErrorI18n(ctx, err)
	}

	req.FromIP = ctx.IP()
	req.UserAgent = ctx.Get("User-Agent")

	// register user
	out, err := h.UserUC.Register(ctx.Context(), req)
	if err != nil {
		return response.BadRequestI18n(ctx, err.Error(), nil)
	}

	return response.SuccessI18n(ctx, "auth.success", out)
}
