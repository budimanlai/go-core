package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/domain/usecase"
)

type countryUsecaseImpl struct {
	repo repository.CountryRepository
}

func NewCountryUsecase(repo repository.CountryRepository) usecase.CountryUsecase {
	return &countryUsecaseImpl{repo: repo}
}

func (u *countryUsecaseImpl) GetByCode(code string) (*entity.Country, error) {
	return u.repo.GetByCode(code)
}

func (u *countryUsecaseImpl) GetAll() ([]*entity.Country, error) {
	return u.repo.GetAll()
}
