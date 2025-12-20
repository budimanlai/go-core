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

type ProvinceHandler struct {
	service service.ProvinceService
}

func NewProvinceHandler(service service.ProvinceService) *ProvinceHandler {
	return &ProvinceHandler{
		service: service,
	}
}

// Index godoc
// @Summary      List Province
// @Description  Get paginated list of Province
// @Tags         Province
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Success      200  {object}  response.PaginationResult{data=[]entity.Province}
// @Router       /api/v1/region/provinces [get]
func (h *ProvinceHandler) Index(c *fiber.Ctx) error {
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
// @Summary      Get Province Detail
// @Description  Get single Province by ID
// @Tags         Province
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Province ID"
// @Success      200  {object}  response.APIResponse{data=entity.Province}
// @Failure      404  {object}  response.APIResponse{data=nil}
// @Router       /api/v1/region/provinces/{id} [get]
func (h *ProvinceHandler) View(c *fiber.Ctx) error {
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
// @Summary      Create new Province
// @Description  Create a new Province record
// @Tags         Province
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateProvinceReq true "Create Request"
// @Success      201  {object}  response.APIResponse{data=entity.Province}
// @Failure      400  {object}  response.APIResponse{data=nil}
// @Router       /api/v1/region/provinces [post]
func (h *ProvinceHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProvinceReq
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorI18n(c, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// Gunakan Validator dari go-pkg
	if err := validator.ValidateStructWithContext(c, &req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	// Mapping menggunakan Copier
	var item entity.Province
	if err := copier.Copy(&item, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Mapping failed")
	}

	if err := h.service.Create(c.Context(), &item); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", item)
}

// Update godoc
// @Summary      Update Province
// @Description  Update existing Province
// @Tags         Province
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Province ID"
// @Param        body body dto.UpdateProvinceReq true "Update Request"
// @Success      200  {object}  response.APIResponse{data=entity.Province}
// @Router       /api/v1/region/provinces/{id} [put]
func (h *ProvinceHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	// 1. Cek Eksistensi Data (Sesuai Pattern BaseHandler)
	existing, err := h.service.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
	}

	// 2. Parse Body
	var req dto.UpdateProvinceReq
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
// @Summary      Delete Province
// @Description  Delete Province by ID
// @Tags         Province
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Province ID"
// @Success      200  {object}  response.APIResponse{data=entity.Province}
// @Router       /api/v1/region/provinces/{id} [delete]
func (h *ProvinceHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.service.Delete(c.Context(), id); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", fiber.Map{"deleted": true})
}
