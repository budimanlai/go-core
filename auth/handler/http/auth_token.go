package http

import (
	"github.com/budimanlai/go-pkg/response"
	"github.com/gofiber/fiber/v2"
)

// VerifyToken godoc
// @Summary      Verify Token
// @Description  Verify JWT token validity
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.LoginResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /auth/verify [get]
func (h *AuthHandler) VerifyToken(ctx *fiber.Ctx) error {
	// get user session from context
	userSession := ctx.Locals("user_token")
	if userSession == nil {
		return response.ErrorI18n(ctx, fiber.StatusUnauthorized, "auth.error.unauthorized", nil)
	}

	// revoke session
	loginResponse, err := h.UserSessionUC.VerifyToken(ctx.Context(), userSession.(string))
	if err != nil {
		return response.ErrorI18n(ctx, fiber.StatusUnauthorized, "auth.error.invalid_credentials", nil)
	}

	return response.SuccessI18n(ctx, "auth.success.token_valid", loginResponse)
}

// RefreshToken godoc
// @Summary      Refresh Token
// @Description  Refresh JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.LoginResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	return response.BadRequestI18n(ctx, "app.error.not_implemented", nil)
}
