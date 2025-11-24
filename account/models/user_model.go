package models

import (
	"time"
)

type User struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	Username           string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	AuthKey            string     `gorm:"type:varchar(32);not null"`
	PasswordHash       string     `gorm:"type:varchar(255);not null"`
	PasswordResetToken *string    `gorm:"type:varchar(255);uniqueIndex"`
	Email              string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Fullname           string     `gorm:"type:varchar(225);not null"`
	Handphone          string     `gorm:"type:varchar(15);not null"`
	Dob                *time.Time `gorm:"type:date"`
	Gender             string     `gorm:"type:enum('M','F');default:'M'"`
	Status             string     `gorm:"type:enum('active','inactive','suspended');default:'active'"`
	MainRole           *string    `gorm:"type:varchar(225)"`
	LoginDashboard     string     `gorm:"type:enum('Y','N');default:'N'"`
	Avatar             *string    `gorm:"type:varchar(200)"`
	Address            *string    `gorm:"type:text"`
	Zipcode            string     `gorm:"type:varchar(10);not null"`
	DistrictID         uint       `gorm:"not null"`
	SubdistrictID      uint       `gorm:"not null"`
	CityID             uint       `gorm:"not null"`
	ProvinceID         uint       `gorm:"not null"`
	CountryID          string     `gorm:"type:char(2);not null"`
	CreatedAt          time.Time  `gorm:"autoCreateTime;not null"`
	CreatedBy          uint       `gorm:"not null"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime;not null"`
	UpdatedBy          uint       `gorm:"not null"`
	VerificationToken  *string    `gorm:"type:varchar(255);uniqueIndex"`
}

func (User) TableName() string {
	return "user"
}
