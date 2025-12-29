package usecase

import (
	"context"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/base"
)

type UserUsecase interface {
	base.BaseUsecase[entity.User]

	// ResetPassword resets the user's password using OTP verification
	ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error
}
