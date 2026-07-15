// internal/service/eventExpiryService.go
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
	EventRepo          repository.EventsRepo
	UpdatedEventRepo   repository.EventsUpdateRepo
	EmailTaskPublisher NewTaskEmail
	rdb                *redis.Client
}

func NewEventExpiryService(
	eventRepo repository.EventsRepo,
	updatedEventRepo repository.EventsUpdateRepo,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *EventExpiryService {
	return &EventExpiryService{
		EventRepo:          eventRepo,
		UpdatedEventRepo:   updatedEventRepo,
		EmailTaskPublisher: emailTaskPublisher,
		rdb:                rdb,
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

func (s *EventExpiryService) RevertToDraft(eventID string) error {
	event, err := s.EventRepo.FindByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Status != string(model.Pending) {
		return nil
	}

	reason := "not reviewed within the allotted time"
	event.Status = string(model.Draft)
	event.SubmittedAt = nil
	event.RejectedReason = &reason

	var categories []*model.Categories
	for _, ec := range event.Categories {
		categories = append(categories, &model.Categories{
			ID: ec.CategoryID,
		})
	}

	if err := s.EventRepo.Update(event, categories); err != nil {
		return err
	}

	utils.InvalidateCache(s.rdb, eventCache)

	// Notify the organizer that their event was reverted to draft
	if event.Profile.User.Email != "" {
		payload := &model.EventEmailPayload{
			To:        event.Profile.User.Email,
			EOName:    event.Profile.Name,
			EventName: event.Name,
		}
		if err := s.EmailTaskPublisher.Enqueue(model.TypeEventEORevertNotification, payload); err != nil {
			log.Printf("[DRAFT-REVERT] failed to send revert notification email for event %s: %v", event.ID, err)
		}
	}

	return nil
}

func (s *EventExpiryService) UpdateRevertToDraft(updateEventID string) error {
	updatedEvent, err := s.UpdatedEventRepo.FindByID(updateEventID)
	if err != nil {
		return errors.New("event not found")
	}

	if updatedEvent.Status != string(model.Pending) {
		return nil
	}

	tx := s.EventRepo.GetDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] UpdateRevertToDraft transaction rolled back: %v", r)
		}
	}()

	// Find the main event data, then reset the request_updates flag
	event, err := s.EventRepo.FindByID(updatedEvent.EventID)
	if err != nil {
		tx.Rollback()
		return errors.New("the main event data not found")
	}

	event.RequestUpdates = false

	if err := tx.Model(&event).Updates(map[string]interface{}{"request_updates": false}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update event")
	}

	if err := s.UpdatedEventRepo.Cancel(updateEventID); err != nil {
		tx.Rollback()
		return errors.New("failed to cancel updated event")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update event")
	}

	utils.InvalidateCache(s.rdb, updatedEventCache)

	// Notify the organizer that their update request was discarded
	if event.Profile.User.Email != "" {
		payload := &model.EventEmailPayload{
			To:        event.Profile.User.Email,
			EOName:    event.Profile.Name,
			EventName: event.Name,
		}
		if err := s.EmailTaskPublisher.Enqueue(model.TypeEventEOUpdateRevertNotification, payload); err != nil {
			log.Printf("[DRAFT-REVERT] failed to send update-revert notification email for updated event %s: %v", updatedEvent.ID, err)
		}
	}

	return nil
}
