package usecase

import (
	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/auth/domain/repository"
	"github.com/budimanlai/go-core/auth/domain/usecase"
	"github.com/budimanlai/go-core/base"
	"gorm.io/gorm"
)

type UserUsecaseImpl struct {
	base.BaseUsecase[entity.User]
}

func NewUserUsecaseImpl(db *gorm.DB, repo repository.UserRepository) usecase.UserUsecase {
	return &UserUsecaseImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
	}
}
