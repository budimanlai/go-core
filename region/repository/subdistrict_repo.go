package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/model"
)

type SubdistrictRepository interface {
	base.BaseRepository[entity.Subdistrict, model.SubdistrictModel]
}

type subdistrictRepositoryImpl struct {
	base.BaseRepository[entity.Subdistrict, model.SubdistrictModel]
}

func NewSubdistrictRepository(f *base.Factory) SubdistrictRepository {
	return &subdistrictRepositoryImpl{
		BaseRepository: base.NewRepository[entity.Subdistrict, model.SubdistrictModel](f),
	}
}
