package repository

import (
	"github.com/budimanlai/go-core/account/domain/entity"
)

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByHandphone(handphone string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	FindByVerificationToken(token string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
	List(limit, offset int) ([]*entity.User, int64, error)
}
