package repository

import "ngevent/internal/model"

type OtpRepo interface {
	Create(otp *model.OtpVerifications) (*model.OtpVerifications, error)
	FindByID(id string) (*model.OtpVerifications, error)
	FindByUserID(id string) (*model.OtpVerifications, error)
	Update(otp *model.OtpVerifications) (*model.OtpVerifications, error)
	Delete(id string) error
}
