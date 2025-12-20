package model

import ()

type CityModel struct {
	CityId   int     `gorm:"column:city_id;primaryKey"`
	CityName *string `gorm:"column:city_name"`
	ProvId   int     `gorm:"column:prov_id"`
}

func (CityModel) TableName() string {
	return "cities"
}
