package tasks

import (
	"context"
	"encoding/json"
	"log"
	"ngevent/internal/model"
	"ngevent/internal/service"

	"github.com/hibiken/asynq"
)

type EventExpiryHandler struct {
	EventExpiryService *service.EventExpiryService
	EventService       *service.EventService
}

func NewEventExpiryHandler(eventExpiryService *service.EventExpiryService, eventService *service.EventService) *EventExpiryHandler {
	return &EventExpiryHandler{EventExpiryService: eventExpiryService, EventService: eventService}
}

func (h *EventExpiryHandler) HandleEventExpired(ctx context.Context, t *asynq.Task) error {
	var p model.EventExpiredPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	if err := h.EventExpiryService.MarkEventAsDone(p.EventID); err != nil {
		log.Printf("[EXPIRY] error marking event %s as done: %v", p.EventID, err)
		return err
	}

	return nil
}

func (h *EventExpiryHandler) HandleUpdatedEventExpired(ctx context.Context, t *asynq.Task) error {
	var p model.UpdatedEventExpiredPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	if err := h.EventExpiryService.MarkUpdatedEventAsDone(p.UpdatedEventID, p.EventID); err != nil {
		log.Printf("[EXPIRY] error marking event %s as done via update %s: %v", p.EventID, p.UpdatedEventID, err)
		return err
	}

	return nil
}

func (h *EventExpiryHandler) HandleDraftRevert(ctx context.Context, t *asynq.Task) error {
	var p model.DraftRevertPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	if err := h.EventExpiryService.RevertToDraft(p.EventID); err != nil {
		log.Printf("[DRAFT-REVERT] error reverting event %s: %v", p.EventID, err)
		return err
	}
	return nil
}


func (h *EventExpiryHandler) HandleUpdatedDraftRevert(ctx context.Context, t *asynq.Task) error {
	var p model.UpdatedDraftRevertPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	if err := h.EventExpiryService.UpdateRevertToDraft(p.UpdatedEventID); err != nil {
		log.Printf("[DRAFT-REVERT] error reverting updated event %s: %v", p.UpdatedEventID, err)
		return err
	}
	return nil
}