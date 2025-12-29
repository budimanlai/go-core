package repository

import (
	entity "github.com/budimanlai/go-core/auth/domain/entity"
	repository "github.com/budimanlai/go-core/auth/domain/repository"
	model "github.com/budimanlai/go-core/auth/models"
	"github.com/budimanlai/go-core/base"
)

type userSessionRepositoryImpl struct {
	base.BaseRepository[entity.UserSession, model.UserSession]
}

func NewUserSessionRepository(f *base.Factory) repository.UserSessionRepository {
	return &userSessionRepositoryImpl{
		BaseRepository: base.NewRepository[entity.UserSession, model.UserSession](f),
	}
}
