package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type EventsUpdateRepo interface {
	GetDB() *gorm.DB
	Create(event *model.UpdatedEvents, categories []*model.Categories) error
	FindAll(filter *dto.UpdatedEventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error)
	FindAllByEventID(filter *dto.UpdatedEventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error)
	FindByID(id string) (*model.UpdatedEvents, error)
	FindByEventID(eventID, status string) (*model.UpdatedEvents, error)
	ReviewEvent(id, status string) error
	Cancel(id string) error
	SoftDeleteEventUpdates(tx *gorm.DB, profileID string) error
}
