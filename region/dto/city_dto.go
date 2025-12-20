package dto

import ()

type CreateCityReq struct {
	CityId   int     `json:"city_id" validate:"omitempty"`
	CityName *string `json:"city_name" validate:"omitempty"`
	ProvId   int     `json:"prov_id" validate:"required"`
}

type UpdateCityReq struct {
	CityId   int     `json:"city_id"`
	CityName *string `json:"city_name"`
	ProvId   int     `json:"prov_id"`
}
