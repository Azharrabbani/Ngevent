package tasks

import (
	"context"
	"encoding/json"
	"log"
	"ngevent/internal/model"
	"ngevent/internal/service"

	"github.com/hibiken/asynq"
)

type UserTaskHandler struct {
	UserService *service.UserService
}

func NewUserTaskHandler(userService *service.UserService) *UserTaskHandler {
	return &UserTaskHandler{UserService: userService}
}

func (h *UserTaskHandler) HandlerUnverifiedTask(ctx context.Context, t *asynq.Task) error {
	var p *model.UnverifiedUserPayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	if err := h.UserService.DeleteUnverifiedUser(p.UserID); err != nil {
		return err
	}

	log.Printf(" [*] \"%s\" has been deleted.", p.UserID)

	return nil

}
