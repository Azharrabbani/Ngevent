package tasks

import (
	"encoding/json"
	"fmt"
	"ngevent/internal/model"
	"time"

	"github.com/hibiken/asynq"
)

type EventExpiryPublisher struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewEventExpiryPublisher(client *asynq.Client, inspector *asynq.Inspector) *EventExpiryPublisher {
	return &EventExpiryPublisher{
		Client:    client,
		Inspector: inspector,
	}
}

// EnqueueEventExpiry schedules a task to mark the event as done at endTime.
func (p *EventExpiryPublisher) EnqueueEventExpiry(payload *model.EventExpiredPayload, endTime time.Time) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	delay := time.Until(endTime)
	if delay <= 0 {
		delay = 0
	}

	taskID := fmt.Sprintf("event_expiry:%s", payload.EventID)

	_ = p.cancelByID(taskID)

	task := asynq.NewTask(model.TypeEventExpired, b)
	_, err = p.Client.Enqueue(task,
		asynq.ProcessIn(delay),
		asynq.TaskID(taskID),
	)
	return err
}

// EnqueueUpdatedEventExpiry schedules a task to mark the event done when an
// approved updated event's end time is reached.
func (p *EventExpiryPublisher) EnqueueUpdatedEventExpiry(payload *model.UpdatedEventExpiredPayload, endTime time.Time) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	delay := time.Until(endTime)
	if delay <= 0 {
		delay = 0
	}

	taskID := fmt.Sprintf("updated_event_expiry:%s", payload.UpdatedEventID)

	_ = p.cancelByID(taskID)

	task := asynq.NewTask(model.TypeUpdatedEventExpired, b)
	_, err = p.Client.Enqueue(task,
		asynq.ProcessIn(delay),
		asynq.TaskID(taskID),
	)
	return err
}

// CancelEventExpiry removes a pending expiry task
func (p *EventExpiryPublisher) CancelEventExpiry(eventID string) error {
	return p.cancelByID(fmt.Sprintf("event_expiry:%s", eventID))
}

func (p *EventExpiryPublisher) EnqueueDraftRevert(payload *model.DraftRevertPayload, revertAt time.Time) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	delay := time.Until(revertAt)
	if delay <= 0 {
		delay = 0
	}

	taskID := fmt.Sprintf("draft_revert:%s", payload.EventID)
	_ = p.cancelByID(taskID)

	task := asynq.NewTask(model.TypeDraftRevert, b)
	_, err = p.Client.Enqueue(task, asynq.ProcessIn(delay), asynq.TaskID(taskID))
	return err
}

func (p *EventExpiryPublisher) CancelDraftRevert(eventID string) error {
	return p.cancelByID(fmt.Sprintf("draft_revert:%s", eventID))
}

func (p *EventExpiryPublisher) EnqueueUpdatedDraftRevert(payload *model.UpdatedDraftRevertPayload, revertAt time.Time) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	delay := time.Until(revertAt)
	if delay <= 0 {
		delay = 0
	}

	taskID := fmt.Sprintf("updated_draft_revert:%s", payload.UpdatedEventID)
	_ = p.cancelByID(taskID)

	task := asynq.NewTask(model.TypeUpdatedDraftRevert, b)
	_, err = p.Client.Enqueue(task, asynq.ProcessIn(delay), asynq.TaskID(taskID))
	return err
}

func (p *EventExpiryPublisher) CancelUpdatedDraftRevert(updatedEventID string) error {
	return p.cancelByID(fmt.Sprintf("updated_draft_revert:%s", updatedEventID))
}

func (p *EventExpiryPublisher) cancelByID(taskID string) error {
	queues, err := p.Inspector.Queues()
	if err != nil {
		return err
	}

	for _, q := range queues {
		scheduled, err := p.Inspector.ListScheduledTasks(q, asynq.PageSize(100))
		if err != nil {
			continue
		}
		for _, t := range scheduled {
			if t.ID == taskID {
				return p.Inspector.DeleteTask(q, taskID)
			}
		}
	}

	return nil
}
