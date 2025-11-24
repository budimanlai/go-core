package repository

import (
	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/models"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var model models.User
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByUsername(username string) (*entity.User, error) {
	var model models.User
	if err := r.db.Where("username = ?", username).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByHandphone(handphone string) (*entity.User, error) {
	var model models.User
	if err := r.db.Where("handphone = ?", handphone).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByID(id uint) (*entity.User, error) {
	var model models.User
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByVerificationToken(token string) (*entity.User, error) {
	var model models.User
	if err := r.db.Where("verification_token = ?", token).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) Create(user *entity.User) error {
	model := toModel(user)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	return nil
}

func (r *userRepositoryImpl) Update(user *entity.User) error {
	model := toModel(user)
	return r.db.Save(&model).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepositoryImpl) List(limit, offset int) ([]*entity.User, int64, error) {
	var modelList []models.User
	var total int64

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&modelList).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(modelList))
	if err := copier.Copy(&users, &modelList); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func toEntity(model *models.User) *entity.User {
	var user entity.User
	copier.Copy(&user, model)
	return &user
}

func toModel(user *entity.User) *models.User {
	var model models.User
	copier.Copy(&model, user)
	return &model
}
