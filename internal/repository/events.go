package repository

import (
	"errors"
	"fmt"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EventsRepository struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) EventsRepo {
	return &EventsRepository{db: db}
}

// GetDB implements EventsRepo.
func (r *EventsRepository) GetDB() *gorm.DB {
	return r.db
}

// UpdateStatus implements [EventsRepo].
func (r *EventsRepository) UpdateStatus(id, status string) error {
	return r.db.Model(&model.Events{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateBannerEvent implements EventsRepo.
func (r *EventsRepository) UpdateBannerEvent(id, banner string) error {
	return r.db.Model(&model.Events{}).
		Where("id = ?", id).
		Updates(&model.Events{
			Banner:    &banner,
			UpdatedAt: time.Now().UTC()}).Error
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
func (r *EventsRepository) Create(event *model.Events, categories []*model.Categories) (*model.Events, error) {
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

	// Draft events cannot have event_updates because
	// staging updates are only created for active events.
	if err := tx.Model(&model.EventCategories{}).
		Where("event_id = ?", id).
		Update("deleted_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Events{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *EventsRepository) HasBlockingEvents(profileID string) (bool, error) {
	var count int64

	err := r.db.Model(&model.Events{}).
		Where(
			"profile_id = ? AND status IN ? AND deleted_at IS NULL",
			profileID,
			[]string{string(model.Pending), string(model.Active)},
		).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *EventsRepository) SoftDeleteEvents(tx *gorm.DB, profileID string) error {
	now := time.Now().UTC()

	// Collect event IDs belonging to this profile
	var eventIDs []string
	if err := tx.Model(&model.Events{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Pluck("id", &eventIDs).Error; err != nil {
		return err
	}

	if len(eventIDs) == 0 {
		return nil
	}

	// Soft delete event_categories
	if err := tx.Model(&model.EventCategories{}).
		Where("event_id IN ? AND deleted_at IS NULL", eventIDs).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	// Soft delete events
	if err := tx.Model(&model.Events{}).
		Where("id IN ? AND deleted_at IS NULL", eventIDs).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	return nil
}

// CancelEvent implements EventsRepo.
func (r *EventsRepository) CancelEvent(id string) error {
	return r.db.Model(&model.Events{}).
		Where("id = ?", id).
		Update("status", model.Cancelled).Error
}

// FindAll implements EventsRepo.
func (r *EventsRepository) FindAll(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	query := r.db.Select(
		`events.*,
		ST_Y(events.coordinates::geometry) AS lat,
		ST_X(events.coordinates::geometry) AS lon`,
	).Scopes(filterEventList(filter))

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Order("created_at DESC").
		Preload("Profile.User").
		Preload("Categories.Category").
		Preload("Reviewer").
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

	query := r.db.Select(
		`events.*,
		ST_Y(events.coordinates::geometry) AS lat,
		ST_X(events.coordinates::geometry) AS lon`,
	).Scopes(filterEventList(filter))

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Order("created_at DESC").
		Preload("Profile.User").
		Preload("Categories.Category").
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

// FindByCity implements EventsRepo.
func (r *EventsRepository) FindNearestEvents(lat, lon float64, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	query := r.db.Select(`
			events.*,
			ST_Y(coordinates::geometry) AS lat,
			ST_X(coordinates::geometry) AS lon,
			ST_Distance(
				coordinates,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography
			) AS distance
		`, lon, lat).
		Where(`
			ST_DWithin(
				coordinates,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
				8000
			) AND status = ?
		`, lon, lat, model.Active)

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Preload("Profile.User").
		Preload("Categories.Category").
		Order("distance ASC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	resp, err := toEventResponse(events)
	if err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.EventsResp]{
		Pagination: pagination,
		Rows:       resp,
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
		Preload("Reviewer").
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

// FindByProfileIDPublic implements EventsRepo.
func (r *EventsRepository) FindByProfileIDPublic(filter *dto.EventFilterPublic, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events []*model.Events

	query := r.db.Select(
		`events.*,
		ST_Y(coordinates::geometry) AS lat,
		ST_X(coordinates::geometry) AS lon`,
	).Scopes(filterEventListPublic(filter))

	if err := query.
		Scopes(Paginate(events, &pagination, query)).
		Preload("Profile.User").
		Preload("Categories.Category").
		Preload("Reviewer").
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
		Preload("Reviewer").
		First(&event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

// FindBySlug implements EventsRepo.
func (r *EventsRepository) FindBySlug(slug string) (*model.Events, error) {
	var event *model.Events

	if err := r.db.
		Select(`
			events.*,
			ST_Y(coordinates::geometry) AS lat,
			ST_X(coordinates::geometry) AS lon
		`).
		Where("slug = ?", slug).
		Preload("Profile.User").
		Preload("Categories.Category").
		First(&event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

// Update implements EventsRepo.
func (r *EventsRepository) Update(event *model.Events, categories []*model.Categories) error {
	// Make transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] Transaction rolled back: %v", r)
		}
	}()

	// Update event
	updateMap := map[string]interface{}{
		"profile_id":     event.ProfileID,
		"banner":         event.Banner,
		"name":           event.Name,
		"slug":           event.Slug,
		"description":    event.Description,
		"address":        event.Address,
		"city":           event.City,
		"country":        event.Country,
		"detail_address": event.DetailAddress,
		"coordinates":    event.Coordinates,
		"start_time":     event.StartTime,
		"end_time":       event.EndTime,
		"status":         event.Status,
	}

	if err := tx.Model(&model.Events{ID: event.ID}).Updates(updateMap).Error; err != nil {
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *EventsRepository) CreateStagedUpdate(event *model.Events, updatedEvent *model.UpdatedEvents, updatedCategories []*model.Categories) error {
	// Make transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] Transaction rolled back: %v", r)
		}
	}()

	// Update event
	if err := tx.Model(&model.Events{ID: event.ID}).Update("request_updates", true).Error; err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("Cannot request an update while a previous update is still in review")
		}
		return err
	}

	// Create stage update
	if err := tx.Create(updatedEvent).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create stage  udpate categories
	if len(updatedCategories) > 0 {
		var eventCategories []*model.EventCategoriesUpdate
		for _, category := range updatedCategories {
			eventCategories = append(eventCategories, &model.EventCategoriesUpdate{
				EventUpdateID: updatedEvent.ID,
				CategoryID:    category.ID,
			})
		}
		if err := tx.Create(eventCategories).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func filterEventList(filter *dto.EventFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Role != nil && *filter.Role == "admin" {
			db = db.Unscoped()
		} else {
			db = db.Where("events.deleted_at IS NULL")
		}

		if filter.ProfileID != nil {
			db = db.Where("events.profile_id = ?", filter.ProfileID)
		}

		if filter.Search != nil {
			search := "%" + strings.ToLower(*filter.Search) + "%"

			db = db.Joins("JOIN organizer_profiles ON organizer_profiles.id = events.profile_id").
				Where(`
			(
				LOWER(organizer_profiles.name) LIKE ?
				OR LOWER(organizer_profiles.email) LIKE ?
				OR LOWER(events.name) LIKE ?
			)
		`,
					search,
					search,
					search,
				)
		}

		if filter.Sort != nil {
			db = db.Order(fmt.Sprintf("events.created_at %s", *filter.Sort))
		} else {
			db = db.Order("events.created_at DESC")
		}

		if filter.Date != nil {
			switch helper.StringValue(filter.Date) {
			case "week":
				db = db.Where(
					"events.created_at >= ?",
					time.Now().AddDate(0, 0, -7),
				)
			case "month":
				db = db.Where(
					"events.created_at >= ?",
					time.Now().AddDate(0, -1, 0),
				)
			case "year":
				db = db.Where(
					"events.created_at >= ?",
					time.Now().AddDate(-1, 0, 0),
				)
			}
		}

		if filter.Status != nil {
			db = db.Where("events.status = ?", *filter.Status)
		} else {
			status := []string{string(model.Active), string(model.Pending), string(model.Done), string(model.Rejected)}
			db = db.Where("status IN ?", status)
		}

		if filter.Title != nil {
			db = db.Where("LOWER(events.slug) LIKE LOWER(?)", "%"+*filter.Title+"%")
		}

		if filter.Location != nil {
			location := "%" + strings.ToLower(*filter.Location) + "%"
			db = db.Where("LOWER(city) LIKE LOWER(?)", location).Or("LOWER(country) LIKE LOWER(?)", location)
		}

		if len(filter.Category) > 0 {
			db = db.Joins("JOIN event_categories ec ON ec.event_id = events.id AND ec.deleted_at IS NULL").
				Where("ec.category_id IN ?", filter.Category).
				Group("events.id")
		}

		if filter.Start != nil {
			db = db.Where("events.start_time >= ?", filter.Start)
		}

		if filter.End != nil {
			db = db.Where("events.start_time < ?", filter.End)
		}

		db = db.Group("events.id, events.created_at, events.name, events.start_time")

		return db
	}
}

func filterEventListPublic(filter *dto.EventFilterPublic) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		db = db.Where("deleted_at IS NULL AND profile_id = ?", filter.ProfileID)

		if filter.Status != nil {
			db = db.Where("status = ?", *filter.Status)
		}

		if filter.Title != nil {
			db = db.Where("LOWER(slug) LIKE LOWER(?)", "%"+*filter.Title+"%")
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

		var reviewer *dto.Reviewer
		if event.Reviewer != nil {
			reviewer = &dto.Reviewer{
				ID:    event.Reviewer.ID,
				Email: event.Reviewer.Email,
			}
		}

		eventsResp = append(eventsResp, &dto.EventsResp{
			ID: event.ID,
			EOProfile: dto.EOProfiles{
				ID:         event.Profile.ID,
				IsVerified: event.Profile.User.IsVerified,
				Status:     event.Profile.Status.Status,
				Email:      event.Profile.User.Email,
				Name:       event.Profile.Name,
				PhotoProfile: helper.StrPointerIfNotEmpty(
					func() string {
						if event.Profile.PhotoProfile == nil {
							return ""
						}
						return fmt.Sprintf("http://localhost:8080/api/v1/organizer/photo/%s", *event.Profile.PhotoProfile)
					}(),
				),
				PhoneNumber: event.Profile.PhoneNumber,
			},
			Event: dto.EventDetail{
				Banner: helper.StrPointerIfNotEmpty(
					func() string {
						if event.Banner == nil {
							return ""
						}
						return fmt.Sprintf("http://localhost:8080/api/v1/event/banner/%s", *event.Banner)
					}(),
				),
				Name:           event.Name,
				Categories:     categories,
				Slug:           event.Slug,
				Status:         event.Status,
				RequestUpdates: event.RequestUpdates,
				Description:    event.Description,
				RejectedReason: event.RejectedReason,
				ReviewedBy:     reviewer,
				ReviewedAt:     helper.TimePtrToUnix(event.ReviewedAt),
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
			StartTime: helper.ConvertDatetoUnix(event.StartTime.Format(time.RFC3339)),
			EndTime:   helper.ConvertDatetoUnix(event.EndTime.Format(time.RFC3339)),
			CreatedAt: helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
			UpdatedAt: helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
			DeletedAt: helper.TimePtrToUnix(event.DeletedAt),
		})
	}

	return eventsResp, nil
}
