package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/models"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type provinceRepositoryImpl struct {
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) repository.ProvinceRepository {
	return &provinceRepositoryImpl{db: db}
}

func (r *provinceRepositoryImpl) GetAll() ([]*entity.Province, error) {
	var modelList []models.Province
	if err := r.db.Where("status = ?", 1).Order("prov_name ASC").Find(&modelList).Error; err != nil {
		return nil, err
	}

	var provinces []*entity.Province
	if err := copier.Copy(&provinces, &modelList); err != nil {
		return nil, err
	}
	return provinces, nil
}
