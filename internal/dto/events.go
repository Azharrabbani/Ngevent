package dto

import "mime/multipart"

type CreateEventReq struct {
	Banner multipart.FileHeader
	Name   string `json:"name" validate:"required"`
}

type EventReq struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Categories  []int           `json:"categories" validate:"require,min=1"`
	Tickets     []TicketsReq    `json:"tickets" validate:"required,min=1"`
	Date        int64           `json:"date" validate:"required"`
	Address     EventAddressReq `json:"address" validate:"required"`
}

type EventsResp struct {
	ID           string       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EOProfile    EOProfiles   `json:"eo_profile"`
	Event        EventDetail  `json:"event"`
	EventAddress EventAddress `json:"event_address"`
	Date         int64        `json:"date"`
	CreatedAt    int64        `json:"created_at" gorm:"default:now()"`
	UpdatedAt    int64        `json:"updated_at" gorm:"default:now()"`
	DeletedAt    int64        `json:"deleted_at"`
}

type TicketsReq struct {
	Name       string `json:"name"`
	Price      string `json:"price"`
	Quantity   int    `json:"quantity"`
	TicketType string `json:"ticket_type"`
}

type EventAddressReq struct {
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	DetailAddress string `json:"detail_address"`
	Lat           string `json:"lat"`
	Long          string `json:"long"`
}

type EOProfiles struct {
	ID           string  `json:"id"`
	IsVerified   bool    `json:"is_verified"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	PhotoProfile *string `json:"photo_profile,omitempty"`
	PhoneNumber  string  `json:"phone_number"`
}

type EventDetail struct {
	Banner      string            `json:"banner"`
	Name        string            `json:"name"`
	Categories  []EventCategories `json:"categories"`
	Tickets     []Tickets         `json:"tickets"`
	Slug        string            `json:"slug"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
}

type EventAddress struct {
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	DetailAddress string `json:"detail_address"`
	Coordinates   string `json:"coordinates"`
}

type EventCategories struct {
	ID   string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string `json:"name"`
}

type Tickets struct {
	ID         string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name       string `json:"name"`
	Price      string `json:"price"`
	Quantity   int    `json:"quantity"`
	TicketType string `json:"ticket_type"`
}
