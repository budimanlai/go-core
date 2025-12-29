package auth

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/middleware/auth"
	"github.com/gofiber/fiber/v2"

	dom_auth_repository "github.com/budimanlai/go-core/auth/domain/repository"
	dom_auth_usecase "github.com/budimanlai/go-core/auth/domain/usecase"
	dom_auth_handler "github.com/budimanlai/go-core/auth/handler/http"
	auth_service "github.com/budimanlai/go-core/auth/service"
	"github.com/budimanlai/go-core/auth/usecase"

	impl_auth_repository "github.com/budimanlai/go-core/auth/repository"
	impl_auth_usecase "github.com/budimanlai/go-core/auth/usecase"
)

type AuthManagerDefaultImpl struct {
	factory *base.Factory

	// repository
	UserRepo        dom_auth_repository.UserRepository
	UserSessionRepo dom_auth_repository.UserSessionRepository
	OtpRepo         dom_auth_repository.OtpRepository

	// usecase
	UserUsecase        dom_auth_usecase.UserUsecase
	UserSessionUsecase dom_auth_usecase.UserSessionUsecase
	OtpUsecase         dom_auth_usecase.OtpUsecase

	// handler
	AuthHandler *dom_auth_handler.AuthHandler

	// service
	JwtService       auth.JWTService
	OtpSenderService auth_service.OtpSenderService
	OtpConfig        impl_auth_usecase.OtpConfig

	// middleware
	// PublicMiddleware is for routes that do not require user session
	PublicMiddleware fiber.Handler

	// PrivateMiddleware is for routes that require valid user session
	PrivateMiddleware fiber.Handler
}

func NewAuthManagerDefaultImpl(factory *base.Factory) *AuthManagerDefaultImpl {
	return &AuthManagerDefaultImpl{
		factory: factory,
	}
}

func (m *AuthManagerDefaultImpl) SetJwtService(jwtService auth.JWTService) {
	m.JwtService = jwtService
}

func (m *AuthManagerDefaultImpl) SetOtpSenderService(otpSenderService auth_service.OtpSenderService, config impl_auth_usecase.OtpConfig) {
	m.OtpSenderService = otpSenderService
	m.OtpConfig = config
}

func (m *AuthManagerDefaultImpl) SetPublicMiddleware(middleware fiber.Handler) {
	m.PublicMiddleware = middleware
}

func (m *AuthManagerDefaultImpl) SetPrivateMiddleware(middleware fiber.Handler) {
	m.PrivateMiddleware = middleware
}

func (m *AuthManagerDefaultImpl) InitManager() {
	if m.JwtService == nil {
		panic("JWT Service is not set in AuthManager")
	}

	m.initContainer()
	m.initUsecase()
}

func (m *AuthManagerDefaultImpl) initContainer() {
	m.UserRepo = impl_auth_repository.NewUserRepositoryImpl(m.factory)
	m.UserSessionRepo = impl_auth_repository.NewUserSessionRepositoryImpl(m.factory)
	m.OtpRepo = impl_auth_repository.NewOtpRepositoryImpl(m.factory)
}

func (m *AuthManagerDefaultImpl) initUsecase() {
	m.UserSessionUsecase = impl_auth_usecase.NewUserSessionUsecaseImpl(m.factory.DB, m.UserSessionRepo, m.UserRepo, m.JwtService)
	m.UserSessionUsecase.SetMultipleLoginAllowed(false) // allow multiple login

	m.OtpUsecase = usecase.NewOtpUsecaseImpl(m.factory.DB, m.OtpRepo, m.OtpConfig)
	m.OtpUsecase.SetSender(m.OtpSenderService)

	m.UserUsecase = impl_auth_usecase.NewUserUsecaseImpl(m.factory.DB, m.UserRepo, m.OtpUsecase)
}

func (m *AuthManagerDefaultImpl) SetRoute(app fiber.Router) {
	m.AuthHandler = dom_auth_handler.NewAuthHandler(m.UserUsecase, m.UserSessionUsecase, m.OtpUsecase)

	// Basic Auth Middleware
	authEndpoint := app.Group("/auth")
	authEndpoint.Post("/login", m.PublicMiddleware, m.AuthHandler.Login)
	authEndpoint.Post("/otp/request", m.PublicMiddleware, m.AuthHandler.RequestOtp)
	authEndpoint.Post("/otp/status", m.PublicMiddleware, m.AuthHandler.StatusOTP)
	authEndpoint.Post("/otp/verify", m.PublicMiddleware, m.AuthHandler.VerifyOTP)
	authEndpoint.Post("/password/reset", m.PublicMiddleware, m.AuthHandler.ResetPassword)

	// JWT Auth Middleware
	jwtRestAPI := app.Group("/auth", m.PrivateMiddleware)
	jwtRestAPI.Post("/logout", m.AuthHandler.Logout)
	jwtRestAPI.Post("/token/verify", m.AuthHandler.VerifyToken)
	jwtRestAPI.Post("/token/refresh", m.AuthHandler.RefreshToken)
}
