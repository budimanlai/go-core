package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
)

type CityUsecase interface {
	// GetByID retrieves a city by its ID.
	GetByID(id uint) (*entity.City, error)

	// GetAllByProvince retrieves cities by their province ID.
	GetAllByProvince(prov_id uint) ([]*entity.City, error)
}
