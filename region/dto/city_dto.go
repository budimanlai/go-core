package dto

type CityResponse struct {
	CityID   uint   `json:"city_id"`
	CityName string `json:"city_name"`
}

type CityRequest struct {
	ProvinceID uint `json:"prov_id" validate:"required"`
}
