package handler

import (
	"time"

	"github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/types"
	"github.com/gofiber/fiber/v2"
)

type CommonHandler struct {
}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) Ping(c *fiber.Ctx) error {
	return response.SuccessI18n(c, "app.success", map[string]interface{}{
		"message": "pong",
		"time":    types.UTCTime(time.Now()).String(),
	})
}
