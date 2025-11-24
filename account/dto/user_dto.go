package dto

import (
	"time"

	"github.com/budimanlai/go-pkg/types"
)

type RegisterRequest struct {
	Username      string     `json:"username" validate:"required,min=3,max=50"`
	Email         string     `json:"email" validate:"required,email"`
	Password      string     `json:"password" validate:"required,min=6"`
	Fullname      string     `json:"fullname" validate:"required,min=2,max=100"`
	Handphone     string     `json:"handphone" validate:"required,min=10,max=15"`
	Dob           *time.Time `json:"dob"`
	Gender        string     `json:"gender" validate:"omitempty,oneof=M F"`
	Address       *string    `json:"address"`
	Zipcode       string     `json:"zipcode"`
	DistrictID    uint       `json:"district_id"`
	SubdistrictID uint       `json:"subdistrict_id"`
	CityID        uint       `json:"city_id"`
	ProvinceID    uint       `json:"province_id"`
	CountryID     string     `json:"country_id" validate:"required,len=2"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type UpdateUserRequest struct {
	Fullname      *string    `json:"fullname"`
	Handphone     *string    `json:"handphone"`
	Dob           *time.Time `json:"dob"`
	Gender        *string    `json:"gender" validate:"omitempty,oneof=M F"`
	Address       *string    `json:"address"`
	Zipcode       *string    `json:"zipcode"`
	DistrictID    *uint      `json:"district_id"`
	SubdistrictID *uint      `json:"subdistrict_id"`
	CityID        *uint      `json:"city_id"`
	ProvinceID    *uint      `json:"province_id"`
	CountryID     *string    `json:"country_id" validate:"omitempty,len=2"`
	Avatar        *string    `json:"avatar"`
}

type UserResponse struct {
	ID             uint          `json:"id"`
	Username       string        `json:"username"`
	Email          string        `json:"email"`
	Fullname       string        `json:"fullname"`
	Handphone      string        `json:"handphone"`
	Dob            *time.Time    `json:"dob,omitempty"`
	Gender         string        `json:"gender"`
	Status         string        `json:"status"`
	MainRole       *string       `json:"main_role,omitempty"`
	LoginDashboard string        `json:"login_dashboard"`
	Avatar         *string       `json:"avatar,omitempty"`
	Address        *string       `json:"address,omitempty"`
	Zipcode        string        `json:"zipcode"`
	DistrictID     uint          `json:"district_id"`
	SubdistrictID  uint          `json:"subdistrict_id"`
	CityID         uint          `json:"city_id"`
	ProvinceID     uint          `json:"province_id"`
	CountryID      string        `json:"country_id"`
	CreatedAt      types.UTCTime `json:"created_at"`
	UpdatedAt      types.UTCTime `json:"updated_at"`
}

type ListUserResponse struct {
	Users      []*UserResponse `json:"users"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
