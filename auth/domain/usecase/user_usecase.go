package usecase

import (
	"github.com/budimanlai/go-core/auth/domain/entity"
	"github.com/budimanlai/go-core/base"
)

type UserUsecase interface {
	base.BaseUsecase[entity.User]
}
