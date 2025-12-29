package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type SubdistrictService interface {
	base.BaseUsecase[entity.Subdistrict]
}

type subdistrictServiceImpl struct {
	base.BaseUsecase[entity.Subdistrict]
	repo repository.SubdistrictRepository
}

func NewSubdistrictService(repo repository.SubdistrictRepository, db *gorm.DB) SubdistrictService {
	return &subdistrictServiceImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
		repo:        repo,
	}
}
