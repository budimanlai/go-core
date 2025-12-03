package usecase

import "github.com/budimanlai/go-core/region/domain/entity"

type CountryUsecase interface {
	// GetByCode retrieves a country by its code.
	GetByCode(code string) (*entity.Country, error)

	// GetAll retrieves all countries.
	GetAll() ([]*entity.Country, error)
}
