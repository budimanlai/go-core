package models

type District struct {
	DisID   uint   `gorm:"column:dis_id;primaryKey;autoIncrement" json:"dis_id"`
	DisName string `gorm:"column:dis_name;type:varchar(255)" json:"dis_name"`
	CityID  uint   `gorm:"column:city_id;not null" json:"city_id"`
}

func (District) TableName() string {
	return "districts"
}
