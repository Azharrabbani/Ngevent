package helper

import (
	"errors"
	"fmt"
	"math/rand"
	"ngevent/internal/model"
	"time"
)

func NewOTP(otpCode, userID, action string, exp time.Time) *model.OtpVerifications {
	return &model.OtpVerifications{
		OTP:              otpCode,
		UserID:           userID,
		TypeVerification: action,
		ExpiredAt:        exp,
	}
}
func GenerateOTP() (string, error) {
	rand.Seed(time.Now().UnixNano())

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	if otp == "" {
		return "", errors.New("failed to load otp")
	}

	return otp, nil
}
