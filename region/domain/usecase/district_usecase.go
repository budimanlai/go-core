package usecase

import (
	"github.com/budimanlai/go-core/region/domain/entity"
)

type DistrictUsecase interface {
	// GetByID retrieves a district by its ID.
	GetByID(id uint) (*entity.District, error)

	// GetAllByCity retrieves all districts by city ID.
	GetAllByCity(city_id uint) ([]*entity.District, error)
}
