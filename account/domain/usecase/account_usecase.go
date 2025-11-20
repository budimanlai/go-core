package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAccountInactive      = errors.New("account is inactive")
)

type AccountUsecase interface {
	Register(ctx context.Context, email, username, password, fullName string) (*entity.Account, error)
	Login(ctx context.Context, identifier, password string) (*entity.Account, error)
	GetByID(ctx context.Context, id string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Account, int64, error)
	Activate(ctx context.Context, id string) error
	Deactivate(ctx context.Context, id string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

type accountUsecase struct {
	repo           repository.AccountRepository
	passwordHasher PasswordHasher
}

func NewAccountUsecase(repo repository.AccountRepository, hasher PasswordHasher) AccountUsecase {
	return &accountUsecase{
		repo:           repo,
		passwordHasher: hasher,
	}
}

func (u *accountUsecase) Register(ctx context.Context, email, username, password, fullName string) (*entity.Account, error) {
	existing, _ := u.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrAccountAlreadyExists
	}

	existing, _ = u.repo.FindByUsername(ctx, username)
	if existing != nil {
		return nil, ErrAccountAlreadyExists
	}

	hashedPassword, err := u.passwordHasher.Hash(password)
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

func (u *accountUsecase) Login(ctx context.Context, identifier, password string) (*entity.Account, error) {
	account, err := u.repo.FindByEmail(ctx, identifier)
	if err != nil || account == nil {
		account, err = u.repo.FindByUsername(ctx, identifier)
		if err != nil || account == nil {
			return nil, ErrInvalidCredentials
		}
	}

	if !account.IsActive {
		return nil, ErrAccountInactive
	}

	if !u.passwordHasher.Verify(account.Password, password) {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}

func (u *accountUsecase) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (u *accountUsecase) Update(ctx context.Context, account *entity.Account) error {
	existing, err := u.repo.FindByID(ctx, account.ID)
	if err != nil || existing == nil {
		return ErrAccountNotFound
	}

	account.UpdatedAt = time.Now()
	return u.repo.Update(ctx, account)
}

func (u *accountUsecase) Delete(ctx context.Context, id string) error {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil || account == nil {
		return ErrAccountNotFound
	}

	return u.repo.SoftDelete(ctx, id)
}

func (u *accountUsecase) List(ctx context.Context, limit, offset int) ([]*entity.Account, int64, error) {
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

func (u *accountUsecase) Activate(ctx context.Context, id string) error {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil || account == nil {
		return ErrAccountNotFound
	}

	account.Activate()
	return u.repo.Update(ctx, account)
}

func (u *accountUsecase) Deactivate(ctx context.Context, id string) error {
	account, err := u.repo.FindByID(ctx, id)
	if err != nil || account == nil {
		return ErrAccountNotFound
	}

	account.Deactivate()
	return u.repo.Update(ctx, account)
}
