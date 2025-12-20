package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/model"
)

type CountryinfoRepository interface {
	base.BaseRepository[entity.Countryinfo, model.CountryinfoModel]
}

type countryinfoRepositoryImpl struct {
	base.BaseRepository[entity.Countryinfo, model.CountryinfoModel]
}

func NewCountryinfoRepository(f *base.Factory) CountryinfoRepository {
	return &countryinfoRepositoryImpl{
		BaseRepository: base.NewRepository[entity.Countryinfo, model.CountryinfoModel](f),
	}
}
