package persistence

import (
	"context"

	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/models"
	"gorm.io/gorm"
)

type accountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &accountRepositoryImpl{
		db: db,
	}
}

func toEntity(model *models.AccountModel) *entity.Account {
	if model == nil {
		return nil
	}
	return &entity.Account{
		ID:        model.ID,
		Email:     model.Email,
		Username:  model.Username,
		Password:  model.Password,
		FullName:  model.FullName,
		Role:      model.Role,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
	}
}

func toModel(account *entity.Account) *models.AccountModel {
	if account == nil {
		return nil
	}
	return &models.AccountModel{
		ID:        account.ID,
		Email:     account.Email,
		Username:  account.Username,
		Password:  account.Password,
		FullName:  account.FullName,
		Role:      account.Role,
		IsActive:  account.IsActive,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		DeletedAt: account.DeletedAt,
	}
}

func (r *accountRepositoryImpl) Create(ctx context.Context, account *entity.Account) error {
	model := toModel(account)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return result.Error
	}
	account.ID = model.ID
	return nil
}

func (r *accountRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Account, error) {
	var model models.AccountModel
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toEntity(&model), nil
}

func (r *accountRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.Account, error) {
	var model models.AccountModel
	result := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toEntity(&model), nil
}

func (r *accountRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.Account, error) {
	var model models.AccountModel
	result := r.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toEntity(&model), nil
}

func (r *accountRepositoryImpl) Update(ctx context.Context, account *entity.Account) error {
	model := toModel(account)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *accountRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&models.AccountModel{}, "id = ?", id).Error
}

func (r *accountRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AccountModel{}, "id = ?", id).Error
}

func (r *accountRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Account, error) {
	var models []models.AccountModel
	result := r.db.WithContext(ctx).Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	accounts := make([]*entity.Account, len(models))
	for i, model := range models {
		accounts[i] = toEntity(&model)
	}
	return accounts, nil
}

func (r *accountRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.AccountModel{}).Where("deleted_at IS NULL").Count(&count)
	return count, result.Error
}
