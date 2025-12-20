package entity

import ()

type District struct {
	DisId   int     `json:"dis_id"`
	DisName *string `json:"dis_name"`
	CityId  int     `json:"city_id"`
}
