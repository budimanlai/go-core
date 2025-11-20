package repository

import (
	"context"

	"github.com/budimanlai/go-core/account/domain/entity"
)

// AccountRepository defines the interface for account data operations
// This is part of the domain layer and should not depend on any specific implementation
type AccountRepository interface {
	// Create creates a new account
	Create(ctx context.Context, account *entity.Account) error

	// FindByID finds an account by ID
	FindByID(ctx context.Context, id string) (*entity.Account, error)

	// FindByEmail finds an account by email
	FindByEmail(ctx context.Context, email string) (*entity.Account, error)

	// FindByUsername finds an account by username
	FindByUsername(ctx context.Context, username string) (*entity.Account, error)

	// Update updates an existing account
	Update(ctx context.Context, account *entity.Account) error

	// Delete hard deletes an account
	Delete(ctx context.Context, id string) error

	// SoftDelete soft deletes an account
	SoftDelete(ctx context.Context, id string) error

	// List lists accounts with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.Account, error)

	// Count counts total accounts
	Count(ctx context.Context) (int64, error)
}
