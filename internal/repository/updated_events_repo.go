package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type EventsUpdateRepo interface {
	GetDB() *gorm.DB
	Create(event *model.UpdatedEvents, categories []*model.Categories, tickets []*model.TicketsUpdate) error
	FindAll(pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error)
	FindByID(id string) (*model.UpdatedEvents, error)
	ReviewEvent(id, status string) error
	Delete(id string) error
}
