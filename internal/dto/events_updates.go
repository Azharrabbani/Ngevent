package dto

import (
	"errors"
	"ngevent/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GetUpdateReq struct {
	Status  string `json:"status" query:"status" validate:"required,oneof=pending approved rejected canceled"`
	EventID string `json:"event_id" validate:"required"`
	UserID  string `json:"user_id" validate:"required"`
	Role    string `json:"role" validate:"required"`
}

type UpdatedEventFilterReq struct {
	Title       string `json:"title"  query:"title"`
	Search      string `json:"search" query:"search"`
	Sort        string `json:"sort"   query:"sort"`
	WithDeleted bool   `json:"with_deleted" query:"with_deleted"`
	Date        string `json:"date"   query:"date"`
	StartTime   int64  `json:"start_time" query:"start_time"`
	Status      string `json:"status" query:"status" validate:"omitempty,oneof=pending approved rejected canceled"`
}

type UpdatedEventFilter struct {
	EventID     *string    `json:"event_id"`
	Title       *string    `json:"title"`
	WithDeleted *bool      `json:"with_deleted"`
	Search      *string    `json:"search"`
	Sort        *string    `json:"sort"`
	Date        *string    `json:"date"`
	Status      *string    `json:"status"`
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
}

type UpdatedEventRespReq struct {
	UpdatedEvent    *model.UpdatedEvents
	EventID         string
	EventCategories []EventCategories
	Tickets         []Tickets
	StartTime       int64
	EndTime         int64
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       *int64
}

type UpdateEvent struct {
	EventTx      *gorm.DB
	UpdatedEvent *model.UpdatedEvents
	Event        *model.Events
}

type ReviewUpdatedEventReq struct {
	ID         string  `json:"id"`
	Status     string  `json:"status" validate:"required,oneof=approved rejected"`
	Reason     *string `json:"reason"`
	ReviewedBy *string `json:"reviewed_by"`
}

func (r *ReviewUpdatedEventReq) ValidateReason() error {
	if r.Status == "rejected" && (r.Reason == nil || strings.TrimSpace(*r.Reason) == "") {
		return errors.New("reason is required when rejecting an update request")
	}
	return nil
}

type EventsUpdatesResp struct {
	ID                string            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID           string            `json:"event_id"`
	EventTitle        string            `json:"event_title"`
	EOProfile         EOProfiles        `json:"eo_profile"`
	UpdatedDetails    UpdatedDetails    `json:"updated_details"`
	UpdatedAddress    UpdatedAddress    `json:"updated_address"`
	UpdatedCategories []EventCategories `json:"updated_categories"`
	CreatedAt         int64             `json:"created_at" gorm:"default:now()"`
	UpdatedAt         int64             `json:"updated_at" gorm:"default:now()"`
	DeletedAt         *int64            `json:"deleted_at,omitempty"`
}

type EventUpdatesResp struct {
	ID                string            `json:"id"`
	EventID           string            `json:"event_id"`
	EventTitle        string            `json:"event_title"`
	EOProfile         EOProfiles        `json:"eo_profile"`
	UpdatedDetails    UpdatedDetails    `json:"updated_details"`
	UpdatedAddress    UpdatedAddress    `json:"updated_address"`
	UpdatedCategories []EventCategories `json:"updated_categories"`
	CreatedAt         int64             `json:"created_at"`
	UpdatedAt         int64             `json:"updated_at"`
	DeletedAt         *int64            `json:"deleted_at,omitempty"`
}

type UpdatedDetails struct {
	Banner         *string   `json:"banner,omitempty"`
	Status         string    `json:"status"`
	Description    string    `json:"description"`
	StartTime      int64     `json:"start_time"`
	EndTime        int64     `json:"end_time"`
	RejectedReason *string   `json:"rejected_reason,omitempty"`
	ReviewedBy     *Reviewer `json:"reviewed_by,omitempty"`
	ReviewedAt     *int64    `json:"reviewed_at,omitempty"`
}

type UpdatedAddress struct {
	Address       string      `json:"address"`
	City          string      `json:"city"`
	Country       string      `json:"country"`
	DetailAddress string      `json:"detail_address"`
	Coordinates   Coordinates `json:"coordinates"`
}

func ToEventUpdateResp(req *UpdatedEventRespReq) (*EventUpdatesResp, error) {
	if req.UpdatedEvent == nil {
		return nil, errors.New("updated event is nil")
	}

	var reviewer *Reviewer
	if req.UpdatedEvent.Reviewer != nil {
		reviewer = &Reviewer{
			ID:    req.UpdatedEvent.Reviewer.ID,
			Email: req.UpdatedEvent.Reviewer.Email,
		}
	}

	profile := req.UpdatedEvent.Event.Profile

	resp := &EventUpdatesResp{
		ID:         req.UpdatedEvent.ID,
		EventID:    req.UpdatedEvent.EventID,
		EventTitle: req.UpdatedEvent.Name,

		EOProfile: EOProfiles{
			ID:           profile.ID,
			IsVerified:   profile.User.IsVerified,
			Email:        profile.User.Email,
			Name:         profile.Name,
			PhotoProfile: profile.PhotoProfile,
			PhoneNumber:  profile.PhoneNumber,
		},

		UpdatedDetails: UpdatedDetails{
			Banner:         req.UpdatedEvent.Banner,
			Status:         req.UpdatedEvent.Status,
			Description:    req.UpdatedEvent.Description,
			StartTime:      req.StartTime,
			EndTime:        req.EndTime,
			RejectedReason: req.UpdatedEvent.RejectedReason,
			ReviewedBy:     reviewer,
			ReviewedAt:     timePtrToUnix(req.UpdatedEvent.ReviewedAt),
		},

		UpdatedAddress: UpdatedAddress{
			Address:       req.UpdatedEvent.Address,
			City:          req.UpdatedEvent.City,
			Country:       req.UpdatedEvent.Country,
			DetailAddress: req.UpdatedEvent.DetailAddress,
			Coordinates: Coordinates{
				Lat: req.UpdatedEvent.Lat,
				Lon: req.UpdatedEvent.Lon,
			},
		},

		UpdatedCategories: req.EventCategories,
		CreatedAt:         req.CreatedAt,
		UpdatedAt:         req.UpdatedAt,
		DeletedAt:         req.DeletedAt,
	}

	return resp, nil
}
