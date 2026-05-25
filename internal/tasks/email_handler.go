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

	return utils.VerifyEmailMail(p.OTP, p.To)
}

// Forget password email
func (h *EmailTaskHandler) HandlerUserForgetPassword(ctx context.Context, t *asynq.Task) error {
	var p *model.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.ForgotPasswordMail(p.To, p.OTPID)
}

// Email to admin for verify the organizer's profile
func (h *EmailTaskHandler) HandlerAdminVerifyProfile(ctx context.Context, t *asynq.Task) error {
	var p *model.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerProfileAdminNotificationEmail(p.To, p.Name, p.UserEmail, p.Action)
}

// Organize's profile register
func (h *EmailTaskHandler) HandlerOrganizerProfileNotif(ctx context.Context, t *asynq.Task) error {
	var p *model.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerProfileVerificationEmail(p.To, p.Name)
}

// Organizer's profile verified
func (h *EmailTaskHandler) HandlerOrganizerProfileVerified(ctx context.Context, t *asynq.Task) error {
	var p *model.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerProfileVerifiedEmail(p.To, p.Name)
}

// Organizer's rejected profile
func (h *EmailTaskHandler) HandlerOrganizerProfileRejected(ctx context.Context, t *asynq.Task) error {
	var p *model.RejectedEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerProfileRejectedEmail(p.To, p.Name, p.Reason)
}

// ============ Email event handler ===============
// Admin new event notification
func (h *EmailTaskHandler) HandlerEventAdminNotification(ctx context.Context, t *asynq.Task) error {
	var p *model.EventEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.AdminEventNotification(p.To, p.EOName, p.EventName, p.EOEmail, p.Status)
}

// Organizer new event notification
func (h *EmailTaskHandler) HandlerEventOrganizerNotification(ctx context.Context, t *asynq.Task) error {
	var p *model.EventEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerEventNotification(p.To, p.EOName, p.EventName)
}

// Organizer event verification
func (h *EmailTaskHandler) HandlerEventOrganizerVerification(ctx context.Context, t *asynq.Task) error {
	var p *model.EventEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerEventVerification(p.To, p.EOName, p.EventName, p.Status, p.Reason)
}

func (h *EmailTaskHandler) HandlerUpdateEventOrganizerNotif(ctx context.Context, t *asynq.Task) error {
	var p *model.EventEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return utils.OrganizerUpdatedEventNotif(p.To, p.EOName, p.EventName, p.Status, p.Reason)
}
