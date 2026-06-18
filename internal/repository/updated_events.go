package repository

import (
	"errors"
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UpdatedEventsRepository struct {
	db *gorm.DB
}

func NewUpdatedEventsRepository(db *gorm.DB) EventsUpdateRepo {
	return &UpdatedEventsRepository{db: db}
}

// Delete implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) Cancel(id string) error {
	return r.db.
		Model(&model.UpdatedEvents{}).
		Where("id = ?", id).
		Updates(&model.UpdatedEvents{
			Status:    string(model.Cancelled),
			DeletedAt: helper.TimeToPointer(time.Now().UTC()),
		}).Error
}

func (r *UpdatedEventsRepository) SoftDeleteEventUpdates(tx *gorm.DB, profileID string) error {
	now := time.Now().UTC()

	// Collect event_update IDs
	var eventUpdateIDs []string
	if err := tx.Model(&model.UpdatedEvents{}).
		Joins("JOIN events ON events.id = event_updates.event_id").
		Where("events.profile_id = ? AND event_updates.deleted_at IS NULL", profileID).
		Pluck("event_updates.id", &eventUpdateIDs).Error; err != nil {
		return err
	}

	if len(eventUpdateIDs) == 0 {
		return nil
	}

	if err := tx.Model(&model.EventCategoriesUpdate{}).
		Where("event_update_id IN ? AND deleted_at IS NULL", eventUpdateIDs).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	if err := tx.Model(&model.UpdatedEvents{}).
		Where("id IN ? AND deleted_at IS NULL", eventUpdateIDs).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements EventsUpdateRepo.
func (r *UpdatedEventsRepository) FindAll(filter *dto.UpdatedEventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error) {
	var updatedEvents []*model.UpdatedEvents

	query := r.db.Select(
		`event_updates.*,
		ST_Y(event_updates.coordinates::geometry) AS lat,
		ST_X(event_updates.coordinates::geometry) AS lon`,
	).
		Scopes(filterUpdatedEventList(filter))

	if err := query.
		Scopes(Paginate(updatedEvents, &pagination, query)).
		Preload("Event.Profile.User").
		Preload("Categories.Category").
		Preload("Reviewer").
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
		Preload("Reviewer").
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

	if err := r.db.
		Select(`
			event_updates.*,
			ST_Y(coordinates::geometry) AS lat,
			ST_X(coordinates::geometry) AS lon
		`).
		Where("id = ?", id).
		Preload("Event").
		Preload("Categories.Category").
		Preload("Reviewer").
		First(&updatedEvent).Error; err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (r *UpdatedEventsRepository) FindByEventID(eventID, status string) (*model.UpdatedEvents, error) {
	var updatedEvent *model.UpdatedEvents

	if err := r.db.
		Select(`
			event_updates.*,
			ST_Y(coordinates::geometry) AS lat,
			ST_X(coordinates::geometry) AS lon
		`).
		Where("event_id = ? AND status = ?", eventID, status).
		Preload("Event").
		Preload("Categories.Category").
		Preload("Reviewer").
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
			db = db.Where("event_updates.event_id = ?", *filter.EventID)
		}

		if filter.Status != nil {
			db = db.Where("event_updates.status = ?", *filter.Status)
		}

		if filter.Title != nil {
			db = db.Where("LOWER(event_updates.slug) LIKE LOWER(?)", "%"+*filter.Title+"%")
		}

		if filter.Search != nil {
			search := "%" + strings.ToLower(*filter.Search) + "%"
			db = db.Joins(`
                JOIN events ON events.id = event_updates.event_id
                JOIN organizer_profiles ON organizer_profiles.id = events.profile_id
            `).Where(`
                LOWER(organizer_profiles.name) LIKE ?
                OR LOWER(organizer_profiles.email) LIKE ?
                OR LOWER(event_updates.name) LIKE ?
            `, search, search, search)
		}

		if filter.Sort != nil {
			db = db.Order(fmt.Sprintf("event_updates.submitted_at %s", *filter.Sort))
		} else {
			db = db.Order("event_updates.submitted_at DESC")
		}

		if filter.Date != nil {
			switch helper.StringValue(filter.Date) {
			case "week":
				db = db.Where("event_updates.submitted_at >= ?", time.Now().AddDate(0, 0, -7))
			case "month":
				db = db.Where("event_updates.submitted_at >= ?", time.Now().AddDate(0, -1, 0))
			case "year":
				db = db.Where("event_updates.submitted_at >= ?", time.Now().AddDate(-1, 0, 0))
			}
		}

		if filter.RangeStart != nil &&
			filter.RangeEnd != nil {

			db = db.Where(
				`
        			events.start_date < ?
        			AND
        			events.end_date >= ?
        		`,
				filter.RangeEnd,
				filter.RangeStart,
			)
		}

		db = db.Group("event_updates.id, event_updates.submitted_at, event_updates.name, event_updates.start_time")

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

		var reviewer *dto.Reviewer
		if event.Reviewer != nil {
			reviewer = &dto.Reviewer{
				ID:    event.Reviewer.ID,
				Email: event.Reviewer.Email,
			}
		}

		profile := event.Event.Profile

		updatedEventResp = append(updatedEventResp, &dto.EventsUpdatesResp{
			ID:         event.ID,
			EventID:    event.Event.ID,
			EventTitle: event.Name,
			EOProfile: dto.EOProfiles{
				ID:           profile.ID,
				IsVerified:   profile.User.IsVerified,
				Status:       profile.Status.Status,
				Email:        profile.User.Email,
				Name:         profile.Name,
				PhotoProfile: profile.PhotoProfile,
				PhoneNumber:  profile.PhoneNumber,
			},
			UpdatedDetails: dto.UpdatedDetails{
				Banner: helper.StrPointerIfNotEmpty(
					func() string {
						if event.Banner == nil {
							return ""
						}
						return fmt.Sprintf("http://localhost:8080/api/v1/updated-event/banner/%s", *event.Banner)
					}(),
				),
				Status:         event.Status,
				Description:    event.Description,
				StartDate:      helper.ConvertDatetoUnix(event.StartDate.Format(time.RFC3339)),
				EndDate:        helper.ConvertDatetoUnix(event.EndDate.Format(time.RFC3339)),
				StartTime:      helper.ConvertDatetoUnix(event.StartTime.Format(time.RFC3339)),
				EndTime:        helper.ConvertDatetoUnix(event.EndTime.Format(time.RFC3339)),
				RejectedReason: event.RejectedReason,
				ReviewedBy:     reviewer,
				ReviewedAt:     helper.TimePtrToUnix(event.ReviewedAt),
			},
			UpdatedAddress: dto.UpdatedAddress{
				Address:       event.Address,
				City:          event.City,
				Country:       event.Country,
				DetailAddress: event.DetailAddress,
				Coordinates: dto.Coordinates{
					Lat: event.Lat,
					Lon: event.Lon,
				},
			},
			UpdatedCategories: categories,
			CreatedAt:         helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:         helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
			SubmittedAt:       helper.ConvertDatetoUnix(event.SubmittedAt.Format(time.RFC3339)),
			DeletedAt:         helper.TimePtrToUnix(event.DeletedAt),
		})
	}

	return updatedEventResp, nil
}
