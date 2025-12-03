package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/models"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type countryRepositoryImpl struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) repository.CountryRepository {
	return &countryRepositoryImpl{db: db}
}

func (r *countryRepositoryImpl) GetByCode(code string) (*entity.Country, error) {
	var model models.Country
	if err := r.db.Where("iso_alpha2 = ? OR iso_alpha3 = ?", code, code).First(&model).Error; err != nil {
		return nil, err
	}

	var country entity.Country
	if err := copier.Copy(&country, &model); err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *countryRepositoryImpl) GetAll() ([]*entity.Country, error) {
	var modelList []models.Country
	if err := r.db.Where("status = ?", "active").Find(&modelList).Error; err != nil {
		return nil, err
	}

	var countries []*entity.Country
	if err := copier.Copy(&countries, &modelList); err != nil {
		return nil, err
	}
	return countries, nil
}
