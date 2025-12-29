package repository

import (
	entity "github.com/budimanlai/go-core/auth/domain/entity"
	repository "github.com/budimanlai/go-core/auth/domain/repository"
	model "github.com/budimanlai/go-core/auth/models"
	"github.com/budimanlai/go-core/base"
)

type OtpRepositoryImpl struct {
	base.BaseRepository[entity.Otp, model.Otp]
}

func NewOtpRepositoryImpl(f *base.Factory) repository.OtpRepository {
	return &OtpRepositoryImpl{
		BaseRepository: base.NewRepository[entity.Otp, model.Otp](f),
	}
}
