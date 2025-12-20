package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/model"
)

type DistrictRepository interface {
	base.BaseRepository[entity.District, model.DistrictModel]
}

type districtRepositoryImpl struct {
	base.BaseRepository[entity.District, model.DistrictModel]
}

func NewDistrictRepository(f *base.Factory) DistrictRepository {
	return &districtRepositoryImpl{
		BaseRepository: base.NewRepository[entity.District, model.DistrictModel](f),
	}
}
