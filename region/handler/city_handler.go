package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"

	"github.com/budimanlai/go-core/region/domain/entity"
	"github.com/budimanlai/go-core/region/dto"
	"github.com/budimanlai/go-core/region/service"

	"github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/validator"
)

type CityHandler struct {
	service service.CityService
}

func NewCityHandler(service service.CityService) *CityHandler {
	return &CityHandler{
		service: service,
	}
}

// Index godoc
// @Summary      List City
// @Description  Get paginated list of City
// @Tags         City
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Success      200  {object}  response.PaginationResult{data=[]entity.City}
// @Router       /api/v1/region/cities [get]
func (h *CityHandler) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.service.FindAll(c.Context(), page, limit)
	if err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessWithPagination(c, "app.success", response.PaginationResult{
		Data:      result.Data,
		Total:     result.Total,
		Page:      result.Page,
		Limit:     result.Limit,
		TotalPage: result.TotalPage,
	})
}

// View godoc
// @Summary      Get City Detail
// @Description  Get single City by ID
// @Tags         City
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "City ID"
// @Success      200  {object}  response.APIResponse{data=entity.City}
// @Failure      404  {object}  response.APIResponse{data=nil}
// @Router       /api/v1/region/cities/{id} [get]
func (h *CityHandler) View(c *fiber.Ctx) error {
	// Asumsi ID Integer (standard SQL), ubah jika UUID
	id, _ := strconv.Atoi(c.Params("id"))

	item, err := h.service.FindByID(c.Context(), id)
	if err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Safety check jika service return nil error tapi data kosong
	if item == nil {
		return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
	}

	return response.SuccessI18n(c, "app.success", item)
}

// Create godoc
// @Summary      Create new City
// @Description  Create a new City record
// @Tags         City
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateCityReq true "Create Request"
// @Success      201  {object}  response.APIResponse{data=entity.City}
// @Failure      400  {object}  response.APIResponse{data=nil}
// @Router       /api/v1/region/cities [post]
func (h *CityHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCityReq
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorI18n(c, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// Gunakan Validator dari go-pkg
	if err := validator.ValidateStructWithContext(c, &req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	// Mapping menggunakan Copier
	var item entity.City
	if err := copier.Copy(&item, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Mapping failed")
	}

	if err := h.service.Create(c.Context(), &item); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", item)
}

// Update godoc
// @Summary      Update City
// @Description  Update existing City
// @Tags         City
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "City ID"
// @Param        body body dto.UpdateCityReq true "Update Request"
// @Success      200  {object}  response.APIResponse{data=entity.City}
// @Router       /api/v1/region/cities/{id} [put]
func (h *CityHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	// 1. Cek Eksistensi Data (Sesuai Pattern BaseHandler)
	existing, err := h.service.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
	}

	// 2. Parse Body
	var req dto.UpdateCityReq
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorI18n(c, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// 3. Validate
	if err := validator.ValidateStructWithContext(c, &req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	// 4. Merge Data (Req -> Existing)
	if err := copier.Copy(existing, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Mapping failed")
	}

	// 5. Save Update
	// Kita kirim object 'existing' yang sudah terupdate field-nya
	if err := h.service.Update(c.Context(), existing); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", existing)
}

// Delete godoc
// @Summary      Delete City
// @Description  Delete City by ID
// @Tags         City
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "City ID"
// @Success      200  {object}  response.APIResponse{data=entity.City}
// @Router       /api/v1/region/cities/{id} [delete]
func (h *CityHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.service.Delete(c.Context(), id); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", fiber.Map{"deleted": true})
}
