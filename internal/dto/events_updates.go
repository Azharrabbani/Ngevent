package dto

import (
	"errors"
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

type UpdatedEventFilterReq struct {
	Title  string `json:"title" query:"title"`
	Status string `json:"status" query:"status" validate:"oneof=pending approved rejected canceled"`
	Date   int64  `json:"date" query:"date"`
}

type UpdatedEventFilter struct {
	EventID *string    `json:"event_id"`
	Title   *string    `json:"title" query:"title"`
	Status  *string    `json:"status" query:"status"`
	Start   *time.Time `json:"start" query:"start"`
	End     *time.Time `json:"end" query:"end"`
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
	ID     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=approved rejected"`
}

type EventsUpdatesResp struct {
	ID                string            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID           string            `json:"event_id"`
	EventTitle        string            `json:"event_title"`
	UpdatedDetails    UpdatedDetails    `json:"updated_details"`
	UpdatedAddress    UpdatedAddress    `json:"updated_address"`
	UpdatedCategories []EventCategories `json:"updated_categories"`
	CreatedAt         int64             `json:"created_at" gorm:"default:now()"`
	UpdatedAt         int64             `json:"updated_at" gorm:"default:now()"`
	DeletedAt         *int64            `json:"deleted_at,omitempty"`
}

type EventUpdatesResp struct {
	ID                string            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID           string            `json:"event_id"`
	EventTitle        string            `json:"event_title"`
	UpdatedDetails    UpdatedDetails    `json:"updated_details"`
	UpdatedAddress    UpdatedAddress    `json:"updated_address"`
	UpdatedCategories []EventCategories `json:"updated_categories"`
	CreatedAt         int64             `json:"created_at" gorm:"default:now()"`
	UpdatedAt         int64             `json:"updated_at" gorm:"default:now()"`
	DeletedAt         *int64            `json:"deleted_at,omitempty"`
}

type UpdatedDetails struct {
	Banner      string `json:"banner,omitempty"`
	Status      string `json:"status"`
	Description string `json:"description"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
}

type UpdatedAddress struct {
	Address       string      `json:"address"`
	City          string      `json:"city"`
	Country       string      `json:"country"`
	DetailAddress string      `json:"detail_address"`
	Coordinates   Coordinates `json:"cooridinates"`
}

func ToEventUpdateResp(req *UpdatedEventRespReq) (*EventUpdatesResp, error) {
	if req.UpdatedEvent == nil {
		return nil, errors.New("updated event is nil")
	}

	resp := &EventUpdatesResp{
		ID:         req.UpdatedEvent.ID,
		EventID:    req.UpdatedEvent.EventID,
		EventTitle: req.UpdatedEvent.Name,
		UpdatedDetails: UpdatedDetails{
			Banner:      *req.UpdatedEvent.Banner,
			Status:      req.UpdatedEvent.Status.Status,
			Description: req.UpdatedEvent.Description,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
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
