package ratelimit

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type RateLimitConfig struct {
	Max          int
	Expiration   time.Duration
	KeyGenerator func(*fiber.Ctx) string
	LimitReached fiber.Handler
	Skip         func(*fiber.Ctx) bool
}

func DefaultConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		},
	}
}

func FiberRateLimitMiddleware(config RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:          config.Max,
		Expiration:   config.Expiration,
		KeyGenerator: config.KeyGenerator,
		LimitReached: config.LimitReached,
	})
}
