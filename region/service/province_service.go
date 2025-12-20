package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type ProvinceService interface {
	base.BaseService[entity.Province]
}

type provinceServiceImpl struct {
	base.BaseService[entity.Province]
	repo repository.ProvinceRepository
}

func NewProvinceService(repo repository.ProvinceRepository, db *gorm.DB) ProvinceService {
	return &provinceServiceImpl{
		BaseService: base.NewBaseService(repo, db),
		repo:        repo,
	}
}
