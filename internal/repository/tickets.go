package repository

import (
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketsRepo {
	return &TicketRepository{db: db}
}

// Create implements TicketsRepo.
func (r *TicketRepository) Create(ticket *model.Tickets) error {
	return r.db.Create(ticket).Error
}

// FindByEventID implements TicketsRepo.
func (r *TicketRepository) FindByEventID(id string) ([]*model.Tickets, error) {
	var tickets []*model.Tickets

	if err := r.db.Where("event_id = ?", id).Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

// Delete implements TicketsRepo.
func (r *TicketRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Tickets{}).Error
}

// Update implements TicketsRepo.
func (r *TicketRepository) Update(ticket *model.Tickets) error {
	return r.db.Updates(ticket).Error
}
