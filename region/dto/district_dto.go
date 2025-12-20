package dto

import ()

type CreateDistrictReq struct {
	DisId   int     `json:"dis_id" validate:"omitempty"`
	DisName *string `json:"dis_name" validate:"omitempty"`
	CityId  int     `json:"city_id" validate:"required"`
}

type UpdateDistrictReq struct {
	DisId   int     `json:"dis_id"`
	DisName *string `json:"dis_name"`
	CityId  int     `json:"city_id"`
}
