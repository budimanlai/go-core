package http

import (
	"strconv"

	"github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-core/account/dto"

	"github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequestI18n(c, "app.invalid_request_body", nil)
	}

	if err := validator.ValidateStruct(req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	user, err := h.usecase.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequestI18n(c, "app.invalid_request_body", nil)
	}

	loginResp, err := h.usecase.Login(&req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.SuccessI18n(c, "app.success", loginResp)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequestI18n(c, "app.invalid_request_body", nil)
	}

	user, err := h.usecase.GetByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.SuccessI18n(c, "app.success", user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequestI18n(c, "app.invalid_request_body", nil)
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequestI18n(c, "app.invalid_request_body", nil)
	}

	user, err := h.usecase.Update(uint(id), &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequestI18n(c, "app.invalid_request_body", nil)
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	listResp, err := h.usecase.List(page, pageSize)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessI18n(c, "app.success", listResp)
}

func (h *UserHandler) Activate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.usecase.Activate(uint(id)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}

func (h *UserHandler) Deactivate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.usecase.Deactivate(uint(id)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}

func (h *UserHandler) Suspend(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.usecase.Suspend(uint(id)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return response.BadRequest(c, "Verification token is required")
	}

	if err := h.usecase.VerifyEmail(token); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}

func (h *UserHandler) EnableDashboard(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.usecase.EnableDashboard(uint(id)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}

func (h *UserHandler) DisableDashboard(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.usecase.DisableDashboard(uint(id)); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessI18n(c, "app.success", nil)
}
