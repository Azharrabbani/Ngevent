package model

import (
	"time"
)

type Status string

const (
	Draft     Status = "draft"
	Pending   Status = "pending"
	Active    Status = "active"
	Done      Status = "done"
	Rejected  Status = "rejected"
	Cancelled Status = "cancelled"
	Approved  Status = "approved"
)

type Events struct {
	ID             string             `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID      string             `json:"profile_id"`
	Banner         *string            `json:"banner"`
	Name           string             `json:"name"`
	Categories     []*EventCategories `json:"categories" gorm:"foreignKey:EventID"`
	Slug           string             `json:"slug"`
	Status         string             `json:"status"`
	RequestUpdates bool               `json:"request_updates" gorm:"default:false"`
	Description    string             `json:"description"`
	Address        string             `json:"address"`
	City           string             `json:"city"`
	Country        string             `json:"country"`
	DetailAddress  string             `json:"detail_address"`
	Coordinates    string             `json:"coordinates" gorm:"type:geography(Point,4326)"`
	Lat            float64            `json:"lat" gorm:"->;column:lat"`
	Lon            float64            `json:"lon" gorm:"->;column:lon"`
	StartTime      time.Time          `json:"start_time"`
	EndTime        time.Time          `json:"end_time"`
	RejectedReason *string            `json:"rejected_reason"`
	ReviewedBy     *string            `json:"reviewed_by"`
	ReviewedAt     *time.Time         `json:"reviewed_at"`
	CreatedAt      time.Time          `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time          `json:"updated_at" gorm:"default:now()"`
	SubmittedAt    *time.Time         `json:"submitted_at"`
	DeletedAt      *time.Time         `json:"deleted_at"`
	Profile        OrganizerProfiles  `gorm:"foreignKey:ProfileID;references:ID"`
	Reviewer       *Users             `gorm:"foreignKey:ReviewedBy;references:ID"`
}

type EventCategories struct {
	ID         string      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID    string      `json:"event_id"`
	CategoryID int64       `json:"category_id"`
	CreatedAt  time.Time   `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"default:now()"`
	DeletedAt  *time.Time  `json:"deleted_at"`
	Category   *Categories `gorm:"foreignKey:CategoryID;references:ID"`
}

type Tickets struct {
	ID         string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID    string     `json:"event_id"`
	Event      *Events    `gorm:"foreignKey:EventID"`
	Name       string     `json:"name"`
	Price      string     `json:"price"`
	Quantity   int        `json:"quantity"`
	TicketType string     `json:"ticket_type"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type UpdatedEvents struct {
	ID             string                   `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID        string                   `json:"event_id"`
	Name           string                   `json:"name"`
	Banner         *string                  `json:"banner"`
	Categories     []*EventCategoriesUpdate `json:"categories" gorm:"foreignKey:EventUpdateID"`
	Slug           string                   `json:"slug"`
	Status         string                   `json:"string"`
	Description    string                   `json:"description"`
	Address        string                   `json:"address"`
	City           string                   `json:"city"`
	Country        string                   `json:"country"`
	DetailAddress  string                   `json:"detail_address"`
	Coordinates    string                   `json:"coordinates" gorm:"type:geography(Point,4326)"`
	Lat            float64                  `json:"lat" gorm:"->;column:lat"`
	Lon            float64                  `json:"lon" gorm:"->;column:lon"`
	StartTime      time.Time                `json:"start_time"`
	EndTime        time.Time                `json:"end_time"`
	RejectedReason *string                  `json:"rejected_reason"`
	ReviewedBy     *string                  `json:"reviewed_by"`
	ReviewedAt     *time.Time               `json:"reviewed_at"`
	CreatedAt      time.Time                `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time                `json:"updated_at" gorm:"default:now()"`
	SubmittedAt    time.Time                `json:"submitted_at" gorm:"default:now()"`
	DeletedAt      *time.Time               `json:"deleted_at"`
	Event          Events                   `gorm:"foreignKey:EventID"`
	Reviewer       *Users                   `gorm:"foreignKey:ReviewedBy;references:ID"`
}

type EventCategoriesUpdate struct {
	ID            string      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventUpdateID string      `json:"event_update_id"`
	CategoryID    int64       `json:"category_id"`
	CreatedAt     time.Time   `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"default:now()"`
	DeletedAt     time.Time   `json:"deleted_at"`
	Category      *Categories `gorm:"foreignKey:CategoryID;references:ID"`
}

type TicketsUpdate struct {
	ID            string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventUpdateID string    `json:"event_update_id" gorm:"foreignKey:EventUpdateID"`
	Name          string    `json:"name"`
	Price         string    `json:"price"`
	Quantity      int       `json:"quantity"`
	TicketType    string    `json:"ticket_type"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:now()"`
	DeletedAt     time.Time `json:"deleted_at"`
}

func (Events) TableName() string {
	return "events"
}

func (EventCategories) TableName() string {
	return "event_categories"
}

func (Tickets) TableName() string {
	return "tickets"
}

func (UpdatedEvents) TableName() string {
	return "event_updates"
}

func (EventCategoriesUpdate) TableName() string {
	return "event_update_categories"
}

func (TicketsUpdate) TableName() string {
	return "event_update_tickets"
}
