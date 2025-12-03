package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/domain/usecase"
)

type provinceUsecaseImpl struct {
	repo repository.ProvinceRepository
}

func NewProvinceUsecase(repo repository.ProvinceRepository) usecase.ProvinceUsecase {
	return &provinceUsecaseImpl{repo: repo}
}

func (u *provinceUsecaseImpl) GetAll() ([]*entity.Province, error) {
	provinces, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return provinces, nil
}
