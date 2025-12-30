package main

import (
	"fmt"
	"time"

	pkg_databases "github.com/budimanlai/go-pkg/databases"
	pkg_middleware "github.com/budimanlai/go-pkg/middleware/auth"

	"github.com/gofiber/fiber/v2"

	auth_nmanager "github.com/budimanlai/go-core/auth"
	auth_service "github.com/budimanlai/go-core/auth/service"
	impl_common_repository "github.com/budimanlai/go-core/common/repository"
	impl_common_usecase "github.com/budimanlai/go-core/common/usecase"
	"github.com/budimanlai/go-core/service"

	"github.com/budimanlai/go-core/auth/usecase"
	"github.com/budimanlai/go-core/base"
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

	// create repo factory
	repoConfig := base.RepoConfig{
		EnableCache:      false,
		EnablePrometheus: false,
		RedisClient:      nil,
	}
	repoFactory := base.NewFactory(db, repoConfig)

	messagingTemplateRepo := impl_common_repository.NewMessagingTemplateRepositoryImpl(repoFactory)
	messagingTemplateUsecase := impl_common_usecase.NewMessagingTemplateUsecaseImpl(db, messagingTemplateRepo)

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
	otpConfig := usecase.OtpConfig{
		UserInitiated:      false,
		BotPhoneNumber:     "1234567890",
		CommandPrefix:      "OTP#",
		MaxPendingRequests: 3,
		ExpiredDuration:    60 * time.Minute, // OTP expires in 60 minutes
	}

	jwtConfig := pkg_middleware.JWTConfig{
		SecretKey:      "your-secret-key",
		Issuer:         "your-app-name",
		ExpirationTime: 72 * time.Hour,
		SigningMethod:  "HS256",
		TokenLookup:    "header:Authorization",
	}

	basicAuthConfig := pkg_middleware.BasicAuthConfig{
		KeyProvider:  pkg_middleware.NewDbKeyProvider(db),
		Unauthorized: func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusUnauthorized) },
	}
	basicAuthMiddleware := pkg_middleware.NewBasicAuth(basicAuthConfig)

	authManager := auth_nmanager.NewAuthManagerDefaultImpl(repoFactory)
	authManager.SetJwtConfig(jwtConfig)
	authManager.SetOtpSenderService(otpSenderService, otpConfig)
	authManager.SetPublicMiddleware(basicAuthMiddleware.Middleware())
	authManager.InitManager()

	app := fiber.New()
	api := app.Group("/api/v1")
	authManager.SetRoute(api)

	api.Get("ping", authManager.PublicMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// setup routes
	if err := app.Listen(":8084"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
