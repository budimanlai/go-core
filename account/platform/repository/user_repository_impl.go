package repository

import (
	"context"
	"time"

	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/models"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByHandphone(ctx context.Context, handphone string) (*entity.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).Where("handphone = ?", handphone).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) FindByVerificationToken(ctx context.Context, token string) (*entity.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).Where("verification_token = ?", token).First(&model).Error; err != nil {
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	model := toModel(user)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	return nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	model := toModel(user)
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.User, int64, error) {
	var modelList []models.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&modelList).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(modelList))
	for i, model := range modelList {
		users[i] = toEntity(&model)
	}

	return users, total, nil
}

func toEntity(model *models.User) *entity.User {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.User{
		ID:                 model.ID,
		Username:           model.Username,
		AuthKey:            model.AuthKey,
		PasswordHash:       model.PasswordHash,
		PasswordResetToken: model.PasswordResetToken,
		Email:              model.Email,
		Fullname:           model.Fullname,
		Handphone:          model.Handphone,
		Dob:                model.Dob,
		Gender:             model.Gender,
		Status:             model.Status,
		MainRole:           model.MainRole,
		LoginDashboard:     model.LoginDashboard,
		Avatar:             model.Avatar,
		Address:            model.Address,
		Zipcode:            model.Zipcode,
		DistrictID:         model.DistrictID,
		SubdistrictID:      model.SubdistrictID,
		CityID:             model.CityID,
		ProvinceID:         model.ProvinceID,
		CountryID:          model.CountryID,
		CreatedAt:          model.CreatedAt,
		CreatedBy:          model.CreatedBy,
		UpdatedAt:          model.UpdatedAt,
		UpdatedBy:          model.UpdatedBy,
		VerificationToken:  model.VerificationToken,
		DeletedAt:          deletedAt,
	}
}

func toModel(user *entity.User) *models.User {
	var deletedAt gorm.DeletedAt
	if user.DeletedAt != nil {
		deletedAt.Time = *user.DeletedAt
		deletedAt.Valid = true
	}

	return &models.User{
		ID:                 user.ID,
		Username:           user.Username,
		AuthKey:            user.AuthKey,
		PasswordHash:       user.PasswordHash,
		PasswordResetToken: user.PasswordResetToken,
		Email:              user.Email,
		Fullname:           user.Fullname,
		Handphone:          user.Handphone,
		Dob:                user.Dob,
		Gender:             user.Gender,
		Status:             user.Status,
		MainRole:           user.MainRole,
		LoginDashboard:     user.LoginDashboard,
		Avatar:             user.Avatar,
		Address:            user.Address,
		Zipcode:            user.Zipcode,
		DistrictID:         user.DistrictID,
		SubdistrictID:      user.SubdistrictID,
		CityID:             user.CityID,
		ProvinceID:         user.ProvinceID,
		CountryID:          user.CountryID,
		CreatedAt:          user.CreatedAt,
		CreatedBy:          user.CreatedBy,
		UpdatedAt:          user.UpdatedAt,
		UpdatedBy:          user.UpdatedBy,
		VerificationToken:  user.VerificationToken,
		DeletedAt:          deletedAt,
	}
}
