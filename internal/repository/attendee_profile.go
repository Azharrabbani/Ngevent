package repository

import (
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

type AttendeeProfileRepository struct {
	db *gorm.DB
}

func NewAttendeeProfileRepository(db *gorm.DB) AttendeeProfilesRepo {
	return &AttendeeProfileRepository{db: db}
}

// Create implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) Create(profile *model.AttendeeProfiles) error {
	return r.db.Create(profile).Error
}

// Delete implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.AttendeeProfiles{}).Error
}

// FindByID implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) FindByID(id string) (*model.AttendeeProfiles, error) {
	var profile *model.AttendeeProfiles

	if err := r.db.Preload("User").Where("id = ?", id).First(&profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

// FindByUserID implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) FindByUserID(id string) (*model.AttendeeProfiles, error) {
	var profile *model.AttendeeProfiles

	if err := r.db.Preload("User").Where("user_id = ?", id).First(&profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

// Update implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) Update(profile *model.AttendeeProfiles) error {
	return r.db.Save(profile).Error
}

// UpdatePhotoProfile implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) UpdatePhotoProfile(userID string, photo string) error {
	return r.db.Model(&model.AttendeeProfiles{}).
		Where("user_id = ?", userID).
		Updates(&model.AttendeeProfiles{
			PhotoProfile: &photo,
			UpdatedAt:    time.Now().UTC()}).Error
}
