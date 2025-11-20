package handler

import (
	"math"
	"strconv"

	"github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-core/account/dto"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	usecase usecase.AccountUsecase
}

func NewAccountHandler(usecase usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{
		usecase: usecase,
	}
}

func (h *AccountHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	account, err := h.usecase.Register(c.Context(), req.Email, req.Username, req.Password, req.FullName)
	if err != nil {
		if err == usecase.ErrAccountAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to register account",
		})
	}

	response := dto.AccountResponse{
		ID:        account.ID,
		Email:     account.Email,
		Username:  account.Username,
		FullName:  account.FullName,
		Role:      account.Role,
		IsActive:  account.IsActive,
		CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": response,
	})
}

func (h *AccountHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	account, err := h.usecase.Login(c.Context(), req.Identifier, req.Password)
	if err != nil {
		if err == usecase.ErrInvalidCredentials || err == usecase.ErrAccountInactive {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to login",
		})
	}

	response := dto.LoginResponse{
		Account: dto.AccountResponse{
			ID:        account.ID,
			Email:     account.Email,
			Username:  account.Username,
			FullName:  account.FullName,
			Role:      account.Role,
			IsActive:  account.IsActive,
			CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		AccessToken: "TODO: generate token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (h *AccountHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	account, err := h.usecase.GetByID(c.Context(), id)
	if err != nil {
		if err == usecase.ErrAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get account",
		})
	}

	response := dto.AccountResponse{
		ID:        account.ID,
		Email:     account.Email,
		Username:  account.Username,
		FullName:  account.FullName,
		Role:      account.Role,
		IsActive:  account.IsActive,
		CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (h *AccountHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	accounts, total, err := h.usecase.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to list accounts",
		})
	}

	accountResponses := make([]dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		accountResponses[i] = dto.AccountResponse{
			ID:        account.ID,
			Email:     account.Email,
			Username:  account.Username,
			FullName:  account.FullName,
			Role:      account.Role,
			IsActive:  account.IsActive,
			CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := dto.ListAccountResponse{
		Data:       accountResponses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (h *AccountHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	if err := h.usecase.Delete(c.Context(), id); err != nil {
		if err == usecase.ErrAccountNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account deleted successfully",
	})
}
