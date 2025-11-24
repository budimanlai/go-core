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
	"github.com/budimanlai/go-pkg/logger"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type CustomUserUsecase struct {
	usecase.UserUsecase
	repo   repository.UserRepository
	hasher security.PasswordHasher
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type CustomUserResponse struct {
	Profile dto.UserResponse `json:"profile"`
	Token   tokenResponse    `json:"token"`
}

func NewCustomUserUsecase(userUsecase usecase.UserUsecase, repo repository.UserRepository, hasher security.PasswordHasher) usecase.UserUsecase {
	obj := &CustomUserUsecase{
		UserUsecase: userUsecase,
		repo:        repo,
		hasher:      hasher,
	}
	obj.SetCustomToResponse(obj.ToResponse)

	return obj
}

func (u *CustomUserUsecase) Register1(req *dto.RegisterRequest) (interface{}, error) {
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
	authKey := helpers.GenerateRandomString(32)

	// Generate verification token
	verificationToken := helpers.GenerateRandomString(32)

	now := time.Now()
	var user entity.User
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

	return u.ToResponse(&user), nil
}

func (u *CustomUserUsecase) ToResponse(user *entity.User) interface{} {
	logger.Printf("Custom toUserResponse")
	var response CustomUserResponse = CustomUserResponse{}
	copier.Copy(&response.Profile, user)
	response.Token.AccessToken = "access-token"
	response.Token.RefreshToken = "refresh-token"
	return &response
}
