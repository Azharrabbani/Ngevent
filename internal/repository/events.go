package repository

import (
	"errors"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"time"

	"gorm.io/gorm"
)

type EventsRepository struct {
	db *gorm.DB
}

// UpdateBannerEvent implements EventsRepo.
func (r *EventsRepository) UpdateBannerEvent(id, banner string) error {
	return r.db.Model(&model.Events{}).
		Where("id = ?", id).
		Updates(&model.Events{
			Banner:    &banner,
			UpdatedAt: time.Now().UTC()}).Error
}

func NewEventsRepository(db *gorm.DB) EventsRepo {
	return &EventsRepository{db: db}
}

// GetDB implements EventsRepo.
func (r *EventsRepository) GetDB() *gorm.DB {
	return r.db
}

// Create implements EventsRepo.
func (r *EventsRepository) Create(event *model.Events, categories []*model.Categories, tickets []*model.Tickets) error {
	// Make transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] Transaction rolled back: %v", r)
		}
	}()

	// Create event
	if err := tx.Create(event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create event categories
	if len(categories) > 0 {
		var eventCategories []*model.EventCategories
		for _, category := range categories {
			eventCategories = append(eventCategories, &model.EventCategories{
				EventID:    event.ID,
				CategoryID: category.ID,
			})
		}
		if err := tx.Create(eventCategories).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create tickets
	if len(tickets) > 0 {
		for _, ticket := range tickets {
			ticket.EventID = event.ID
			if err := tx.Create(ticket).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// ReviewEvent implements EventsRepo.
func (r *EventsRepository) ReviewEvent(event *model.Events) error {
	return r.db.Save(&event).Error
}

// Delete implements EventsRepo.
func (r *EventsRepository) Delete(id string) error {
	// Make transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] Transaction rolled back: %v", r)
		}
	}()

	now := time.Now().UTC()

	// Update event
	var event *model.Events
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return err
	}

	event.DeletedAt = now
	if err := tx.Updates(event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update event categories
	var eventCategories []*model.EventCategories
	if err := r.db.Where("event_id = ?", event.ID).Find(&eventCategories).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, category := range eventCategories {
		category.DeletedAt = now
		if err := tx.Updates(category).Error; err != nil {
			tx.Rollback()
			return err
		}

	}

	// Update tickets
	var tickets []*model.Tickets
	if err := r.db.Where("event_id = ?", event.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, ticket := range tickets {
		ticket.DeletedAt = now
		if err := tx.Updates(ticket).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil

}

// FindAll implements EventsRepo.
func (r *EventsRepository) FindAll(pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	if err := r.db.
		Scopes(Paginate(events, &pagination, r.db)).
		Preload("Profile").
		Preload("Categories").
		Preload("Tickets").
		Find(&events).Error; err != nil {
		return nil, err
	}

	eventsResp, err := toEventResponse(events)
	if err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.EventsResp]{
		Pagination: pagination,
		Rows:       eventsResp,
	}, nil

}

// FindByID implements EventsRepo.
func (r *EventsRepository) FindByID(id string) (*model.Events, error) {
	var event *model.Events

	if err := r.db.Where("id = ?", id).
		Preload("Profile").
		Preload("Categories").
		Preload("Tickets").
		First(&event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

// FindBySlug implements EventsRepo.
func (r *EventsRepository) FindBySlug(slug string, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	query := r.db.Where("LOWER(slug) LIKE LOWER(?)", "%"+slug+"%")

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Preload("Profile").
		Preload("Categories").
		Preload("Tickets").
		Find(&events).Error; err != nil {
		return nil, errors.New("event not found")
	}

	eventsResp, err := toEventResponse(events)
	if err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.EventsResp]{
		Pagination: pagination,
		Rows:       eventsResp,
	}, nil
}

// Update implements EventsRepo.
func (r *EventsRepository) Update(event *model.Events, categories []*model.Categories, tickets []*model.Tickets) error {
	// Make transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] Transaction rolled back: %v", r)
		}
	}()

	// Update event
	if err := tx.Updates(event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update event categories
	if err := r.db.Where("event_id = ?", event.ID).Delete(&model.EventCategories{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(categories) > 0 {
		var eventCategories []*model.EventCategories
		for _, category := range categories {
			eventCategories = append(eventCategories, &model.EventCategories{
				EventID:    event.ID,
				CategoryID: category.ID,
			})
		}
		if err := tx.Updates(eventCategories).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update tickets
	if err := r.db.Where("event_id = ?", event.ID).Delete(&model.Tickets{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update tickets
	if len(tickets) > 0 {
		for _, ticket := range tickets {
			if err := tx.Updates(ticket).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func toEventResponse(events []*model.Events) ([]*dto.EventsResp, error) {
	var eventsResp []*dto.EventsResp

	if len(eventsResp) < 0 {
		return nil, errors.New("no data found")
	}

	for _, event := range events {
		var categories []dto.EventCategories
		for _, cat := range event.Categories {
			categories = append(categories, dto.EventCategories{
				ID:   cat.ID,
				Name: cat.Category.Name,
			})
		}

		var tickets []dto.Tickets
		for _, ticket := range event.Tickets {
			tickets = append(tickets, dto.Tickets{
				ID:         ticket.ID,
				Name:       ticket.Name,
				Price:      ticket.Price,
				Quantity:   ticket.Quantity,
				TicketType: ticket.TicketType,
			})
		}

		eventsResp = append(eventsResp, &dto.EventsResp{
			ID: event.ID,
			EOProfile: dto.EOProfiles{
				ID:           event.Profile.ID,
				IsVerified:   event.Profile.User.IsVerified,
				Email:        event.Profile.User.Email,
				Name:         event.Profile.Name,
				PhotoProfile: event.Profile.PhotoProfile,
				PhoneNumber:  event.Profile.PhoneNumber,
			},
			Event: dto.EventDetail{
				Banner:      *event.Banner,
				Name:        event.Name,
				Categories:  categories,
				Tickets:     tickets,
				Slug:        event.Slug,
				Status:      event.Status,
				Description: event.Description,
			},
			EventAddress: dto.EventAddress{
				Address:       event.Address,
				City:          event.City,
				Country:       event.Country,
				DetailAddress: event.DetailAddress,
				Coordinates:   event.Coordinates,
			},
			Date:      helper.ConvertDatetoUnix(event.Date.Format(time.RFC3339)),
			CreatedAt: helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
			UpdatedAt: helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
			DeletedAt: helper.ConvertDatetoUnix(event.DeletedAt.Format(time.RFC3339)),
		})
	}

	return eventsResp, nil
}
