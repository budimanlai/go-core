package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-core/account/dto"
	"github.com/budimanlai/go-core/account/platform/security"
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
	return &CustomUserUsecase{
		UserUsecase: userUsecase,
		repo:        repo,
		hasher:      hasher,
	}
}

func (u *CustomUserUsecase) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
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
	authKey, err := generateAuthKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth key: %w", err)
	}

	// Generate verification token
	verificationToken, err := generateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

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

	return toUserResponse(&user), nil
}

func toUserResponse(user *entity.User) *dto.UserResponse {
	var response dto.UserResponse
	copier.Copy(&response, user)
	return &response
}

func generateAuthKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
