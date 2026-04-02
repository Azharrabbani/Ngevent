package repository

import "ngevent/internal/model"

type TicketsRepo interface {
	Create(ticket *model.Tickets) error
	FindByEventID(id string) ([]*model.Tickets, error)
	Update(ticket *model.Tickets) error
	Delete(id string) error
}
