package repository

import (
	"context"

	"github.com/budimanlai/go-core/account/domain/entity"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByHandphone(ctx context.Context, handphone string) (*entity.User, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, int64, error)
}
