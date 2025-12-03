package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
)

type SubDistrictRepository interface {
	// GetByID retrieves a subdistrict by its ID.
	GetByID(id uint) (*entity.SubDistrict, error)

	// GetAllByDistrict retrieves all subdistricts by district ID.
	GetAllByDistrict(dis_id uint) ([]*entity.SubDistrict, error)
}
