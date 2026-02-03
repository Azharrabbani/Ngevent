package tasks

import (
	"context"
	"encoding/json"
	"ngevent/internal/model"
	"ngevent/internal/utils"

	"github.com/hibiken/asynq"
)

type EmailTaskHandler struct{}

func NewEmailTaskHandler() *EmailTaskHandler {
	return &EmailTaskHandler{}
}

// Verification email
func (h *EmailTaskHandler) HandlerUserVerification(ctx context.Context, t *asynq.Task) error {
	var p *model.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.VerifyEmailMail(p.OTP, p.To, p.OTPID)
}

// Forget password email
func (h *EmailTaskHandler) HandlerUserForgetPassword(ctx context.Context, t *asynq.Task) error {
	var p *model.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.ForgotPasswordMail(p.To, p.OTPID)
}

