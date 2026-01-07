package model

import "time"

type Users struct {
	ID         string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:now()"`
}

func (Users) TableName() string {
	return "users"
}
