package http

import (
	"strconv"

	"github.com/budimanlai/go-core/region/domain/usecase"
	"github.com/budimanlai/go-core/region/dto"
	"github.com/budimanlai/go-pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type RegionHandler struct {
	countryUC     usecase.CountryUsecase
	provinceUC    usecase.ProvinceUsecase
	cityUC        usecase.CityUsecase
	districtUC    usecase.DistrictUsecase
	subdistrictUC usecase.SubDistrictUsecase
}

func NewRegionHandler(
	countryUC usecase.CountryUsecase,
	provinceUC usecase.ProvinceUsecase,
	cityUC usecase.CityUsecase,
	districtUC usecase.DistrictUsecase,
	subdistrictUC usecase.SubDistrictUsecase,
) *RegionHandler {
	return &RegionHandler{
		countryUC:     countryUC,
		provinceUC:    provinceUC,
		cityUC:        cityUC,
		districtUC:    districtUC,
		subdistrictUC: subdistrictUC,
	}
}

// Country Handlers
func (h *RegionHandler) GetCountries(c *fiber.Ctx) error {
	countries, err := h.countryUC.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get countries")
	}

	var countryResponses []dto.CountryResponse
	if err := copier.Copy(&countryResponses, &countries); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map countries")
	}

	return response.SuccessI18n(c, "app.success", countryResponses)
}

func (h *RegionHandler) GetCountryByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return response.BadRequest(c, "Country code is required")
	}

	country, err := h.countryUC.GetByCode(code)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Country not found")
	}

	var countryResponse dto.CountryResponse
	if err := copier.Copy(&countryResponse, &country); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map country")
	}
	countryResponse.CountryName = country.Name
	countryResponse.CountryCode = country.IsoAlpha2

	return response.SuccessI18n(c, "app.success", countryResponse)
}

// Province Handlers
func (h *RegionHandler) GetProvinces(c *fiber.Ctx) error {
	provinces, err := h.provinceUC.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get provinces")
	}

	var out []dto.ProvinceResponse
	if err := copier.Copy(&out, &provinces); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map provinces")
	}

	return response.SuccessI18n(c, "app.success", out)
}

// City Handlers
func (h *RegionHandler) GetCityByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid city ID")
	}

	city, err := h.cityUC.GetByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "City not found")
	}

	var cityResponse dto.CityResponse
	if err := copier.Copy(&cityResponse, &city); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map city")
	}

	return response.SuccessI18n(c, "app.success", cityResponse)
}

func (h *RegionHandler) GetCitiesByProvince(c *fiber.Ctx) error {
	provIDStr := c.Query("prov_id")
	if provIDStr == "" {
		return response.BadRequest(c, "Province ID is required")
	}

	provID, err := strconv.ParseUint(provIDStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid province ID")
	}

	cities, err := h.cityUC.GetAllByProvince(uint(provID))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get cities")
	}

	var cityResponses []dto.CityResponse
	if err := copier.Copy(&cityResponses, &cities); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map cities")
	}

	return response.SuccessI18n(c, "app.success", cityResponses)
}

// District Handlers
func (h *RegionHandler) GetDistrictByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid district ID")
	}

	district, err := h.districtUC.GetByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "District not found")
	}

	var districtResponse dto.DistrictResponse
	if err := copier.Copy(&districtResponse, &district); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map district")
	}

	return response.SuccessI18n(c, "app.success", districtResponse)
}

func (h *RegionHandler) GetDistrictsByCity(c *fiber.Ctx) error {
	cityIDStr := c.Query("city_id")
	if cityIDStr == "" {
		return response.BadRequest(c, "City ID is required")
	}

	cityID, err := strconv.ParseUint(cityIDStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid city ID")
	}

	districts, err := h.districtUC.GetAllByCity(uint(cityID))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get districts")
	}

	var districtResponses []dto.DistrictResponse
	if err := copier.Copy(&districtResponses, &districts); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map districts")
	}

	return response.SuccessI18n(c, "app.success", districtResponses)
}

// SubDistrict Handlers
func (h *RegionHandler) GetSubDistrictsByDistrict(c *fiber.Ctx) error {
	disIDStr := c.Query("dis_id")
	if disIDStr == "" {
		return response.BadRequest(c, "District ID is required")
	}

	disID, err := strconv.ParseUint(disIDStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid district ID")
	}

	subdistricts, err := h.subdistrictUC.GetAllByDistrict(uint(disID))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get subdistricts")
	}

	var subdistrictResponses []dto.SubDistrictResponse
	if err := copier.Copy(&subdistrictResponses, &subdistricts); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to map subdistricts")
	}

	return response.SuccessI18n(c, "app.success", subdistrictResponses)
}
