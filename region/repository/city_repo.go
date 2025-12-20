package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/model"
)

type CityRepository interface {
	base.BaseRepository[entity.City, model.CityModel]
}

type cityRepositoryImpl struct {
	base.BaseRepository[entity.City, model.CityModel]
}

func NewCityRepository(f *base.Factory) CityRepository {
	return &cityRepositoryImpl{
		BaseRepository: base.NewRepository[entity.City, model.CityModel](f),
	}
}
