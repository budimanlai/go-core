package entity

import "time"

type UserSession struct {
	ID           int
	AppID        int
	UserID       uint
	Tokens       string
	CreateOn     time.Time
	LastAccessOn *time.Time
	RemoveOn     *time.Time
	FromIP       string
	UserAgent    string
}
