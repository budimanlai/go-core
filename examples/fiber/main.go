package main

import (
	"fmt"
	"log"

	accountHandler "github.com/budimanlai/go-core/account/handler"
	accountPersistence "github.com/budimanlai/go-core/account/platform/persistence"
	accountUsecase "github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-core/config"
	"github.com/budimanlai/go-core/middleware/auth"
	"github.com/budimanlai/go-core/middleware/cors"
	"github.com/budimanlai/go-core/middleware/logging"
	"github.com/budimanlai/go-core/middleware/recovery"
	"github.com/budimanlai/go-core/middleware/ratelimit"
	"github.com/budimanlai/go-core/pkg/crypto"
	"github.com/budimanlai/go-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	appLogger := logger.NewSimpleLogger()
	appLogger.Info("Starting application...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	appLogger.Info("Database connected successfully")

	passwordHasher := crypto.NewBcryptHasher(10)
	jwtService := auth.NewJWTService(auth.JWTConfig{
		SecretKey:       cfg.JWTSecret,
		Issuer:          cfg.JWTIssuer,
		ExpirationHours: cfg.JWTExpirationHours,
	})

	accountRepo := accountPersistence.NewAccountRepository(db)
	accountUC := accountUsecase.NewAccountUsecase(accountRepo, passwordHasher)
	accountHTTPHandler := accountHandler.NewAccountHandler(accountUC)

	app := fiber.New(fiber.Config{
		AppName:      "Go Core Example",
		ErrorHandler: customErrorHandler,
	})

	app.Use(recovery.FiberRecoveryMiddleware(recovery.DefaultConfig()))
	app.Use(cors.FiberCORSMiddleware(cors.DefaultConfig()))
	app.Use(logging.FiberLoggerMiddleware(logging.LoggerConfig{
		SkipPaths: []string{"/health"},
		LogFunc: func(entry logging.LogEntry) {
			appLogger.Info(
				"%s %s - Status: %d - Latency: %s - IP: %s",
				entry.Method,
				entry.Path,
				entry.StatusCode,
				entry.Latency,
				entry.IP,
			)
		},
	}))
	app.Use(ratelimit.FiberRateLimitMiddleware(ratelimit.DefaultConfig()))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "go-core",
		})
	})

	api := app.Group("/api/v1")

	public := api.Group("/public")
	public.Post("/register", accountHTTPHandler.Register)
	public.Post("/login", accountHTTPHandler.Login)

	protected := api.Group("/accounts")
	protected.Use(auth.FiberJWTMiddleware(jwtService))
	protected.Get("/", accountHTTPHandler.List)
	protected.Get("/:id", accountHTTPHandler.GetByID)
	protected.Delete("/:id", accountHTTPHandler.Delete)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	appLogger.Info("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    fmt.Sprintf("ERR_%d", code),
			"message": message,
		},
	})
}
