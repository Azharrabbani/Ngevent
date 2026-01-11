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

type RegisterResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type LoginResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AccessToken string `json:"access_token"`
	LoginAt     int64  `json:"login_at"`
}

func (Users) TableName() string {
	return "users"
}
