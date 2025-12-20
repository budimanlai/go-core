package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type CityService interface {
	base.BaseService[entity.City]
}

type cityServiceImpl struct {
	base.BaseService[entity.City]
	repo repository.CityRepository
}

func NewCityService(repo repository.CityRepository, db *gorm.DB) CityService {
	return &cityServiceImpl{
		BaseService: base.NewBaseService(repo, db),
		repo:        repo,
	}
}
