package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

type EventsRepo interface {
	GetDB() *gorm.DB
	Create(event *model.Events, categories []*model.Categories) (*model.Events, error)
	FindAll(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindActiveEvents(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindByProfileID(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindByProfileIDPublic(filter *dto.EventFilterPublic, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindByID(id string) (*model.Events, error)
	FindBySlug(slug string) (*model.Events, error)
	FindNearestEvents(lat, lon float64, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	UpdateStatus(id, status string) error
	FindStaleUnreviewedEvents(reviewCutoff time.Time, now time.Time) ([]*model.Events, error)
	IsCategoriesChanged(eventID string, ids []int64) bool
	Update(event *model.Events, categories []*model.Categories) error
	CreateStagedUpdate(event *model.Events, updatedEvent *model.UpdatedEvents, updatedCategories []*model.Categories) error
	UpdateBannerEvent(id, banner string) error
	ReviewEvent(event *model.Events) error
	CancelEvent(tx *gorm.DB, id string) error
	HasBlockingEvents(profileID string) (bool, error)
	Delete(id string) error
	SoftDeleteEvents(tx *gorm.DB, profileID string) error
}
