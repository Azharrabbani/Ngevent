package model

import "time"

type Events struct {
	ID            string            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID     string            `json:"profile_id"`
	Banner        string            `json:"banner"`
	Name          string            `json:"name"`
	Categories    []EventCategories `json:"categories"`
	Tickets       []Tickets         `json:"tickets" gorm:"foreignKey:EventID"`
	Slug          string            `json:"slug"`
	Status        string            `json:"status"`
	Description   string            `json:"description"`
	Address       string            `json:"address"`
	City          string            `json:"city"`
	Country       string            `json:"country"`
	DetailAddress string            `json:"detail_address"`
	Coordinates   string            `json:"coordinates" gorm:"type:geography(Point,4326)"`
	Date          time.Time         `json:"date"`
	CreatedAt     time.Time         `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"default:now()"`
	DeletedAt     time.Time         `json:"deleted_at"`
	Profile       OrganizerProfiles `gorm:"foreignKey:ProfileID"`
}

type EventCategories struct {
	ID         string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID    string     `json:"event_id"`
	CategoryID int64      `json:"category_id"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt  time.Time  `json:"deleted_at"`
	Category   Categories `gorm:"foreignKey:CategoryID"`
}

type Tickets struct {
	ID         string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID    string    `json:"event_id"`
	Name       string    `json:"name"`
	Price      string    `json:"price"`
	Quantity   int       `json:"quantity"`
	TicketType string    `json:"ticket_type"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:now()"`
	DeletedAt  time.Time `json:"deleted_at"`
}

func (Events) TableName() string {
	return "events"
}

func (EventCategories) TableName() string {
	return "event_categories"
}
