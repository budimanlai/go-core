package entity

import ()

type Subdistrict struct {
	SubdisId   int     `json:"subdis_id"`
	SubdisName *string `json:"subdis_name"`
	DisId      int     `json:"dis_id"`
}
