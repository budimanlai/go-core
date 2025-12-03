package usecase

import "github.com/budimanlai/go-core/region/domain/entity"

type SubDistrictUsecase interface {
	// GetByDistrictID retrieves a list of sub-districts by district ID
	GetAllByDistrict(district_id uint) ([]*entity.SubDistrict, error)
}
