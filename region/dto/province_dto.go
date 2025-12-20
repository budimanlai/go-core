package dto

import ()

type CreateProvinceReq struct {
	ProvId     int     `json:"prov_id" validate:"omitempty"`
	ProvName   *string `json:"prov_name" validate:"omitempty"`
	Locationid *int    `json:"locationid" validate:"omitempty"`
	Status     *int    `json:"status" validate:"omitempty"`
}

type UpdateProvinceReq struct {
	ProvId     int     `json:"prov_id"`
	ProvName   *string `json:"prov_name"`
	Locationid *int    `json:"locationid"`
	Status     *int    `json:"status"`
}
