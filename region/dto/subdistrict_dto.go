package dto

type SubDistrictResponse struct {
	SubdisID   uint   `json:"subdis_id"`
	SubdisName string `json:"subdis_name"`
}

type SubDistrictRequest struct {
	DisID uint `json:"dis_id" validate:"required"`
}
