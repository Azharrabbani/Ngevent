package tasks

import (
	"context"
	"encoding/json"
	"log"
	"ngevent/internal/model"
	"ngevent/internal/service"

	"github.com/hibiken/asynq"
)

type OtpTaskHandler struct {
	AuthService *service.AuthService
}

func NewOTPTaskHandler(authService *service.AuthService) *OtpTaskHandler {
	return &OtpTaskHandler{AuthService: authService}
}

func (h *OtpTaskHandler) HandlerUnusedOTP(ctx context.Context, t *asynq.Task) error {
	var p *model.OTPPayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	if err := h.AuthService.DeleteUnusedOTP(p.OTPID); err != nil {
		return err
	}

	log.Printf(" [*] \"%s\" has been deleted.", p.OTPID)

	return nil

}
