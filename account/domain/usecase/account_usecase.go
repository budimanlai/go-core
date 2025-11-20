package usecase

import (
	"context"
	"errors"

	"github.com/budimanlai/go-core/account/domain/entity"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAccountInactive      = errors.New("account is inactive")
)

// PasswordHasher interface for dependency injection
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

// AccountUsecase defines business logic operations for account management
type AccountUsecase interface {
	Register(ctx context.Context, email, username, password, fullName string) (*entity.Account, error)
	Login(ctx context.Context, identifier, password string) (*entity.Account, error)
	GetByID(ctx context.Context, id string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Account, int64, error)
	Activate(ctx context.Context, id string) error
	Deactivate(ctx context.Context, id string) error
}
