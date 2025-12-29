package entity

type User struct {
	ID           uint
	Username     string
	Fullname     string
	PasswordHash string
	Email        string
	Handphone    string
	Status       string
}

func (u *User) IsActive() bool {
	return u.Status == "active"
}
