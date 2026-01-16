package model

import "time"

type OtpVerifications struct {
	ID               string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID           string    `json:"user_id"`
	OTP              string    `json:"otp"`
	IsUsed           bool      `json:"is_used"`
	ExpiredAt        time.Time `json:"expired_at" gorm:"default:now()"`
	TypeVerification string    `json:"type_verification"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"default:now()"`
	User             Users     `gorm:"foreignKey:UserID"`
}

func (OtpVerifications) TableName() string {
	return "otp_verifications"
}
