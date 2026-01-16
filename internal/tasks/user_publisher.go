package tasks

import (
	"encoding/json"
	"fmt"
	"ngevent/internal/model"
	"time"

	"github.com/hibiken/asynq"
)

type UserTaskPublisher struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewUserTaskPublisher(client *asynq.Client, inspector *asynq.Inspector) *UserTaskPublisher {
	return &UserTaskPublisher{
		Client:    client,
		Inspector: inspector,
	}
}

func (t *UserTaskPublisher) EnqueueUnverifiedUser(taskType string, payload *model.UnverifiedUserPayload) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, p)

	_, err = t.Client.Enqueue(task, asynq.ProcessIn(6*time.Minute), asynq.TaskID(payload.UserID))
	if err != nil {
		return err
	}

	return nil
}

func (t *UserTaskPublisher) CancelUnverifiedUser(id string) error {
	queues, err := t.Inspector.Queues()
	if err != nil {
		return err
	}

	for _, q := range queues {
		scheduledTask, err := t.Inspector.ListScheduledTasks(q, asynq.PageSize(100))
		if err != nil {
			continue
		}

		for _, task := range scheduledTask {
			if task.ID == id {
				return t.Inspector.DeleteTask(q, id)
			}
		}
	}

	return fmt.Errorf("task with ID %s not found in any queue", id)
}
