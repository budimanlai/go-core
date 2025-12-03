package models

type City struct {
	CityID   uint   `gorm:"column:city_id;primaryKey;autoIncrement"`
	CityName string `gorm:"column:city_name;type:varchar(255)"`
	ProvID   uint   `gorm:"column:prov_id;not null"`
}

func (City) TableName() string {
	return "cities"
}
