package entity

import ()

type City struct {
	CityId   int     `json:"city_id"`
	CityName *string `json:"city_name"`
	ProvId   int     `json:"prov_id"`
}
