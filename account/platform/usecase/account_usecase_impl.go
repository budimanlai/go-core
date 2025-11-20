package usecase

import (
	"context"
	"time"

	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/domain/usecase"
)

type accountUsecaseImpl struct {
	repo   repository.AccountRepository
	hasher usecase.PasswordHasher
}

// NewAccountUsecase creates a new account usecase implementation
func NewAccountUsecase(repo repository.AccountRepository, hasher usecase.PasswordHasher) usecase.AccountUsecase {
	return &accountUsecaseImpl{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *accountUsecaseImpl) Register(ctx context.Context, email, username, password, fullName string) (*entity.Account, error) {
	existing, _ := u.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, usecase.ErrAccountAlreadyExists
	}

	existing, _ = u.repo.FindByUsername(ctx, username)
	if existing != nil {
		return nil, usecase.ErrAccountAlreadyExists
	}

	hashedPassword, err := u.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	account := &entity.Account{
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		FullName:  fullName,
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (u *accountUsecaseImpl) Login(ctx context.Context, identifier, password string) (*entity.Account, error) {
	account, err := u.repo.FindByEmail(ctx, identifier)
	if err != nil || account == nil {
		account, err = u.repo.FindByUsername(ctx, identifier)
		if err != nil || account == nil {
			return nil, usecase.ErrInvalidCredentials
		}
	}

	if !account.IsActive {
		return nil, usecase.ErrAccountInactive
	}

	if !u.hasher.Verify(account.Password, password) {
		return nil, usecase.ErrInvalidCredentials
	}

	return account, nil
}

func (u *accountUsecaseImpl) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, usecase.ErrAccountNotFound
	}
	return account, nil
}

func (u *accountUsecaseImpl) Update(ctx context.Context, account *entity.Account) error {
	existing, err := u.repo.FindByID(ctx, account.ID)
	if err != nil || existing == nil {
		return usecase.ErrAccountNotFound
	}

	account.UpdatedAt = time.Now()
	return u.repo.Update(ctx, account)
}

func (u *accountUsecaseImpl) Delete(ctx context.Context, id string) error {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil || account == nil {
		return usecase.ErrAccountNotFound
	}

	return u.repo.SoftDelete(ctx, id)
}

func (u *accountUsecaseImpl) List(ctx context.Context, limit, offset int) ([]*entity.Account, int64, error) {
	accounts, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (u *accountUsecaseImpl) Activate(ctx context.Context, id string) error {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil || account == nil {
		return usecase.ErrAccountNotFound
	}

	account.Activate()
	return u.repo.Update(ctx, account)
}

func (u *accountUsecaseImpl) Deactivate(ctx context.Context, id string) error {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil || account == nil {
		return usecase.ErrAccountNotFound
	}

	account.Deactivate()
	return u.repo.Update(ctx, account)
}
