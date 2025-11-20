package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

type RecoveryConfig struct {
	PanicHandler     func(c *fiber.Ctx, err interface{})
	EnableStackTrace bool
}

func DefaultConfig() RecoveryConfig {
	return RecoveryConfig{
		PanicHandler:     nil,
		EnableStackTrace: false,
	}
}

func FiberRecoveryMiddleware(config RecoveryConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()

				if config.PanicHandler != nil {
					config.PanicHandler(c, r)
				}

				fmt.Printf("PANIC: %v\n%s\n", r, stack)

				response := fiber.Map{
					"error": "internal server error",
				}

				if config.EnableStackTrace {
					response["panic"] = fmt.Sprintf("%v", r)
					response["stack"] = string(stack)
				}

				c.Status(fiber.StatusInternalServerError).JSON(response)
			}
		}()

		return c.Next()
	}
}
