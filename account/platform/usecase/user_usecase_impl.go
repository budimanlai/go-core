package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-core/account/dto"
	"github.com/budimanlai/go-core/account/platform/security"
	"github.com/budimanlai/go-pkg/helpers"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type userUsecaseImpl struct {
	repo   repository.UserRepository
	hasher security.PasswordHasher
}

func NewUserUsecase(repo repository.UserRepository, hasher security.PasswordHasher) usecase.UserUsecase {
	return &userUsecaseImpl{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *userUsecaseImpl) Register(req *dto.RegisterRequest) (interface{}, error) {
	// Check if email exists
	existingUser, err := u.repo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Check if username exists
	existingUser, err = u.repo.FindByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Check if handphone exists
	existingUser, err = u.repo.FindByHandphone(req.Handphone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check handphone: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("handphone already registered")
	}

	// Hash password
	hashedPassword, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate auth key
	authKey := helpers.GenerateRandomString(16)

	// Generate verification token
	verificationToken := helpers.GenerateRandomString(32)

	now := time.Now()
	var user entity.User = entity.User{}
	if err := copier.Copy(&user, req); err != nil {
		return nil, fmt.Errorf("failed to copy request to user: %w", err)
	}
	user.AuthKey = authKey
	user.PasswordHash = hashedPassword
	user.Status = "active"
	user.LoginDashboard = "N"
	user.VerificationToken = &verificationToken
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := u.repo.Create(&user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return toUserResponse(&user), nil
}

func (u *userUsecaseImpl) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !u.hasher.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive() {
		return nil, errors.New("user account is not active")
	}

	// TODO: Generate JWT token here
	token := "jwt_token_placeholder"

	return &dto.LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (u *userUsecaseImpl) GetByID(id uint) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return toUserResponse(user), nil
}

func (u *userUsecaseImpl) Update(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := copier.CopyWithOption(user, req, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, fmt.Errorf("failed to copy request to user: %w", err)
	}

	user.UpdatedAt = time.Now()
	user.UpdatedBy = id

	if err := u.repo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return toUserResponse(user), nil
}

func (u *userUsecaseImpl) Delete(id uint) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.SoftDelete()
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) List(page, pageSize int) (*dto.ListUserResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, total, err := u.repo.List(pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	userResponses := make([]*dto.UserResponse, len(users))
	if err := copier.Copy(&userResponses, &users); err != nil {
		return nil, fmt.Errorf("failed to copy users to response: %w", err)
	}

	return &dto.ListUserResponse{
		Users:      userResponses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (u *userUsecaseImpl) Activate(id uint) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.Activate()
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) Deactivate(id uint) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.Deactivate()
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) Suspend(id uint) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.Suspend()
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) VerifyEmail(token string) error {
	user, err := u.repo.FindByVerificationToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid verification token")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.ClearVerificationToken()
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) EnableDashboard(id uint) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.EnableDashboardAccess()
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) DisableDashboard(id uint) error {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.DisableDashboardAccess()
	return u.repo.Update(user)
}

func toUserResponse(user *entity.User) *dto.UserResponse {
	var response dto.UserResponse
	copier.Copy(&response, user)
	return &response
}
