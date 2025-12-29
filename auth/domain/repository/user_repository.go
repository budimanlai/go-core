package repository

import (
	"github.com/budimanlai/go-core/base"

	entity "github.com/budimanlai/go-core/auth/domain/entity"
	model "github.com/budimanlai/go-core/auth/models"
)

type UserRepository interface {
	base.BaseRepository[entity.User, model.User]
}
