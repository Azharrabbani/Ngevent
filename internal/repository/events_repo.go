package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type EventsRepo interface {
	GetDB() *gorm.DB
	Create(event *model.Events, categories []*model.Categories, tickets []*model.Tickets) error
	FindAll(pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	FindByID(id string) (*model.Events, error)
	FindBySlug(slug string, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error)
	Update(event *model.Events, categories []*model.Categories, tickets []*model.Tickets) error
	Delete(id string) error
}
