package usecase

import "github.com/budimanlai/go-core/region/domain/entity"

type ProvinceUsecase interface {
	// GetAll retrieves a list of all provinces
	GetAll() ([]*entity.Province, error)
}
