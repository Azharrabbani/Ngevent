package model

// Type is a string value that indicates the type of the task.
// A list of task types.
const (
	TypeVerifiedUser = "user:verified"
	TypeVerifiedOTP  = "otp:verified"
)

// Task payload for any unverified user related tasks.
type UnverifiedUserPayload struct {
	UserID string
}

// Task payload for expired otp
type OTPPayload struct {
	OTPID string
}