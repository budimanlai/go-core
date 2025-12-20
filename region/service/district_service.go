package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type DistrictService interface {
	base.BaseService[entity.District]
}

type districtServiceImpl struct {
	base.BaseService[entity.District]
	repo repository.DistrictRepository
}

func NewDistrictService(repo repository.DistrictRepository, db *gorm.DB) DistrictService {
	return &districtServiceImpl{
		BaseService: base.NewBaseService(repo, db),
		repo:        repo,
	}
}
