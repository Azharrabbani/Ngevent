package model

// Type is a string value that indicates the type of the task.
// A list of task types.
const (
	TypeVerifiedUser = "user:verified"
	TypeVerifiedOTP  = "otp:verified"

	// Email payload
	TypeEmailForgetPassword           = "email:user:forgetPassword"
	TypeEMailVerify                   = "email:user:verify"
	TypeEmailAdminVerification        = "email:admin:verify"
	TypeEmailOrganizerProfile         = "email:organizer:profile"
	TypeEmailOrganizerProfileVerified = "email:organizer:profile:verified"
	TypeEmailOrganizerProfileRejected = "email:organizer:profile:rejected"

	// Event email
	TypeEventAdminNotification  = "event:admin:notification"
	TypeEventEONotification     = "event:eo:notification"
	TypeEventEOVerification     = "event:eo:verification"
	TypeEventUpdateNotification = "updated_event:eo:notification"

	TypeEventExpired        = "event:expired"
	TypeUpdatedEventExpired = "updated_event:expired"
)

type EventExpiredPayload struct {
	EventID string `json:"event_id"`
}

type UpdatedEventExpiredPayload struct {
	UpdatedEventID string `json:"updated_event_id"`
	EventID        string `json:"event_id"`
}

// Task payload for any unverified user related tasks.
type UnverifiedUserPayload struct {
	UserID string
}

// Task payload for expired otp
type OTPPayload struct {
	OTPID string
}

// Task payload for email
type EmailPayload struct {
	To        string
	Name      string
	UserEmail string
	Action    string
	OTP       string
	OTPID     string
}

type RejectedEmailPayload struct {
	To     string
	Name   string
	Reason string
}

type EventActionReq string

const (
	Create EventActionReq = "create"
	Update EventActionReq = "update"
)

type EventEmailPayload struct {
	To        string
	EOName    string
	EOEmail   string
	EventName string
	Status    string
	Reason    string
}
