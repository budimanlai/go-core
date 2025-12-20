package model

import ()

type SubdistrictModel struct {
	SubdisId   int     `gorm:"column:subdis_id;primaryKey"`
	SubdisName *string `gorm:"column:subdis_name"`
	DisId      int     `gorm:"column:dis_id"`
}

func (SubdistrictModel) TableName() string {
	return "subdistricts"
}
