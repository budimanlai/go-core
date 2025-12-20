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

type SubdistrictHandler struct {
	service service.SubdistrictService
}

func NewSubdistrictHandler(service service.SubdistrictService) *SubdistrictHandler {
	return &SubdistrictHandler{
		service: service,
	}
}

// Index godoc
// @Summary      List Subdistrict
// @Description  Get paginated list of Subdistrict
// @Tags         Subdistrict
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Success      200  {object}  response.PaginationResult{data=[]entity.Subdistrict}
// @Router       /api/v1/region/subdistricts [get]
func (h *SubdistrictHandler) Index(c *fiber.Ctx) error {
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
// @Summary      Get Subdistrict Detail
// @Description  Get single Subdistrict by ID
// @Tags         Subdistrict
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subdistrict ID"
// @Success      200  {object}  response.APIResponse{data=entity.Subdistrict}
// @Failure      404  {object}  response.APIResponse{data=nil}
// @Router       /api/v1/region/subdistricts/{id} [get]
func (h *SubdistrictHandler) View(c *fiber.Ctx) error {
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
// @Summary      Create new Subdistrict
// @Description  Create a new Subdistrict record
// @Tags         Subdistrict
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateSubdistrictReq true "Create Request"
// @Success      201  {object}  response.APIResponse{data=entity.Subdistrict}
// @Failure      400  {object}  response.APIResponse{data=nil}
// @Router       /api/v1/region/subdistricts [post]
func (h *SubdistrictHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSubdistrictReq
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorI18n(c, fiber.StatusBadRequest, "app.error.invalid_request_body", nil)
	}

	// Gunakan Validator dari go-pkg
	if err := validator.ValidateStructWithContext(c, &req); err != nil {
		return response.ValidationErrorI18n(c, err)
	}

	// Mapping menggunakan Copier
	var item entity.Subdistrict
	if err := copier.Copy(&item, &req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Mapping failed")
	}

	if err := h.service.Create(c.Context(), &item); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", item)
}

// Update godoc
// @Summary      Update Subdistrict
// @Description  Update existing Subdistrict
// @Tags         Subdistrict
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subdistrict ID"
// @Param        body body dto.UpdateSubdistrictReq true "Update Request"
// @Success      200  {object}  response.APIResponse{data=entity.Subdistrict}
// @Router       /api/v1/region/subdistricts/{id} [put]
func (h *SubdistrictHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	// 1. Cek Eksistensi Data (Sesuai Pattern BaseHandler)
	existing, err := h.service.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return response.ErrorI18n(c, fiber.StatusNotFound, "app.error.not_found", nil)
	}

	// 2. Parse Body
	var req dto.UpdateSubdistrictReq
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
// @Summary      Delete Subdistrict
// @Description  Delete Subdistrict by ID
// @Tags         Subdistrict
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subdistrict ID"
// @Success      200  {object}  response.APIResponse{data=entity.Subdistrict}
// @Router       /api/v1/region/subdistricts/{id} [delete]
func (h *SubdistrictHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.service.Delete(c.Context(), id); err != nil {
		return response.ErrorI18n(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessI18n(c, "app.success", fiber.Map{"deleted": true})
}
