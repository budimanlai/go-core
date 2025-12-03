package models

type Province struct {
	ProvID     uint   `gorm:"column:prov_id;primaryKey;autoIncrement" json:"prov_id"`
	ProvName   string `gorm:"column:prov_name;type:varchar(255)" json:"prov_name"`
	LocationID int    `gorm:"column:locationid" json:"locationid"`
	Status     int    `gorm:"column:status;default:1" json:"status"`
}

func (Province) TableName() string {
	return "provinces"
}
