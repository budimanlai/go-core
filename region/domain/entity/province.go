package entity

import ()

type Province struct {
	ProvId     int     `json:"prov_id"`
	ProvName   *string `json:"prov_name"`
	Locationid *int    `json:"locationid"`
	Status     *int    `json:"status"`
}
