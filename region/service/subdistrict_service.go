package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type SubdistrictService interface {
	base.BaseService[entity.Subdistrict]
}

type subdistrictServiceImpl struct {
	base.BaseService[entity.Subdistrict]
	repo repository.SubdistrictRepository
}

func NewSubdistrictService(repo repository.SubdistrictRepository, db *gorm.DB) SubdistrictService {
	return &subdistrictServiceImpl{
		BaseService: base.NewBaseService(repo, db),
		repo:        repo,
	}
}
