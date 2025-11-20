package logging

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoggerConfig struct {
	LogFunc   func(LogEntry)
	SkipPaths []string
}

type LogEntry struct {
	Timestamp    time.Time
	Method       string
	Path         string
	StatusCode   int
	Latency      time.Duration
	IP           string
	UserAgent    string
	ErrorMessage string
}

func FiberLoggerMiddleware(config LoggerConfig) fiber.Handler {
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *fiber.Ctx) error {
		if skipPaths[c.Path()] {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		var errorMessage string
		if err != nil {
			errorMessage = err.Error()
		}

		entry := LogEntry{
			Timestamp:    start,
			Method:       c.Method(),
			Path:         c.Path(),
			StatusCode:   c.Response().StatusCode(),
			Latency:      latency,
			IP:           c.IP(),
			UserAgent:    c.Get("User-Agent"),
			ErrorMessage: errorMessage,
		}

		if config.LogFunc != nil {
			config.LogFunc(entry)
		}

		return err
	}
}
