package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type ProvinceService interface {
	base.BaseUsecase[entity.Province]
}

type provinceServiceImpl struct {
	base.BaseUsecase[entity.Province]
	repo repository.ProvinceRepository
}

func NewProvinceService(repo repository.ProvinceRepository, db *gorm.DB) ProvinceService {
	return &provinceServiceImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
		repo:        repo,
	}
}
