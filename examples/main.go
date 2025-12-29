package main

import (
	"fmt"
	"time"

	pkg_databases "github.com/budimanlai/go-pkg/databases"
	"github.com/gofiber/fiber/v2"

	handler "github.com/budimanlai/go-core/auth/handler/http"
	auth_service "github.com/budimanlai/go-core/auth/service"
	"github.com/budimanlai/go-core/service"

	common_repository "github.com/budimanlai/go-core/common/repository"
	common_usecase "github.com/budimanlai/go-core/common/usecase"

	"github.com/budimanlai/go-core/auth/repository"
	"github.com/budimanlai/go-core/auth/usecase"
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/middleware/auth"
)

func main() {
	fmt.Println("Example login usecase implementation")

	dbConfig := pkg_databases.DbConfig{
		Driver:   pkg_databases.Postgres,
		Host:     "127.0.0.1",
		Port:     "5432",
		Name:     "evoucher",
		Username: "dev",
		Password: "12345678",
	}

	dbManager := pkg_databases.NewDbManager(dbConfig)
	err := dbManager.Open()
	if err != nil {
		panic(err)
	}
	defer dbManager.Close()

	db := dbManager.GetDb()

	jwtConfig := auth.JWTConfig{
		SecretKey:       "your-secret-key",
		Issuer:          "your-app-name",
		ExpirationHours: 72,
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// create repo factory
	repoConfig := base.RepoConfig{
		EnableCache:      false,
		EnablePrometheus: false,
		RedisClient:      nil,
	}
	repoFactory := base.NewFactory(db, repoConfig)

	// create repository
	userRepo := repository.NewUserRepository(repoFactory)
	userSessionRepo := repository.NewUserSessionRepository(repoFactory)
	otpRepo := repository.NewOtpRepository(repoFactory)
	messagingTemplateRepo := common_repository.NewMessagingTemplateRepositoryImpl(repoFactory)

	// create usecase
	userSessionUsecase := usecase.NewUserSessionUsecaseImpl(db, userSessionRepo, userRepo, jwtService)
	userSessionUsecase.SetMultipleLoginAllowed(false) // allow multiple login
	messagingTemplateUsecase := common_usecase.NewMessagingTemplateUsecaseImpl(db, messagingTemplateRepo)

	// mail service
	mailConfig := service.SMTPMailServiceConfig{
		Host:     "127.0.0.1",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "no-reply@example.com",
	}
	emailService := service.NewSMTPMailServiceImpl(mailConfig, messagingTemplateUsecase)
	waviroService := service.NewWaviroServiceImpl()

	// create OTP sender service
	otpSenderService := auth_service.NewOtpSenderServiceImpl(
		emailService,
		waviroService,
	)

	// create OTP usecase
	otpConfig := usecase.OtpConfig{
		UserInitiated:      false,
		BotPhoneNumber:     "1234567890",
		CommandPrefix:      "OTP#",
		MaxPendingRequests: 3,
		ExpiredDuration:    60 * time.Minute, // OTP expires in 60 minutes
	}
	otpUsecase := usecase.NewOtpUsecaseImpl(db, otpRepo, otpConfig)
	otpUsecase.SetSender(otpSenderService)

	// create handler
	authHandler := handler.NewAuthHandler(userSessionUsecase, otpUsecase)

	app := fiber.New()
	api := app.Group("/api/v1")

	basicApiAuthService := auth.NewBasicAuthService(auth.BasicAuthConfig{
		Users: map[string]string{
			"admin": "admin123",
			"user":  "user123",
		},
	})

	// Basic Auth Middleware
	authEndpoint := api.Group("/auth")
	authEndpoint.Post("/login", auth.FiberBasicAuthMiddleware(basicApiAuthService), authHandler.Login)
	authEndpoint.Post("/otp/request", auth.FiberBasicAuthMiddleware(basicApiAuthService), authHandler.RequestOtp)
	authEndpoint.Post("/otp/status", auth.FiberBasicAuthMiddleware(basicApiAuthService), authHandler.StatusOTP)
	authEndpoint.Post("/otp/verify", auth.FiberBasicAuthMiddleware(basicApiAuthService), authHandler.VerifyOTP)

	// JWT Auth Middleware
	jwtRestAPI := api.Group("/auth", auth.FiberJWTMiddleware(jwtService))
	jwtRestAPI.Post("/logout", authHandler.Logout)
	jwtRestAPI.Post("/token/verify", authHandler.VerifyToken)
	jwtRestAPI.Post("/token/refresh", authHandler.RefreshToken)

	// setup routes
	if err := app.Listen(":8084"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
