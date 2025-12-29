package service

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/repository"
	"gorm.io/gorm"
)

type CountryinfoService interface {
	base.BaseUsecase[entity.Countryinfo]
}

type countryinfoServiceImpl struct {
	base.BaseUsecase[entity.Countryinfo]
	repo repository.CountryinfoRepository
}

func NewCountryinfoService(repo repository.CountryinfoRepository, db *gorm.DB) CountryinfoService {
	return &countryinfoServiceImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
		repo:        repo,
	}
}
