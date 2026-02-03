package model

import "time"

type AttendeeProfiles struct {
	ID           string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Username     *string   `json:"username"`
	PhotoProfile *string   `json:"photo_profile"`
	PhoneNumber  string    `json:"phone_number"`
	Country      string    `json:"country"`
	Address      *string   `json:"address"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:now()"`
	User         Users     `gorm:"foreignKey:UserID"`
}

func (AttendeeProfiles) TableName() string {
	return "attendee_profiles"
}
