package dto

import (
	"errors"
	"mime/multipart"
	"ngevent/internal/model"
	"time"
)

type CreateEventReq struct {
	Banner multipart.FileHeader
	Name   string `json:"name" validate:"required"`
}

type EventReq struct {
	ID          *string         `json:"id"`
	Name        string          `json:"name" validate:"required"`
	UserID      string          `json:"user_id"`
	Description string          `json:"description" validate:"required"`
	Categories  []int64         `json:"categories" validate:"required,min=1"`
	StartTime   int64           `json:"start_time" validate:"required"`
	EndTime     int64           `json:"end_time" validate:"required"`
	Address     EventAddressReq `json:"address" validate:"required"`
	Status      string          `json:"status" validate:"oneof=draft pending"`
}

type NearestEventReq struct {
	Lat float64 `json:"lat" validate:"required"`
	Lon float64 `json:"lon" validate:"required"`
}

type NearestResult struct {
	Haversine Haversine   `json:"haversine"`
	Dijkstra  Dijkstra    `json:"dijkstra"`
	Path      []PathPoint `json:"Path"`
}

type Haversine struct {
	Name     string
	Distance string
	Time     string
	Accuracy string
}

type Dijkstra struct {
	Name     string
	Distance string
	Time     string
	Accuracy string
}

type PathPoint struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type RouteResp struct {
	Event    string      `json:"event"`
	Distance string      `json:"distance"`
	Path     []PathPoint `json:"path"`
}

type EventFilterReq struct {
	Title       string `json:"title" query:"title"`
	Search      string `json:"search" query:"search"` // global search
	Sort        string `json:"sort" query:"sort"`
	Date        string `json:"date" query:"date"`
	Category    []int  `json:"category" query:"category"`
	Status      string `json:"status" query:"status"`
	GetUpdate   *bool  `json:"get_update" query:"get_update"`
	WithDeleted bool   `json:"with_deleted" query:"with_deleted"`
	StartTime   int64  `json:"start_time" query:"start_time"`
	Location    string `json:"location" query:"location"`
}

type EventFilter struct {
	ProfileID   *string    `json:"profile_id"`
	Title       *string    `json:"title" query:"title"`
	Search      *string    `json:"search" query:"search"` // global search
	GetUpdate   *bool      `json:"get_update"`
	Sort        *string    `json:"sort"`
	Date        *string    `json:"date"`
	Category    []int      `json:"category" query:"category"`
	Status      *string    `json:"status" query:"status"`
	WithDeleted *bool      `json:"with_deleted" query:"with_deleted"`
	Start       *time.Time `json:"start" query:"start"`
	End         *time.Time `json:"end" query:"end"`
	Location    *string    `json:"location" query:"location"`
}

type EventAddressReq struct {
	DetailAddress string `json:"detail_address" validate:"required"`
	Lat           string `json:"lat" validate:"required"`
	Long          string `json:"long" validate:"required"`
}

type ReviewEventReq struct {
	ID     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=active reject"`
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
	StartTime    int64        `json:"start_time"`
	EndTime      int64        `json:"end_time"`
	Distance     string       `json:"distance,omitempty"`
	Path         []PathPoint  `json:"path,omitempty"`
	CreatedAt    int64        `json:"created_at" gorm:"default:now()"`
	UpdatedAt    int64        `json:"updated_at" gorm:"default:now()"`
	DeletedAt    *int64       `json:"deleted_at,omitempty"`
}

type EventRespReq struct {
	Event           *model.Events
	Organizer       *model.OrganizerProfiles
	EventCategories []EventCategories
	StartTime       int64
	EndTime         int64
	UserLat         float64
	UserLon         float64
	Path            []PathPoint
	Distance        string
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       *int64
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
	Banner      *string           `json:"banner,omitempty"`
	Name        string            `json:"name"`
	Categories  []EventCategories `json:"categories"`
	Slug        string            `json:"slug"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
}

type EventAddress struct {
	Address       string      `json:"address"`
	City          string      `json:"city"`
	Country       string      `json:"country"`
	DetailAddress string      `json:"detail_address"`
	Coordinates   Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type EventCategories struct {
	ID   int64  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
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
			Banner:      req.Event.Banner,
			Name:        req.Event.Name,
			Categories:  req.EventCategories,
			Slug:        req.Event.Slug,
			Status:      req.Event.Status,
			Description: req.Event.Description,
		},
		EventAddress: EventAddress{
			Address:       req.Event.Address,
			City:          req.Event.City,
			Country:       req.Event.Country,
			DetailAddress: req.Event.DetailAddress,
			Coordinates: Coordinates{
				Lat: req.Event.Lat,
				Lon: req.Event.Lon,
			},
		},
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Distance:  req.Distance,
		Path:      req.Path,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		DeletedAt: req.DeletedAt,
	}

	return eventResp, nil
}
