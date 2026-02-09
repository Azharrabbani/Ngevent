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

type OrganizerProfiles struct {
	ID            string               `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        string               `json:"user_id"`
	IsVerified    bool                 `json:"is_verified"`
	Name          string               `json:"name"`
	PhotoProfile  *string              `json:"photo_profile"`
	PhoneNumber   string               `json:"phone_number"`
	Country       string               `json:"country"`
	Address       *string              `json:"address"`
	SocialMedias  OrganizerSocialMedia `json:"social_media"`
	CompanyDetail OrganizerCompDetail  `json:"company_detail"`
	CreatedAt     time.Time            `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time            `json:"updated_at" gorm:"default:now()"`
	User          Users                `gorm:"foreignKey:UserID"`
}

type OrganizerSocialMedia struct {
	Email     *string `json:"email"`
	Instagram *string `json:"instagram"`
}

type OrganizerCompDetail struct {
	Description *string `json:"description"`
	NPWP        string  `json:"npwp"`
	NIB         string  `json:"nib"`
}

func (AttendeeProfiles) TableName() string {
	return "attendee_profiles"
}

func (OrganizerProfiles) TableName() string {
	return "organizer_profiles"
}
