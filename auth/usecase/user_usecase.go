package usecase

import (
	"context"
	"errors"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/domain/repository"
	"github.com/budimanlai/go-core/auth/domain/usecase"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/base"
	"gorm.io/gorm"

	pkg_security "github.com/budimanlai/go-pkg/security"
)

type UserUsecaseImpl struct {
	base.BaseUsecase[entity.User]

	OtpUC usecase.OtpUsecase
}

func NewUserUsecaseImpl(db *gorm.DB, repo repository.UserRepository, otpUC usecase.OtpUsecase) usecase.UserUsecase {
	return &UserUsecaseImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
		OtpUC:       otpUC,
	}
}

// ResetPassword resets the user's password after verifying the OTP.
func (u *UserUsecaseImpl) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error {
	// 1. Check if otp is valid
	valid, err := u.OtpUC.Status(ctx, request.Identifier, request.TrxID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("Invalid OTP or expired")
	}

	// 2. Find user by identifier
	user, err := u.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		if request.Channel == "email" {
			return d.Where("email = ?", request.Identifier)
		} else {
			return d.Where("handphone = ?", request.Identifier)
		}
	})
	if user == nil {
		return errors.New("User not found")
	}

	// 3. Update user password
	user.PasswordHash = pkg_security.HashPassword(request.Password)
	err = u.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
