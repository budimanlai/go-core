package models

import "time"

type Otp struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Handphone string    `gorm:"column:handphone;type:varchar(50);default:''"`
	TrxID     string    `gorm:"column:trx_id;type:varchar(50);default:''"`
	PinCode   string    `gorm:"column:pin_code;type:char(6);default:''"`
	Status    string    `gorm:"column:status;type:varchar(15);default:'waiting'"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (Otp) TableName() string {
	return "wa_otp"
}
