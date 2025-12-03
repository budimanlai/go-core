package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/models"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type districtRepositoryImpl struct {
	db *gorm.DB
}

func NewDistrictRepository(db *gorm.DB) repository.DistrictRepository {
	return &districtRepositoryImpl{db: db}
}

func (r *districtRepositoryImpl) GetByID(id uint) (*entity.District, error) {
	var model models.District
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var district entity.District
	if err := copier.Copy(&district, &model); err != nil {
		return nil, err
	}
	return &district, nil
}

func (r *districtRepositoryImpl) GetAllByCity(cityID uint) ([]*entity.District, error) {
	var modelList []models.District
	if err := r.db.Where("city_id = ?", cityID).Order("dis_name ASC").Find(&modelList).Error; err != nil {
		return nil, err
	}

	var districts []*entity.District
	if err := copier.Copy(&districts, &modelList); err != nil {
		return nil, err
	}
	return districts, nil
}
