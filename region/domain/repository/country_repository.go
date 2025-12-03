package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
)

type CountryRepository interface {
	// GetByCode retrieves a country entity based on its code.
	GetByCode(code string) (*entity.Country, error)

	// GetAll retrieves all country entities.
	GetAll() ([]*entity.Country, error)
}
