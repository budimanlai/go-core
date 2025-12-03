package repository

import "github.com/budimanlai/go-core/region/domain/entity"

type ProvinceRepository interface {
	// GetAll retrieves all provinces.
	GetAll() ([]*entity.Province, error)
}
