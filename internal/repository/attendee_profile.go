package repository

import (
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AttendeeProfileRepository struct {
	db *gorm.DB
}

// FindAll implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) (*model.PaginationRow[*dto.AttendeeProfilesResponse], error) {
	var profiles []*model.AttendeeProfiles

	query := r.db.Scopes(filterAttendeeList(filter))

	if err := query.Preload("User").
		Scopes(Paginate(profiles, &pagination, query)).
		Find(&profiles).Error; err != nil {
		return nil, err
	}

	attendees := toAttendeesResponse(profiles)

	return &model.PaginationRow[*dto.AttendeeProfilesResponse]{
		Pagination: pagination,
		Rows:       attendees,
	}, nil

}

func NewAttendeeProfileRepository(db *gorm.DB) AttendeeProfilesRepo {
	return &AttendeeProfileRepository{db: db}
}

// HasProfile implements AttendeeProfilesRepo.
func (r *AttendeeProfileRepository) HasProfile(userID string) (bool, error) {
	if err := r.db.
		Where("user_id = ?", userID).
		First(&model.AttendeeProfiles{}).Error; err != nil {
		return false, err
	}

	return true, nil
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

func filterAttendeeList(filter *dto.FilterProfileReq) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Filter == nil {
			return db
		}

		query := "%" + strings.ToLower(*filter.Filter) + "%"

		return db.Joins("JOIN users ON users.id = attendee_profiles.user_id").
			Where(
				db.Where("LOWER(users.email) LIKE ?", query).
					Or("LOWER(attendee_profiles.name) LIKE ?", query).
					Or("LOWER(attendee_profiles.username) LIKE ?", query),
			)
	}
}

func toAttendeesResponse(profiles []*model.AttendeeProfiles) []*dto.AttendeeProfilesResponse {
	var attendees []*dto.AttendeeProfilesResponse

	for _, profile := range profiles {
		attendees = append(attendees, &dto.AttendeeProfilesResponse{
			ID:           profile.ID,
			UserID:       profile.UserID,
			Email:        profile.User.Email,
			Name:         profile.Name,
			Username:     profile.Username,
			PhotoProfile: fmt.Sprintf("http://localhost:8080/api/v1/attendee/photo/%s", helper.StringValue(profile.PhotoProfile)),
			PhoneNumber:  profile.PhoneNumber,
			Country:      profile.Country,
			Address:      profile.Address,
		})
	}

	return attendees
}
