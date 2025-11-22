package entity

import (
	"time"
)

type User struct {
	ID                   uint
	Username             string
	AuthKey              string
	PasswordHash         string
	PasswordResetToken   *string
	Email                string
	Fullname             string
	Handphone            string
	Dob                  *time.Time
	Gender               string
	Status               string
	MainRole             *string
	LoginDashboard       string
	Avatar               *string
	Address              *string
	Zipcode              string
	DistrictID           uint
	SubdistrictID        uint
	CityID               uint
	ProvinceID           uint
	CountryID            string
	CreatedAt            time.Time
	CreatedBy            uint
	UpdatedAt            time.Time
	UpdatedBy            uint
	VerificationToken    *string
	DeletedAt            *time.Time
}

// IsDeleted checks if the user has been soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// Activate changes user status to active
func (u *User) Activate() {
	u.Status = "active"
	u.UpdatedAt = time.Now()
}

// Deactivate changes user status to inactive
func (u *User) Deactivate() {
	u.Status = "inactive"
	u.UpdatedAt = time.Now()
}

// Suspend changes user status to suspended
func (u *User) Suspend() {
	u.Status = "suspended"
	u.UpdatedAt = time.Now()
}

// SoftDelete marks user as deleted
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// CanLoginDashboard checks if user can access dashboard
func (u *User) CanLoginDashboard() bool {
	return u.LoginDashboard == "Y"
}

// EnableDashboardAccess enables dashboard login
func (u *User) EnableDashboardAccess() {
	u.LoginDashboard = "Y"
	u.UpdatedAt = time.Now()
}

// DisableDashboardAccess disables dashboard login
func (u *User) DisableDashboardAccess() {
	u.LoginDashboard = "N"
	u.UpdatedAt = time.Now()
}

// IsVerified checks if user email is verified
func (u *User) IsVerified() bool {
	return u.VerificationToken == nil
}

// ClearVerificationToken marks email as verified
func (u *User) ClearVerificationToken() {
	u.VerificationToken = nil
	u.UpdatedAt = time.Now()
}
