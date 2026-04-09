package model

import "time"

type UserRole string

const (
	Admin     UserRole = "admin"
	Attendee  UserRole = "user"
	Organizer UserRole = "event organizer"
)

type Users struct {
	ID         string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Role       string     `json:"role"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type RegisterResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (Users) TableName() string {
	return "users"
}
