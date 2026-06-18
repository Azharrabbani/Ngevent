package dto

import (
	"errors"
	"mime/multipart"
	"ngevent/internal/model"
	"strings"
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
	StartDate   int64           `json:"start_date" validate:"required"`
	EndDate     int64           `json:"end_date" validate:"required"`
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
	Title     string   `json:"title" query:"title"`
	Search    string   `json:"search" query:"search"` // global search
	Sort      string   `json:"sort" query:"sort"`
	Date      string   `json:"date" query:"date"`
	Category  []int    `json:"category" query:"category"`
	Status    string   `json:"status" query:"status"`
	EventDate int64    `json:"event_date" query:"event_date"`
	Location  string   `json:"location" query:"location"`
	Month     int      `json:"month" query:"month"`
	Year      int      `json:"year" query:"year"`
	Lat       *float64 `json:"lat" query:"lat"`
	Lon       *float64 `json:"lon" query:"lon"`
}

type EventFilter struct {
	ProfileID  *string    `json:"profile_id"`
	Title      *string    `json:"title" query:"title"`
	Search     *string    `json:"search" query:"search"` // global search
	Sort       *string    `json:"sort"`
	Date       *string    `json:"date"`
	Role       *string    `json:"role" query:"role"`
	Category   []int      `json:"category" query:"category"`
	Status     *string    `json:"status" query:"status"`
	RangeStart *time.Time `json:"range_start" query:"range_start"`
	RangeEnd   *time.Time `json:"range_end" query:"range_end"`
	Location   *string    `json:"location" query:"location"`
	Month      *int       `json:"month"`
	Year       *int       `json:"year"`
	Lat        *float64   `json:"lat"`
	Lon        *float64   `json:"lon"`
}

type EventFilterPublic struct {
	ProfileID string  `json:"profile_id"`
	Title     *string `json:"title" query:"title"`
	Status    *string `json:"status" query:"status" validate:"oneof=active done"`
}

type EventAddressReq struct {
	DetailAddress string `json:"detail_address" validate:"required"`
	Lat           string `json:"lat" validate:"required"`
	Long          string `json:"long" validate:"required"`
	DisplayName   string `json:"display_name"`
	City          string `json:"city"`
	Country       string `json:"country"`
}

type ReviewEventReq struct {
	ID         string  `json:"id"`
	Status     string  `json:"status" validate:"required,oneof=active rejected"`
	Reason     *string `json:"reason"`
	ReviewedBy *string `json:"reviewed_by"`
}

func (r *ReviewEventReq) ValidateReason() error {
	if r.Status == "reject" && (r.Reason == nil || strings.TrimSpace(*r.Reason) == "") {
		return errors.New("reason is required when rejecting an event")
	}
	return nil
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
	StartDate    int64        `json:"start_date"`
	EndDate      int64        `json:"end_date"`
	StartTime    int64        `json:"start_time"`
	EndTime      int64        `json:"end_time"`
	Distance     string       `json:"distance,omitempty"`
	Path         []PathPoint  `json:"path,omitempty"`
	CreatedAt    int64        `json:"created_at" gorm:"default:now()"`
	UpdatedAt    int64        `json:"updated_at" gorm:"default:now()"`
	SubmittedAt  *int64       `json:"submitted_at,omitempty"`
	DeletedAt    *int64       `json:"deleted_at,omitempty"`
}

type EventRespReq struct {
	Event           *model.Events
	Organizer       *model.OrganizerProfiles
	EventCategories []EventCategories
	StartDate       int64
	EndDate         int64
	StartTime       int64
	EndTime         int64
	UserLat         float64
	UserLon         float64
	Path            []PathPoint
	Distance        string
	CreatedAt       int64
	UpdatedAt       int64
	SubmittedAt     *int64
	DeletedAt       *int64
}

type EOProfiles struct {
	ID           string  `json:"id"`
	IsVerified   bool    `json:"is_verified"`
	Status       string  `json:"status"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	PhotoProfile *string `json:"photo_profile,omitempty"`
	PhoneNumber  string  `json:"phone_number"`
}

type EventDetail struct {
	Banner         *string           `json:"banner,omitempty"`
	Name           string            `json:"name"`
	Categories     []EventCategories `json:"categories"`
	Slug           string            `json:"slug"`
	Status         string            `json:"status"`
	RequestUpdates bool              `json:"request_updates"`
	Description    string            `json:"description"`
	RejectedReason *string           `json:"rejected_reason,omitempty"`
	ReviewedBy     *Reviewer         `json:"reviewed_by,omitempty"`
	ReviewedAt     *int64            `json:"reviewed_at,omitempty"`
}

type Reviewer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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

	var reviewer *Reviewer
	if req.Event.Reviewer != nil {
		reviewer = &Reviewer{
			ID:    req.Event.Reviewer.ID,
			Email: req.Event.Reviewer.Email,
		}
	}

	eventResp := &EventsResp{
		ID: req.Event.ID,
		EOProfile: EOProfiles{
			ID:           req.Organizer.ID,
			IsVerified:   req.Organizer.User.IsVerified,
			Status:       req.Organizer.Status.Status,
			Email:        req.Organizer.User.Email,
			Name:         req.Organizer.Name,
			Slug:         req.Organizer.Slug,
			PhotoProfile: req.Organizer.PhotoProfile,
			PhoneNumber:  req.Organizer.PhoneNumber,
		},
		Event: EventDetail{
			Banner:         req.Event.Banner,
			Name:           req.Event.Name,
			Categories:     req.EventCategories,
			Slug:           req.Event.Slug,
			Status:         req.Event.Status,
			RequestUpdates: req.Event.RequestUpdates,
			Description:    req.Event.Description,
			RejectedReason: req.Event.RejectedReason,
			ReviewedBy:     reviewer,
			ReviewedAt:     timePtrToUnix(req.Event.ReviewedAt),
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
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Distance:    req.Distance,
		Path:        req.Path,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
		SubmittedAt: req.SubmittedAt,
		DeletedAt:   req.DeletedAt,
	}

	return eventResp, nil
}

func timePtrToUnix(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	val := t.Unix()
	return &val
}
