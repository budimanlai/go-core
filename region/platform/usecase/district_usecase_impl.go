package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/domain/usecase"
)

type districtUsecaseImpl struct {
	repo repository.DistrictRepository
}

func NewDistrictUsecase(repo repository.DistrictRepository) usecase.DistrictUsecase {
	return &districtUsecaseImpl{repo: repo}
}

func (u *districtUsecaseImpl) GetByID(id uint) (*entity.District, error) {
	return u.repo.GetByID(id)
}

func (u *districtUsecaseImpl) GetAllByCity(cityID uint) ([]*entity.District, error) {
	return u.repo.GetAllByCity(cityID)
}
