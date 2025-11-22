package usecase

import (
	"context"

	"github.com/budimanlai/go-core/account/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) (*dto.ListUserResponse, error)
	Activate(ctx context.Context, id uint) error
	Deactivate(ctx context.Context, id uint) error
	Suspend(ctx context.Context, id uint) error
	VerifyEmail(ctx context.Context, token string) error
	EnableDashboard(ctx context.Context, id uint) error
	DisableDashboard(ctx context.Context, id uint) error
}
