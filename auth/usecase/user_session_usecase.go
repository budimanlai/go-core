package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/domain/repository"
	"github.com/budimanlai/go-core/auth/domain/usecase"
	"github.com/budimanlai/go-core/auth/dto"
	"github.com/budimanlai/go-core/auth/models"
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/middleware/auth"

	"gorm.io/gorm"

	pkg_helpers "github.com/budimanlai/go-pkg/helpers"
	pkg_security "github.com/budimanlai/go-pkg/security"
)

type UserSessionUsecaseImpl struct {
	base.BaseUsecase[entity.UserSession]

	UserRepository repository.UserRepository

	// MultipleLoginAllowed indicates whether multiple logins are allowed for a user
	MultipleLoginAllowed bool

	JWTService auth.JWTService
}

func NewUserSessionUsecaseImpl(db *gorm.DB, repo repository.UserSessionRepository,
	userRepo repository.UserRepository, jwtService auth.JWTService) usecase.UserSessionUsecase {
	return &UserSessionUsecaseImpl{
		BaseUsecase:          base.NewBaseUsecase(repo, db),
		UserRepository:       userRepo,
		MultipleLoginAllowed: false,
		JWTService:           jwtService,
	}
}

// IsMultipleLoginAllowed returns whether multiple logins are allowed for a user
func (u *UserSessionUsecaseImpl) IsMultipleLoginAllowed() bool {
	return u.MultipleLoginAllowed
}

// SetMultipleLoginAllowed sets whether multiple logins are allowed for a user
func (u *UserSessionUsecaseImpl) SetMultipleLoginAllowed(allowed bool) {
	u.MultipleLoginAllowed = allowed
}

// RevokeSessionsByUserID revokes all sessions for the given user ID
func (u *UserSessionUsecaseImpl) RevokeSessionsByUserID(ctx context.Context, userID uint) {
	// Revoke all sessions for the given user ID
	u.GetDB().Model(&models.UserSession{}).Where("user_id = ? AND remove_on IS NULL", userID).
		Update("remove_on", time.Now())
}

// GenerateSession creates a new user session for the given user ID
func (u *UserSessionUsecaseImpl) GenerateSession(ctx context.Context, userID uint, fromIP, userAgent string) (*entity.UserSession, error) {
	var out *entity.UserSession
	err := u.GetDB().Transaction(func(tx *gorm.DB) error {
		// check if multiple login is allowed
		if !u.IsMultipleLoginAllowed() {
			// revoke other sessions
			u.RevokeSessionsByUserID(ctx, userID)
		}

		// create new session
		sessionEntity := &entity.UserSession{
			UserID:       userID,
			AppID:        1,
			Tokens:       pkg_helpers.GenerateRandomString(32),
			FromIP:       fromIP,
			UserAgent:    userAgent,
			LastAccessOn: pkg_helpers.Pointer(time.Now()),
		}

		// save to db
		err := u.Create(ctx, sessionEntity)
		if err != nil {
			return err
		}
		out = sessionEntity

		return nil
	})

	return out, err
}

// Login authenticates a user and returns a dto.LoginResponse if successful
func (u *UserSessionUsecaseImpl) Login(ctx context.Context, username, password, fromIP, userAgent string) (*dto.LoginResponse, error) {
	// 1. check if username or password is not empty
	if username == "" || password == "" {
		return nil, errors.New("username or password can't blank")
	}

	// 2. find user by email or handphone
	user, err := u.UserRepository.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("email = ? or handphone = ? and status = ?", username, username, "active")
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 3. check if user.status is active
	if !user.IsActive() {
		return nil, errors.New("user is not active")
	}

	// 4. verify password
	if ok, err := pkg_security.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("invalid password")
	}

	// 5. Generate user session token
	sessionEntity, err := u.GenerateSession(ctx, user.ID, fromIP, userAgent)
	if err != nil {
		return nil, err
	}

	// 6. generate jwt token
	accessToken, err := u.JWTService.GenerateToken(sessionEntity.Tokens)
	if err != nil {
		return nil, err
	}

	// 7. prepare response
	var out dto.LoginResponse = dto.LoginResponse{
		UserID:    user.ID,
		Email:     user.Email,
		Handphone: user.Handphone,
		Fullname:  user.Fullname,
		Token: dto.Token{
			AccessToken: accessToken,
		},
	}

	return &out, nil
}

// Logout revokes the user session associated with the given token string
func (u *UserSessionUsecaseImpl) Logout(ctx context.Context, tokenString string) error {
	// revoke session
	return u.GetDB().Model(&models.UserSession{}).
		Where("tokens = ? AND remove_on IS NULL", tokenString).
		Update("remove_on", time.Now()).Error
}

// check if token is valid and return user session
func (u *UserSessionUsecaseImpl) VerifyToken(ctx context.Context, tokenString string) (*dto.LoginResponse, error) {
	// check if token isExists in user sessions
	result, err := u.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("tokens = ? AND remove_on IS NULL", tokenString)
	})
	if err != nil {
		return nil, err
	}

	// check if token is exists
	if errors.Is(err, gorm.ErrRecordNotFound) || result == nil {
		return nil, errors.New("invalid token")
	}

	// update last used on
	result.LastAccessOn = pkg_helpers.Pointer(time.Now())
	if err := u.Update(ctx, result); err != nil {
		return nil, err
	}

	// tokenString is valid, now get user info base on userID
	user, err := u.UserRepository.FindByID(ctx, result.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := u.JWTService.GenerateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// prepare response
	var out dto.LoginResponse = dto.LoginResponse{
		UserID:    user.ID,
		Email:     user.Email,
		Handphone: user.Handphone,
		Fullname:  user.Fullname,
		Token: dto.Token{
			AccessToken: accessToken,
		},
	}

	return &out, nil
}
