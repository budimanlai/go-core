package models

import "time"

type AccountModel struct {
	ID        string     `db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string     `db:"email" gorm:"uniqueIndex;not null;type:varchar(255)"`
	Username  string     `db:"username" gorm:"uniqueIndex;not null;type:varchar(50)"`
	Password  string     `db:"password" gorm:"not null;type:varchar(255)"`
	FullName  string     `db:"full_name" gorm:"not null;type:varchar(100)"`
	Role      string     `db:"role" gorm:"not null;type:varchar(20);default:'user'"`
	IsActive  bool       `db:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt *time.Time `db:"deleted_at" gorm:"index"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
