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

type UpdatedEventsRepository struct {
	db *gorm.DB
}

func NewUpdatedEventsRepository(db *gorm.DB) EventsUpdateRepo {
	return &UpdatedEventsRepository{db: db}
}

// Create implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) Create(event *model.UpdatedEvents, categories []*model.Categories, tickets []*model.TicketsUpdate) error {
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
		var eventCategories []*model.EventCategoriesUpdate
		for _, category := range categories {
			eventCategories = append(eventCategories, &model.EventCategoriesUpdate{
				EventUpdateID: event.ID,
				CategoryID:    category.ID,
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
			ticket.EventUpdateID = event.ID
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

// Delete implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) Cancel(id string) error {
	return r.db.
		Model(&model.UpdatedEvents{}).
		Where("id = ?", id).
		Updates(&model.UpdatedEvents{
			Status:    string(model.UpdatedCanceled),
			DeletedAt: helper.TimeToPointer(time.Now().UTC()),
		}).Error
}

// FindAll implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) FindAll(filter *dto.UpdatedEventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error) {
	var updatedEvents []*model.UpdatedEvents

	query := r.db.Scopes(filterUpdatedEventList(filter))

	if err := query.
		Scopes(Paginate(updatedEvents, &pagination, query)).
		Preload("Event").
		Preload("Categories.Category").
		Preload("Tickets").
		Find(&updatedEvents).Error; err != nil {
		return nil, err
	}

	updatedEventsResp, err := toUpdatedEventsResp(updatedEvents)
	if err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.EventsUpdatesResp]{
		Pagination: pagination,
		Rows:       updatedEventsResp,
	}, nil
}

// FindAllByEventID implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) FindAllByEventID(filter *dto.UpdatedEventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error) {
	var updatedEvents []*model.UpdatedEvents

	query := r.db.Scopes(filterUpdatedEventList(filter))

	if err := query.
		Scopes(Paginate(updatedEvents, &pagination, query)).
		Preload("Event").
		Preload("Categories.Category").
		Preload("Tickets").
		Find(&updatedEvents).Error; err != nil {
		return nil, err
	}

	updatedEventsResp, err := toUpdatedEventsResp(updatedEvents)
	if err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.EventsUpdatesResp]{
		Pagination: pagination,
		Rows:       updatedEventsResp,
	}, nil
}

// FindByID implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) FindByID(id string) (*model.UpdatedEvents, error) {
	var updatedEvent *model.UpdatedEvents

	if err := r.db.Where("id = ?", id).
		Preload("Event").
		Preload("Categories.Category").
		Preload("Tickets").
		First(&updatedEvent).Error; err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

// GetDB implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) GetDB() *gorm.DB {
	return r.db
}

// ReviewEvent implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) ReviewEvent(id string, status string) error {
	return r.db.Model(&model.UpdatedEvents{}).Where("id = ?", id).Update("status", status).Error
}

func filterUpdatedEventList(filter *dto.UpdatedEventFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.EventID != nil {
			db = db.Where("event_id = ?", filter.EventID)
		}

		if filter.Status != nil {
			db = db.Where("status = ?", *filter.Status)
		}

		if filter.Title != nil {
			db = db.Where("LOWER(slug) LIKE LOWER(?)", "%"+*filter.Title+"%")
		}

		if filter.Start != nil {
			db = db.Where("created_at >= ?", filter.Start)
		}

		if filter.End != nil {
			db = db.Where("created_at < ?", filter.End)
		}

		return db
	}
}

func toUpdatedEventsResp(updatedEvents []*model.UpdatedEvents) ([]*dto.EventsUpdatesResp, error) {
	var updatedEventResp []*dto.EventsUpdatesResp

	if len(updatedEvents) == 0 {
		return nil, errors.New("no data found")
	}

	for _, event := range updatedEvents {
		var categories []dto.EventCategories
		for _, cat := range event.Categories {
			categories = append(categories, dto.EventCategories{
				ID:   cat.Category.ID,
				Name: cat.Category.Name,
			})
		}

		updatedEventResp = append(updatedEventResp, &dto.EventsUpdatesResp{
			ID:         event.ID,
			EventID:    event.Event.ID,
			EventTitle: event.Name,
			UpdatedDetails: dto.UpdatedDetails{
				Banner:      *event.Banner,
				Status:      event.Status,
				Description: event.Description,
				Date:        helper.ConvertDatetoUnix(event.Date.Format(time.RFC3339)),
			},
			UpdatedAddress: dto.UpdatedAddress{
				Address:       event.Address,
				City:          event.City,
				Country:       event.Country,
				DetailAddress: event.DetailAddress,
				Coordinates:   event.Coordinates,
			},
			UpdatedCategories: categories,
			UpdatedTickets:    len(event.Tickets),
			CreatedAt:         helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:         helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
			DeletedAt:         helper.TimePtrToUnix(event.DeletedAt),
		})
	}

	return updatedEventResp, nil
}
