package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
)

type UpdatedEventService struct {
	UpdatedEventRepo     repository.EventsUpdateRepo
	EventRepo            repository.EventsRepo
	OrganizerProfileRepo repository.OrganizerProfileRepo
	rdb                  *redis.Client
	EmailTaskPublisher   NewTaskEmail
}

func NewUpdatedEventService(
	updatedEventRepo repository.EventsUpdateRepo,
	eventRepo repository.EventsRepo,
	organizerProfileRepo repository.OrganizerProfileRepo,
	rdb *redis.Client,
	emailTaskPublisher NewTaskEmail,
) *UpdatedEventService {
	return &UpdatedEventService{
		UpdatedEventRepo:     updatedEventRepo,
		EventRepo:            eventRepo,
		OrganizerProfileRepo: organizerProfileRepo,
		rdb:                  rdb,
		EmailTaskPublisher:   emailTaskPublisher,
	}
}

var updatedEventCache []string = []string{
	"updated_events:all:*",
	"updated_event_list:all:*",
}

func (s *UpdatedEventService) InvalidateUpdatedEventCache() {
	ctx := context.Background()

	patterns := []string{
		"updated_events:all:*",
		"updated_event_list:all:*",
	}

	for _, pattern := range patterns {
		iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			s.rdb.Del(ctx, iter.Val())
		}
	}

	// Use SCAN for pattern keys to avoid blocking
	log.Println("[CACHE] updated events cache invalidated")
}

func (s *UpdatedEventService) ListAllUpdatedEvents(filter *dto.UpdatedEventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsUpdatesResp], error) {
	var updatedEvents *model.PaginationRow[*dto.EventsUpdatesResp]

	filterBytes, _ := json.Marshal(filter)
	hash := sha1.Sum(filterBytes)
	filterHash := hex.EncodeToString(hash[:])

	cacheKey := fmt.Sprintf("updated_events:all:%d:%d:%s:%s", pagination.Limit, pagination.Page, pagination.Sort, filterHash)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &updatedEvents)
	}

	if updatedEvents == nil {
		// if cache miss, get from db
		updatedEvents, err = s.UpdatedEventRepo.FindAll(filter, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(updatedEvents); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return updatedEvents, nil
}

func (s *UpdatedEventService) ListAllUpdatedEventsByEventID(filter *dto.UpdatedEventFilter, pagination model.Pagination, eventID string) (*model.PaginationRow[*dto.EventsUpdatesResp], error) {
	var updatedEvents *model.PaginationRow[*dto.EventsUpdatesResp]

	filter.EventID = &eventID

	filterBytes, _ := json.Marshal(filter)
	hash := sha1.Sum(filterBytes)
	filterHash := hex.EncodeToString(hash[:])

	cacheKey := fmt.Sprintf("updated_event_list:all:%d:%d:%s:%s", pagination.Limit, pagination.Page, pagination.Sort, filterHash)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &updatedEvents)
	}

	if updatedEvents == nil {
		// if cache miss, get from db
		updatedEvents, err = s.UpdatedEventRepo.FindAllByEventID(filter, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(updatedEvents); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return updatedEvents, nil
}

func (s *UpdatedEventService) GetUpdateEventByID(id, userID, role string) (*dto.EventUpdatesResp, error) {
	// Validate event
	event, err := s.UpdatedEventRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("updated event request not found")
	}

	var categories []dto.EventCategories
	for _, cat := range event.Categories {
		categories = append(categories, dto.EventCategories{
			ID:   cat.Category.ID,
			Name: cat.Category.Name,
		})
	}

	startTime := helper.ConvertDatetoUnix(event.StartTime.Format(time.RFC3339))
	endTime := helper.ConvertDatetoUnix(event.EndTime.Format(time.RFC3339))
	req := &dto.UpdatedEventRespReq{
		UpdatedEvent:    event,
		EventID:         event.EventID,
		EventCategories: categories,
		StartTime:       startTime,
		EndTime:         endTime,
		CreatedAt:       helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:       helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
		DeletedAt:       helper.TimePtrToUnix(event.DeletedAt),
	}

	resp, err := dto.ToEventUpdateResp(req)

	return resp, nil
}

