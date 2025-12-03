package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/models"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type cityRepositoryImpl struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) repository.CityRepository {
	return &cityRepositoryImpl{db: db}
}

func (r *cityRepositoryImpl) GetByID(id uint) (*entity.City, error) {
	var model models.City
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var city entity.City
	if err := copier.Copy(&city, &model); err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *cityRepositoryImpl) GetAllByProvince(provID uint) ([]*entity.City, error) {
	var modelList []models.City
	if err := r.db.Where("prov_id = ?", provID).Order("city_name ASC").Find(&modelList).Error; err != nil {
		return nil, err
	}

	var cities []*entity.City
	if err := copier.Copy(&cities, &modelList); err != nil {
		return nil, err
	}
	return cities, nil
}
