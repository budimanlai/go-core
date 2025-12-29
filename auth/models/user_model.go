package models

type User struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string `gorm:"column:username;type:varchar(255);not null"`
	Fullname     string `gorm:"column:fullname;type:varchar(255);not null"`
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`
	Email        string `gorm:"column:email;type:varchar(255);not null;uniqueIndex"`
	Handphone    string `gorm:"column:handphone;type:varchar(20);not null;default:'';uniqueIndex"`
	Status       string `gorm:"column:status;type:varchar(15);not null;default:'active'"`
}

func (User) TableName() string {
	return "users"
}
