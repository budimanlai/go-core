package dto

import ()

type CreateSubdistrictReq struct {
	SubdisId   int     `json:"subdis_id" validate:"omitempty"`
	SubdisName *string `json:"subdis_name" validate:"omitempty"`
	DisId      int     `json:"dis_id" validate:"required"`
}

type UpdateSubdistrictReq struct {
	SubdisId   int     `json:"subdis_id"`
	SubdisName *string `json:"subdis_name"`
	DisId      int     `json:"dis_id"`
}
