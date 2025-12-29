package models

import (
	"time"
)

type UserSession struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement"`
	AppID        int        `gorm:"column:app_id;not null"`
	UserID       int        `gorm:"column:user_id;not null"`
	Tokens       string     `gorm:"column:tokens;type:varchar(32);default:'';not null"`
	CreateOn     time.Time  `gorm:"column:create_on;default:CURRENT_TIMESTAMP;not null"`
	LastAccessOn *time.Time `gorm:"column:last_access_on"`
	RemoveOn     *time.Time `gorm:"column:remove_on"`
	FromIP       string     `gorm:"column:from_ip;type:varchar(15);default:'';not null"`
	UserAgent    string     `gorm:"column:user_agent;type:varchar(256)"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
