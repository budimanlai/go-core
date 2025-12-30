package usecase

import (
	"context"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/base"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserSessionUsecase interface {
	base.BaseUsecase[entity.UserSession]

	// IsMultipleLoginAllowed indicates whether multiple logins are allowed for a user
	IsMultipleLoginAllowed() bool

	// SetMultipleLoginAllowed sets whether multiple logins are allowed for a user
	SetMultipleLoginAllowed(allowed bool)

	// RevokeSessionsByUserID revokes all sessions for a given user ID
	RevokeSessionsByUserID(ctx context.Context, userID uint)

	// GenerateSession generates a new user session
	GenerateSession(ctx context.Context, userID uint, fromIP, userAgent string) (*entity.UserSession, error)

	// Login authenticates a user and returns a token if successful
	Login(ctx context.Context, username, password, fromIP, userAgent string) (*dto.LoginResponse, error)

	// Logout logs out a user by revoking their session token
	Logout(ctx context.Context, tokenString string) error

	// VerifyToken verifies a JWT token and returns the associated login response
	VerifyToken(ctx context.Context, tokenString string) (*dto.LoginResponse, error)

	// GetUserIDByToken retrieves the user ID associated with the given token
	GetUserIDByToken(ctx context.Context, tokenString string) (*entity.UserSession, error)

	// SuccessHandler handles successful JWT authentication
	SuccessHandler(c *fiber.Ctx, claims jwt.MapClaims) error

	// GenerateToken creates a new user session and generates a JWT token for the given user ID
	GenerateToken(ctx context.Context, user_id uint, fromIP, userAgent string) (string, error)
}
