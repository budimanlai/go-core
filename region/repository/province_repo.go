package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/model"
)

type ProvinceRepository interface {
	base.BaseRepository[entity.Province, model.ProvinceModel]
}

type provinceRepositoryImpl struct {
	base.BaseRepository[entity.Province, model.ProvinceModel]
}

func NewProvinceRepository(f *base.Factory) ProvinceRepository {
	return &provinceRepositoryImpl{
		BaseRepository: base.NewRepository[entity.Province, model.ProvinceModel](f),
	}
}
