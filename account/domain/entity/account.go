package entity

import "time"

type Account struct {
	ID        string
	Email     string
	Username  string
	Password  string
	FullName  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (a *Account) IsDeleted() bool {
	return a.DeletedAt != nil
}

func (a *Account) Activate() {
	a.IsActive = true
	a.UpdatedAt = time.Now()
}

func (a *Account) Deactivate() {
	a.IsActive = false
	a.UpdatedAt = time.Now()
}

func (a *Account) SoftDelete() {
	now := time.Now()
	a.DeletedAt = &now
	a.UpdatedAt = now
}
