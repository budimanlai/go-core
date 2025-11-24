package usecase

import (
	"context"
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

func (u *userUsecaseImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
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
	user := &entity.User{
		Username:          req.Username,
		AuthKey:           authKey,
		PasswordHash:      hashedPassword,
		Email:             req.Email,
		Fullname:          req.Fullname,
		Handphone:         req.Handphone,
		Dob:               req.Dob,
		Gender:            req.Gender,
		Status:            "active",
		LoginDashboard:    "N",
		Zipcode:           req.Zipcode,
		DistrictID:        req.DistrictID,
		SubdistrictID:     req.SubdistrictID,
		CityID:            req.CityID,
		ProvinceID:        req.ProvinceID,
		CountryID:         req.CountryID,
		Address:           req.Address,
		VerificationToken: &verificationToken,
		CreatedAt:         now,
		CreatedBy:         0,
		UpdatedAt:         now,
		UpdatedBy:         0,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return toUserResponse(user), nil
}

func (u *userUsecaseImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
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

func (u *userUsecaseImpl) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return toUserResponse(user), nil
}

func (u *userUsecaseImpl) Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}
	if req.Handphone != nil {
		user.Handphone = *req.Handphone
	}
	if req.Dob != nil {
		user.Dob = req.Dob
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.Zipcode != nil {
		user.Zipcode = *req.Zipcode
	}
	if req.DistrictID != nil {
		user.DistrictID = *req.DistrictID
	}
	if req.SubdistrictID != nil {
		user.SubdistrictID = *req.SubdistrictID
	}
	if req.CityID != nil {
		user.CityID = *req.CityID
	}
	if req.ProvinceID != nil {
		user.ProvinceID = *req.ProvinceID
	}
	if req.CountryID != nil {
		user.CountryID = *req.CountryID
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	user.UpdatedAt = time.Now()
	user.UpdatedBy = id

	if err := u.repo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return toUserResponse(user), nil
}

func (u *userUsecaseImpl) Delete(ctx context.Context, id uint) error {
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

func (u *userUsecaseImpl) List(ctx context.Context, page, pageSize int) (*dto.ListUserResponse, error) {
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
	for i, user := range users {
		userResponses[i] = toUserResponse(user)
	}

	return &dto.ListUserResponse{
		Users:      userResponses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (u *userUsecaseImpl) Activate(ctx context.Context, id uint) error {
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

func (u *userUsecaseImpl) Deactivate(ctx context.Context, id uint) error {
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

func (u *userUsecaseImpl) Suspend(ctx context.Context, id uint) error {
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

func (u *userUsecaseImpl) VerifyEmail(ctx context.Context, token string) error {
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

func (u *userUsecaseImpl) EnableDashboard(ctx context.Context, id uint) error {
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

func (u *userUsecaseImpl) DisableDashboard(ctx context.Context, id uint) error {
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
	return &dto.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Fullname:       user.Fullname,
		Handphone:      user.Handphone,
		Dob:            user.Dob,
		Gender:         user.Gender,
		Status:         user.Status,
		MainRole:       user.MainRole,
		LoginDashboard: user.LoginDashboard,
		Avatar:         user.Avatar,
		Address:        user.Address,
		Zipcode:        user.Zipcode,
		DistrictID:     user.DistrictID,
		SubdistrictID:  user.SubdistrictID,
		CityID:         user.CityID,
		ProvinceID:     user.ProvinceID,
		CountryID:      user.CountryID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
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
