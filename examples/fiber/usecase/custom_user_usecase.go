package usecase

import (
	"github.com/budimanlai/go-core/account/domain/entity"
	"github.com/budimanlai/go-core/account/domain/repository"
	"github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-core/account/dto"
	"github.com/budimanlai/go-core/account/platform/security"
	"github.com/budimanlai/go-pkg/logger"
	"github.com/jinzhu/copier"
)

type CustomUserUsecase struct {
	usecase.UserUsecase
	repo   repository.UserRepository
	hasher security.PasswordHasher
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type CustomUserResponse struct {
	Profile dto.UserResponse `json:"profile"`
	Token   tokenResponse    `json:"token"`
}

func NewCustomUserUsecase(userUsecase usecase.UserUsecase, repo repository.UserRepository, hasher security.PasswordHasher) usecase.UserUsecase {
	obj := &CustomUserUsecase{
		UserUsecase: userUsecase,
		repo:        repo,
		hasher:      hasher,
	}
	obj.SetCustomResponse(obj.toResponse)

	return obj
}

func (u *CustomUserUsecase) toResponse(user *entity.User) interface{} {
	logger.Printf("Custom toUserResponse")
	var response CustomUserResponse = CustomUserResponse{}
	copier.Copy(&response.Profile, user)
	response.Token.AccessToken = "access-token"
	response.Token.RefreshToken = "refresh-token"
	return &response
}
