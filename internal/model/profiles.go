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
	Status        OrganizerStatus      `json:"status" gorm:"embedded"`
	Name          string               `json:"name"`
	PhotoProfile  *string              `json:"photo_profile"`
	PhoneNumber   string               `json:"phone_number"`
	Country       string               `json:"country"`
	Address       *string              `json:"address"`
	SocialMedias  OrganizerSocialMedia `json:"social_media" gorm:"embedded"`
	CompanyDetail OrganizerCompDetail  `json:"company_detail" gorm:"embedded"`
	CreatedAt     time.Time            `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time            `json:"updated_at" gorm:"default:now()"`
	User          Users                `gorm:"foreignKey:UserID"`
}

type OrganizerStatus struct {
	Status         string     `json:"status" gorm:"type:organizer_profile_status;default:'pending'"`
	RejectedReason *string    `json:"rejected_reason"`
	ReviewedBy     *string    `json:"reviewed_by"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
}

type OrganizerSocialMedia struct {
	Email     *string `json:"email"`
	Instagram *string `json:"instagram"`
}

type OrganizerCompDetail struct {
	Description  *string `json:"description"`
	NPWPNumber   string  `json:"npwp_number"`
	NPWPDocument string  `json:"npwp_document"`
	NIBNumber    string  `json:"nib"`
	NIBDocument  string  `json:"nib_document"`
}

func (AttendeeProfiles) TableName() string {
	return "attendee_profiles"
}

func (OrganizerProfiles) TableName() string {
	return "organizer_profiles"
}
