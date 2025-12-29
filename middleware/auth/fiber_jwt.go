package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type JWTMiddleware struct {
	JwtService JWTService
}

func FiberJWTMiddleware(jwtService JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		tokenString := authHeader[len(bearerPrefix):]

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("user_token", claims.UserToken)
		c.Locals("claims", claims)

		return c.Next()
	}
}

func FiberBasicAuthMiddleware(basicAuthService BasicAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		username, err := basicAuthService.Validate(authHeader)
		if err != nil {
			c.Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}

		c.Locals("username", username)

		return c.Next()
	}
}
