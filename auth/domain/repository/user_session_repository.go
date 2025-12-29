package repository

import (
	entity "github.com/budimanlai/go-core/auth/domain/entity"
	model "github.com/budimanlai/go-core/auth/models"
	"github.com/budimanlai/go-core/base"
)

type UserSessionRepository interface {
	base.BaseRepository[entity.UserSession, model.UserSession]
}
