package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/domain/repository"
	"github.com/budimanlai/go-core/auth/domain/usecase"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/base"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	pkg_helpers "github.com/budimanlai/go-pkg/helpers"
	pkg_security "github.com/budimanlai/go-pkg/security"
)

type UserUsecaseImpl struct {
	base.BaseUsecase[entity.User]

	OtpUC         usecase.OtpUsecase
	UserSessionUC usecase.UserSessionUsecase
}

func NewUserUsecaseImpl(db *gorm.DB, repo repository.UserRepository, otpUC usecase.OtpUsecase, userSessionUC usecase.UserSessionUsecase) usecase.UserUsecase {
	return &UserUsecaseImpl{
		BaseUsecase:   base.NewBaseUsecase(repo, db),
		OtpUC:         otpUC,
		UserSessionUC: userSessionUC,
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

	// 4. Revoke OTP
	u.OtpUC.Revoke(ctx, request.Identifier, request.TrxID)

	return nil
}

// Register registers a new user.
func (u *UserUsecaseImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// 1. Check if OTP is valid
	valid, err := u.OtpUC.Status(ctx, req.Email, req.TrxID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("Invalid OTP or expired")
	}

	// 2. Check if email is already registered
	existingUser, err := u.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("email = ?", req.Email)
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("Email is already registered")
	}

	// 3. check if handphone is already registered
	existingUser, err = u.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("handphone = ?", req.Handphone)
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("Handphone is already registered")
	}

	// 4. Create new user
	var out dto.LoginResponse
	err = u.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. create user entity
		newUser := &entity.User{
			Username:     req.Email,
			Email:        req.Email,
			Handphone:    req.Handphone,
			AuthKey:      pkg_helpers.GenerateRandomString(32),
			PasswordHash: pkg_security.HashPassword(req.Password),
			Fullname:     req.Fullname,
			Status:       "active",
			CreatedBy:    1,
			UpdatedBy:    1,
			UpdatedAt:    time.Now(),
			CreatedAt:    time.Now(),
		}
		err = u.Create(ctx, newUser)
		if err != nil {
			return err
		}

		// 2. generate jwt token
		accessToken, err := u.UserSessionUC.GenerateToken(ctx, newUser.ID, req.FromIP, req.UserAgent)
		if err != nil {
			return err
		}

		// 3. prepare output
		copier.Copy(&out, &req)
		out.UserID = newUser.ID
		out.Token = dto.Token{
			AccessToken: accessToken,
		}

		// 4. revoke otp
		if req.Channel == "email" {
			u.OtpUC.Revoke(ctx, req.Email, req.TrxID)
		} else {
			u.OtpUC.Revoke(ctx, req.Handphone, req.TrxID)
		}

		return nil
	})

	return &out, err
}
