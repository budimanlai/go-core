package entity

import "time"

type User struct {
	ID           uint
	Username     string
	Fullname     string
	AuthKey      string
	PasswordHash string
	Email        string
	Handphone    string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    uint
	UpdatedBy    uint
}

func (u *User) IsActive() bool {
	return u.Status == "active"
}