func (s *UpdatedEventService) ReviewUpdated(req *dto.ReviewUpdatedEventReq) error {
	// Begin trasaction
	updatedX := s.UpdatedEventRepo.GetDB().Begin()
	eventX := s.EventRepo.GetDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			updatedX.Rollback()
			eventX.Rollback()
		}
	}()

	// Validate updated event
	updatedEvent, err := s.UpdatedEventRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("update event not found")
	}

	if updatedEvent.Status != string(model.Pending) {
		return errors.New(fmt.Sprintf("Updated already %s", updatedEvent.Status))
	}

	event, err := s.EventRepo.FindByID(updatedEvent.EventID)
	if err != nil {
		return errors.New("event not found")
	}

	updatedEvent.Status = req.Status

	// Review the update
	if err := updatedX.Updates(updatedEvent).Error; err != nil {
		log.Printf("[ERROR] review event failed with %v error", err)
		updatedX.Rollback()
		return errors.New(fmt.Sprintf("review updated event failed"))
	}

	// If the status approved
	// Update the event with the updated event data
	if req.Status == "approved" {
		update := &dto.UpdateEvent{
			EventTx:      eventX,
			UpdatedEvent: updatedEvent,
			Event:        event,
		}

		oldBanner, err := updateEventWithUpdated(update)
		if err != nil {
			updatedX.Rollback()
			eventX.Rollback()
			return err
		}

		if oldBanner != "" {
			if err := os.Remove(filepath.Join(eventBannerPath, oldBanner)); err != nil {
				updatedX.Rollback()
				eventX.Rollback()
				return err
			}
		}
	}

	// Commit transaction
	if err := updatedX.Commit().Error; err != nil {
		return err
	}

	if err := eventX.Commit().Error; err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, updatedEventCache)

	// Send email to organizer
	organizer, err := s.OrganizerProfileRepo.FindByID(updatedEvent.Event.ProfileID)
	if err != nil {
		log.Println("[ERROR] organizer data not found")
	}

	payload := &model.EventEmailPayload{
		To:        organizer.User.Email,
		EOName:    organizer.Name,
		EventName: updatedEvent.Name,
		Status:    updatedEvent.Status,
	}

	if err := s.EmailTaskPublisher.Enqueue(model.TypeEventUpdateNotification, payload); err != nil {
		log.Printf("[ERROR] failed sending email to %s with error %v", organizer.User.Email, err)
	}

	return nil
}

func (s *UpdatedEventService) CancelUpdate(id string) error {
	if err := s.UpdatedEventRepo.Cancel(id); err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, updatedEventCache)

	return nil
}

func updateEventWithUpdated(update *dto.UpdateEvent) (string, error) {
	var oldBanner string

	// Save the update event
	if update.UpdatedEvent.Banner != nil && update.Event.Banner != nil &&
		*update.UpdatedEvent.Banner != *update.Event.Banner {
		fileName := *update.UpdatedEvent.Banner

		updatedBannerSrc := filepath.Join(updatedEventBannerPath, fileName)
		dstPath := filepath.Join(eventBannerPath, fileName)

		bannerFile, err := helper.CopyFile(updatedBannerSrc, dstPath)
		if err != nil {
			return "", err
		}

		oldBanner = *update.Event.Banner
		update.Event.Banner = &bannerFile
	}

	update.Event.Name = update.UpdatedEvent.Name
	update.Event.Slug = update.UpdatedEvent.Slug
	update.Event.Description = update.UpdatedEvent.Description
	update.Event.Address = update.UpdatedEvent.Address
	update.Event.City = update.UpdatedEvent.City
	update.Event.Country = update.UpdatedEvent.Country
	update.Event.DetailAddress = update.UpdatedEvent.DetailAddress
	update.Event.Coordinates = update.UpdatedEvent.Coordinates
	update.Event.StartTime = update.UpdatedEvent.StartTime
	update.Event.EndTime = update.UpdatedEvent.EndTime
	update.Event.UpdatedAt = time.Now().UTC()

	if err := update.EventTx.Updates(update.Event).Error; err != nil {
		log.Printf("[ERROR] update event failed with error %v", err)
		return "", errors.New("failed to update the event")
	}

	return oldBanner, nil
}
