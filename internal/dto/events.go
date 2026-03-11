package dto

import (
	"errors"
	"mime/multipart"
	"ngevent/internal/model"
)

type CreateEventReq struct {
	Banner multipart.FileHeader
	Name   string `json:"name" validate:"required"`
}

type EventReq struct {
	Name        string          `json:"name" validate:"required"`
	UserID      string          `json:"user_id"`
	EventID     string          `json:"event_id"`
	Description string          `json:"description" validate:"required"`
	Categories  []int64         `json:"categories" validate:"required,min=1"`
	Tickets     []TicketsReq    `json:"tickets" validate:"required,min=1"`
	Date        int64           `json:"date" validate:"required"`
	Address     EventAddressReq `json:"address" validate:"required"`
	Status      *string         `json:"status"`
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

type ReviewEventReq struct {
	ID     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type UpdateEventReq struct {
	UserID string
	ID     string
	Banner *multipart.FileHeader
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

type EventRespReq struct {
	Event           *model.Events
	Organizer       *model.OrganizerProfiles
	EventCategories []EventCategories
	Tickets         []Tickets
	Date            int64
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       int64
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

func ToEventResp(req *EventRespReq) (*EventsResp, error) {
	if req.Event == nil {
		return nil, errors.New("event is nil")
	}

	eventResp := &EventsResp{
		ID: req.Event.ID,
		EOProfile: EOProfiles{
			ID:           req.Organizer.ID,
			IsVerified:   req.Organizer.User.IsVerified,
			Email:        req.Organizer.User.Email,
			Name:         req.Organizer.Name,
			PhotoProfile: req.Organizer.PhotoProfile,
			PhoneNumber:  req.Organizer.PhoneNumber,
		},
		Event: EventDetail{
			Banner:      *req.Event.Banner,
			Name:        req.Event.Name,
			Categories:  req.EventCategories,
			Tickets:     req.Tickets,
			Slug:        req.Event.Slug,
			Status:      req.Event.Status,
			Description: req.Event.Description,
		},
		EventAddress: EventAddress{
			Address:       req.Event.Address,
			City:          req.Event.City,
			Country:       req.Event.Country,
			DetailAddress: req.Event.DetailAddress,
			Coordinates:   req.Event.Coordinates,
		},
		Date:      req.Date,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		DeletedAt: req.DeletedAt,
	}

	return eventResp, nil
}
