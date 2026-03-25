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

// IsCategoriesChanged implements EventsRepo.
func (r *EventsRepository) IsCategoriesChanged(eventID string, ids []int64) bool {

	var old []int64

	r.db.
		Model(&model.EventCategories{}).
		Where("event_id = ?", eventID).
		Pluck("category_id", &old)

	if len(old) != len(ids) {
		return true
	}

	set := map[int64]bool{}

	for _, id := range old {
		set[id] = true
	}

	for _, id := range ids {
		if !set[id] {
			return true
		}
	}

	return false
}

// Create implements EventsRepo.
func (r *EventsRepository) Create(event *model.Events, categories []*model.Categories, tickets []*model.Tickets) (*model.Events, error) {
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
		return nil, err
	}

	// Create event categories
	if len(categories) > 0 {
		var eventCategories []*model.EventCategories
		for _, category := range categories {
			eventCategories = append(eventCategories, &model.EventCategories{
				EventID:    event.ID,
				CategoryID: category.ID,
				DeletedAt:  nil,
			})
		}
		if err := tx.Create(eventCategories).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Create tickets
	if len(tickets) > 0 {
		for _, ticket := range tickets {
			ticket.EventID = event.ID
			ticket.DeletedAt = nil
			if err := tx.Create(ticket).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return event, nil
}

// ReviewEvent implements EventsRepo.
func (r *EventsRepository) ReviewEvent(event *model.Events) error {
	return r.db.Updates(&event).Error
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

	event.DeletedAt = &now
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
		category.DeletedAt = &now
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
		ticket.DeletedAt = &now
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

// CancelEvent implements EventsRepo.
func (r *EventsRepository) CancelEvent(id string) error {
	return r.db.Model(&model.Events{}).
		Where("id = ?", id).
		Update("status", model.Cancel).Error
}

// FindAll implements EventsRepo.
func (r *EventsRepository) FindAll(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	query := r.db.Select(
		`events.*,
		ST_Y(coordinates::geometry) AS lat,
		ST_X(coordinates::geometry) AS lon`,
	).Scopes(filterEventList(filter))

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Preload("Profile.User").
		Preload("Categories.Category").
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

// FindActiveEvents implements EventsRepo.
func (r *EventsRepository) FindActiveEvents(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	if filter.Status != nil {
		return nil, errors.New("unauthorized action")
	}

	query := r.db.Scopes(filterEventList(filter))

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Preload("Profile.User").
		Preload("Categories.Category").
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

// FindByProfileID implements EventsRepo.
func (r *EventsRepository) FindByProfileID(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	query := r.db.Select(
		`events.*,
		ST_Y(coordinates::geometry) AS lat,
		ST_X(coordinates::geometry) AS lon`,
	).Scopes(filterEventList(filter))

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Preload("Profile.User").
		Preload("Categories.Category").
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

	if err := r.db.
		Select(`
			events.*,
			ST_Y(coordinates::geometry) AS lat,
			ST_X(coordinates::geometry) AS lon
		`).
		Where("id = ?", id).
		Preload("Profile.User").
		Preload("Categories.Category").
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
		Preload("Profile.User").
		Preload("Categories.Category").
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

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("Cannot request an update while a previous update is still pending review for this event")
		}

		return err
	}

	// Update event categories
	if err := tx.Where("event_id = ?", event.ID).Delete(&model.EventCategories{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(categories) > 0 {
		var eventCategories []*model.EventCategories
		for _, category := range categories {
			eventCategories = append(eventCategories, &model.EventCategories{
				EventID:    event.ID,
				CategoryID: category.ID,
				DeletedAt:  nil,
			})
		}
		if err := tx.Save(eventCategories).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update tickets
	if err := tx.Where("event_id = ?", event.ID).Delete(&model.Tickets{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(tickets) > 0 {
		for _, ticket := range tickets {
			ticket.EventID = event.ID
			ticket.DeletedAt = nil
			if err := tx.Save(ticket).Error; err != nil {
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

func filterEventList(filter *dto.EventFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.ProfileID != nil {
			db = db.Where("profile_id = ?", filter.ProfileID)
		}

		if filter.Status != nil {
			db = db.Where("status = ?", *filter.Status)
		} else {
			db = db.Where("status = ?", model.Active)
		}

		if filter.Title != nil {
			db = db.Where("LOWER(slug) LIKE LOWER(?)", "%"+*filter.Title+"%")
		}

		if filter.City != nil {
			db = db.Where("LOWER(city) LIKE LOWER(?)", "%"+*filter.City+"%")
		}

		if filter.Country != nil {
			db = db.Where("LOWER(country) LIKE LOWER(?)", "%"+*filter.Country+"%")
		}

		if filter.Category != nil {
			db = db.InnerJoins("Categories", db.Where("(category_id) IN ?", filter.Category))
		}

		if filter.Start != nil {
			db = db.Where("date >= ?", filter.Start)
		}

		if filter.End != nil {
			db = db.Where("date < ?", filter.End)
		}

		return db
	}
}

func toEventResponse(events []*model.Events) ([]*dto.EventsResp, error) {
	var eventsResp []*dto.EventsResp

	if len(events) == 0 {
		return nil, errors.New("no data found")
	}

	for _, event := range events {
		var categories []dto.EventCategories
		for _, cat := range event.Categories {
			categories = append(categories, dto.EventCategories{
				ID:   cat.Category.ID,
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
				Banner:      event.Banner,
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
				Coordinates: dto.Coordinates{
					Lat: event.Lat,
					Lon: event.Lon,
				},
			},
			Date:      helper.ConvertDatetoUnix(event.Date.Format(time.RFC3339)),
			CreatedAt: helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
			UpdatedAt: helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
			DeletedAt: helper.TimePtrToUnix(event.DeletedAt),
		})
	}

	return eventsResp, nil
}
