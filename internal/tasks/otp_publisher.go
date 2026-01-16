package tasks

import (
	"encoding/json"
	"fmt"
	"ngevent/internal/model"
	"time"

	"github.com/hibiken/asynq"
)

type OtpTaskPublisher struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewOTPTaskPublisher(client *asynq.Client, inspector *asynq.Inspector) *OtpTaskPublisher {
	return &OtpTaskPublisher{
		Client:    client,
		Inspector: inspector,
	}
}

func (t *OtpTaskPublisher) EnqueueOTPVerification(taskType string, payload *model.OTPPayload) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, p)

	_, err = t.Client.Enqueue(task, asynq.ProcessIn(6*time.Minute), asynq.TaskID(payload.OTPID))
	if err != nil {
		return err
	}

	return nil
}

func (t *OtpTaskPublisher) CancelOTPVerification(id string) error {
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
			if id == task.ID {
				t.Inspector.DeleteTask(q, id)
			}
		}
	}

	return fmt.Errorf("task with ID %s not found in any queue", id)
}
