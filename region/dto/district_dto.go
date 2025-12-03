package dto

type DistrictResponse struct {
	DisID   uint   `json:"dis_id"`
	DisName string `json:"dis_name"`
}

type DistrictRequest struct {
	CityID uint `json:"city_id" validate:"required"`
}
