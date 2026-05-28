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
}

func NewEventExpiryHandler(eventExpiryService *service.EventExpiryService) *EventExpiryHandler {
	return &EventExpiryHandler{EventExpiryService: eventExpiryService}
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
