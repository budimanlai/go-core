package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/domain/usecase"
)

type cityUsecaseImpl struct {
	repo repository.CityRepository
}

func NewCityUsecase(repo repository.CityRepository) usecase.CityUsecase {
	return &cityUsecaseImpl{repo: repo}
}

func (u *cityUsecaseImpl) GetByID(id uint) (*entity.City, error) {
	return u.repo.GetByID(id)
}

func (u *cityUsecaseImpl) GetAllByProvince(provID uint) ([]*entity.City, error) {
	return u.repo.GetAllByProvince(provID)
}
