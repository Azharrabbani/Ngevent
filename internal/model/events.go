package model

import "time"

type Status int

const (
	Draft     Status = 1
	Pending   Status = 2
	Active    Status = 3
	Done      Status = 4
	Rejected  Status = 5
	Cancelled Status = 6
)

type EventsStatuses struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	Status    string     `json:"name"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:now()"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Events struct {
	ID            string             `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID     string             `json:"profile_id"`
	Banner        *string            `json:"banner"`
	Name          string             `json:"name"`
	Categories    []*EventCategories `json:"categories" gorm:"foreignKey:EventID"`
	Slug          string             `json:"slug"`
	StatusID      int64              `json:"status_id"`
	Description   string             `json:"description"`
	Address       string             `json:"address"`
	City          string             `json:"city"`
	Country       string             `json:"country"`
	DetailAddress string             `json:"detail_address"`
	Coordinates   string             `json:"coordinates" gorm:"type:geography(Point,4326)"`
	Lat           float64            `json:"lat" gorm:"->;column:lat"`
	Lon           float64            `json:"lon" gorm:"->;column:lon"`
	StartTime     time.Time          `json:"start_time"`
	EndTime       time.Time          `json:"end_time"`
	CreatedAt     time.Time          `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time          `json:"updated_at" gorm:"default:now()"`
	DeletedAt     *time.Time         `json:"deleted_at"`
	Profile       OrganizerProfiles  `gorm:"foreignKey:ProfileID;references:ID"`
	Status        EventsStatuses     `gorm:"foreignKey:StatusID;references:ID"`
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
	ID            string                   `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID       string                   `json:"event_id"`
	Name          string                   `json:"name"`
	Banner        *string                  `json:"banner"`
	Categories    []*EventCategoriesUpdate `json:"categories" gorm:"foreignKey:EventUpdateID"`
	Slug          string                   `json:"slug"`
	StatusID      int64                    `json:"status_id"`
	Description   string                   `json:"description"`
	Address       string                   `json:"address"`
	City          string                   `json:"city"`
	Country       string                   `json:"country"`
	DetailAddress string                   `json:"detail_address"`
	Coordinates   string                   `json:"coordinates" gorm:"type:geography(Point,4326)"`
	Lat           float64                  `json:"lat" gorm:"column:lat"`
	Lon           float64                  `json:"lon" gorm:"column:lon"`
	StartTime     time.Time                `json:"start_time"`
	EndTime       time.Time                `json:"end_time"`
	CreatedAt     time.Time                `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time                `json:"updated_at" gorm:"default:now()"`
	DeletedAt     *time.Time               `json:"deleted_at"`
	Event         Events                   `gorm:"foreignKey:EventID"`
	Status        EventsStatuses           `gorm:"foreignKey:StatusID"`
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

func (EventsStatuses) TableName() string {
	return "events_statuses"
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
