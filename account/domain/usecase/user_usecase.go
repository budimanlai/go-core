package usecase

import (
	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/dto"
)

type UserUsecase interface {
	Register(req *dto.RegisterRequest) (interface{}, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetByID(id uint) (*dto.UserResponse, error)
	Update(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id uint) error
	List(page, pageSize int) (*dto.ListUserResponse, error)
	Activate(id uint) error
	Deactivate(id uint) error
	Suspend(id uint) error
	VerifyEmail(token string) error
	EnableDashboard(id uint) error
	DisableDashboard(id uint) error
	SetCustomResponse(customToResponse func(*entity.User) interface{})
}
