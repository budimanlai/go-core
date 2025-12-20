package model

import ()

type ProvinceModel struct {
	ProvId     int     `gorm:"column:prov_id;primaryKey"`
	ProvName   *string `gorm:"column:prov_name"`
	Locationid *int    `gorm:"column:locationid"`
	Status     *int    `gorm:"column:status"`
}

func (ProvinceModel) TableName() string {
	return "provinces"
}
