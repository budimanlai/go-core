package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/domain/usecase"
)

type subdistrictUsecaseImpl struct {
	repo repository.SubDistrictRepository
}

func NewSubDistrictUsecase(repo repository.SubDistrictRepository) usecase.SubDistrictUsecase {
	return &subdistrictUsecaseImpl{repo: repo}
}

func (u *subdistrictUsecaseImpl) GetAllByDistrict(districtID uint) ([]*entity.SubDistrict, error) {
	return u.repo.GetAllByDistrict(districtID)
}
