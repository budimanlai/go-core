package model

import ()

type DistrictModel struct {
	DisId   int     `gorm:"column:dis_id;primaryKey"`
	DisName *string `gorm:"column:dis_name"`
	CityId  int     `gorm:"column:city_id"`
}

func (DistrictModel) TableName() string {
	return "districts"
}
