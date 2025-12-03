package repository

import (
	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/domain/repository"
	"github.com/budimanlai/go-core/region/models"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type subdistrictRepositoryImpl struct {
	db *gorm.DB
}

func NewSubDistrictRepository(db *gorm.DB) repository.SubDistrictRepository {
	return &subdistrictRepositoryImpl{db: db}
}

func (r *subdistrictRepositoryImpl) GetByID(id uint) (*entity.SubDistrict, error) {
	var model models.SubDistrict
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var subdistrict entity.SubDistrict
	if err := copier.Copy(&subdistrict, &model); err != nil {
		return nil, err
	}
	return &subdistrict, nil
}

func (r *subdistrictRepositoryImpl) GetAllByDistrict(disID uint) ([]*entity.SubDistrict, error) {
	var modelList []models.SubDistrict
	if err := r.db.Where("dis_id = ?", disID).Order("subdis_name ASC").Find(&modelList).Error; err != nil {
		return nil, err
	}

	var subdistricts []*entity.SubDistrict
	if err := copier.Copy(&subdistricts, &modelList); err != nil {
		return nil, err
	}
	return subdistricts, nil
}
