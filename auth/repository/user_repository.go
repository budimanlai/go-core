package repository

import (
	"github.com/budimanlai/go-core/base"

	entity "github.com/budimanlai/go-core/auth/domain/entity"
	repository "github.com/budimanlai/go-core/auth/domain/repository"
	model "github.com/budimanlai/go-core/auth/models"
)

type userRepositoryImpl struct {
	base.BaseRepository[entity.User, model.User]
}

func NewUserRepositoryImpl(f *base.Factory) repository.UserRepository {
	return &userRepositoryImpl{
		BaseRepository: base.NewRepository[entity.User, model.User](f),
	}
}
