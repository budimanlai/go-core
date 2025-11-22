package main

import (
	"fmt"
	"log"

	userHTTP "github.com/budimanlai/go-core/account/platform/http"
	userRepository "github.com/budimanlai/go-core/account/platform/repository"
	userSecurity "github.com/budimanlai/go-core/account/platform/security"
	userUsecase "github.com/budimanlai/go-core/account/platform/usecase"
	"github.com/budimanlai/go-core/config"
	"github.com/budimanlai/go-core/middleware/auth"
	"github.com/budimanlai/go-core/middleware/cors"
	"github.com/budimanlai/go-core/middleware/logging"
	"github.com/budimanlai/go-core/middleware/ratelimit"
	"github.com/budimanlai/go-core/middleware/recovery"
	"github.com/budimanlai/go-pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	logger.Printf("Starting application...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Printf("Database connected successfully")

	passwordHasher := userSecurity.NewBcryptHasher()
	jwtService := auth.NewJWTService(auth.JWTConfig{
		SecretKey:       cfg.JWTSecret,
		Issuer:          cfg.JWTIssuer,
		ExpirationHours: cfg.JWTExpirationHours,
	})

	userRepo := userRepository.NewUserRepository(db)
	userUC := userUsecase.NewUserUsecase(userRepo, passwordHasher)
	userHTTPHandler := userHTTP.NewUserHandler(userUC)

	app := fiber.New(fiber.Config{
		AppName:      "Go Core Example",
		ErrorHandler: customErrorHandler,
	})

	app.Use(recovery.FiberRecoveryMiddleware(recovery.DefaultConfig()))
	app.Use(cors.FiberCORSMiddleware(cors.DefaultConfig()))
	app.Use(logging.FiberLoggerMiddleware(logging.LoggerConfig{
		SkipPaths: []string{"/health"},
		LogFunc: func(entry logging.LogEntry) {
			logger.Infof(
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
	public.Post("/register", userHTTPHandler.Register)
	public.Post("/login", userHTTPHandler.Login)
	public.Get("/verify", userHTTPHandler.VerifyEmail)

	protected := api.Group("/users")
	protected.Use(auth.FiberJWTMiddleware(jwtService))
	protected.Get("/", userHTTPHandler.List)
	protected.Get("/:id", userHTTPHandler.GetByID)
	protected.Put("/:id", userHTTPHandler.Update)
	protected.Delete("/:id", userHTTPHandler.Delete)
	protected.Post("/:id/activate", userHTTPHandler.Activate)
	protected.Post("/:id/deactivate", userHTTPHandler.Deactivate)
	protected.Post("/:id/suspend", userHTTPHandler.Suspend)
	protected.Post("/:id/dashboard/enable", userHTTPHandler.EnableDashboard)
	protected.Post("/:id/dashboard/disable", userHTTPHandler.DisableDashboard)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	logger.Infof("Server starting on %s", addr)
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
