package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type EmailTaskPublisher struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewEmailTaskPublisher(client *asynq.Client, inspector *asynq.Inspector) *EmailTaskPublisher {
	return &EmailTaskPublisher{
		Client:    client,
		Inspector: inspector,
	}
}

func (t *EmailTaskPublisher) Enqueue(taskType string, payload interface{}) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, p)

	_, err = t.Client.Enqueue(
		task,
		asynq.MaxRetry(3),
		asynq.Timeout(30*time.Second),
		asynq.Queue("emails"),
	)

	return err
}
