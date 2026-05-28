package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type EventsRepo interface {
	GetDB() *gorm.DB

	// Events
	Create(event *model.Events, categories []*model.Categories) (*model.Events, error)
	FindAll(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindActiveEvents(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindByProfileID(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindByID(id string) (*model.Events, error)
	FindBySlug(slug string, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindNearestEvents(lat, lon float64, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	UpdateStatus(id, status string) error
	IsCategoriesChanged(eventID string, ids []int64) bool
	Update(event *model.Events, categories []*model.Categories) error
	UpdateBannerEvent(id, banner string) error
	ReviewEvent(event *model.Events) error
	CancelEvent(id string) error
	Delete(id string) error
}
