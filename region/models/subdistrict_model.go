package models

type SubDistrict struct {
	SubdisID   uint   `gorm:"primaryKey;column:subdis_id;type:int(11) unsigned;autoIncrement" json:"subdis_id"`
	SubdisName string `gorm:"column:subdis_name;type:varchar(255)" json:"subdis_name"`
	DisID      uint   `gorm:"column:dis_id;type:int(11) unsigned;not null" json:"dis_id"`
}

func (SubDistrict) TableName() string {
	return "subdistricts"
}
