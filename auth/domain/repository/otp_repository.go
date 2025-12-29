package repository

import (
	entity "github.com/budimanlai/go-core/auth/domain/entity"
	model "github.com/budimanlai/go-core/auth/models"
	"github.com/budimanlai/go-core/base"
)

type OtpRepository interface {
	base.BaseRepository[entity.Otp, model.Otp]
}
