package service

import (
	"errors"
	"log"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"

	"github.com/redis/go-redis/v9"
)

type EventExpiryService struct {
	EventRepo        repository.EventsRepo
	UpdatedEventRepo repository.EventsUpdateRepo
	rdb              *redis.Client
}

func NewEventExpiryService(
	eventRepo repository.EventsRepo,
	updatedEventRepo repository.EventsUpdateRepo,
	rdb *redis.Client,
) *EventExpiryService {
	return &EventExpiryService{
		EventRepo:        eventRepo,
		UpdatedEventRepo: updatedEventRepo,
		rdb:              rdb,
	}
}

func (s *EventExpiryService) MarkEventAsDone(eventID string) error {
	event, err := s.EventRepo.FindByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Only active events should be marked done
	if event.Status != string(model.Active) {
		log.Printf("[EXPIRY] skipping event %s: status is %s, expected active", eventID, event.Status)
		return nil
	}

	if err := s.EventRepo.UpdateStatus(eventID, string(model.Done)); err != nil {
		return errors.New("failed to mark event as done")
	}

	utils.InvalidateCache(s.rdb, eventCache)

	log.Printf("[EXPIRY] event %s marked as done", eventID)
	return nil
}

func (s *EventExpiryService) MarkUpdatedEventAsDone(updatedEventID, eventID string) error {
	updatedEvent, err := s.UpdatedEventRepo.FindByID(updatedEventID)
	if err != nil {
		return errors.New("updated event not found")
	}

	// Only approved updated events need expiry tracking
	if updatedEvent.Status != string(model.Approved) {
		log.Printf("[EXPIRY] skipping updated event %s: status is %s", updatedEventID, updatedEvent.Status)
		return nil
	}

	// Mark the underlying event as done
	event, err := s.EventRepo.FindByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Status != string(model.Active) {
		log.Printf("[EXPIRY] skipping event %s from updated event expiry: status is %s", eventID, event.Status)
		return nil
	}

	if err := s.EventRepo.UpdateStatus(eventID, "done"); err != nil {
		return errors.New("failed to mark event as done via updated event expiry")
	}

	utils.InvalidateCache(s.rdb, eventCache)

	log.Printf("[EXPIRY] event %s marked as done (triggered by updated event %s)", eventID, updatedEventID)
	return nil
}
