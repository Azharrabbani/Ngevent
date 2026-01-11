package repository

import (
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type OtpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) *OtpRepository {
	return &OtpRepository{db: db}
}

func (r *OtpRepository) Create(otp *model.OtpVerifications) (*model.OtpVerifications, error) {
	if err := r.db.Create(otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *OtpRepository) FindByID(id string) (*model.OtpVerifications, error) {
	var otp *model.OtpVerifications

	if err := r.db.Where("id = ?", id).First(&otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *OtpRepository) FindByUserID(id string) (*model.OtpVerifications, error) {
	var otp *model.OtpVerifications

	if err := r.db.Where("user_id = ?", id).First(&otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *OtpRepository) Update(otp *model.OtpVerifications) (*model.OtpVerifications, error) {

	if err := r.db.Save(otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *OtpRepository) Delete(id string) error {
	var otp *model.OtpVerifications

	return r.db.Where("id = ?", id).Delete(&otp).Error
}
