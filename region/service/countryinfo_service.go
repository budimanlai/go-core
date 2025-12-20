package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type CountryinfoService interface {
	base.BaseService[entity.Countryinfo]
}

type countryinfoServiceImpl struct {
	base.BaseService[entity.Countryinfo]
	repo repository.CountryinfoRepository
}

func NewCountryinfoService(repo repository.CountryinfoRepository, db *gorm.DB) CountryinfoService {
	return &countryinfoServiceImpl{
		BaseService: base.NewBaseService(repo, db),
		repo:        repo,
	}
}
