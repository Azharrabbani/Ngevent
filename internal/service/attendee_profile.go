package service

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
)

type AttendeeProfileService struct {
	AttendeeRepo repository.AttendeeProfilesRepo
}

func NewAttendeeProfileService(attendeeRepo repository.AttendeeProfilesRepo) *AttendeeProfileService {
	return &AttendeeProfileService{AttendeeRepo: attendeeRepo}
}

func (s *AttendeeProfileService) Create(profile *model.AttendeeProfiles) error {
	return s.AttendeeRepo.Create(profile)
}

func (s *AttendeeProfileService) FindByID(id string) (*dto.AttendeeProfilesResponse, error) {
	profile, err := s.AttendeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	attendee := toProfileResponse(profile)

	return attendee, nil
}

func (s *AttendeeProfileService) FindByUserID(id string) (*dto.AttendeeProfilesResponse, error) {
	profile, err := s.AttendeeRepo.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	attendee := toProfileResponse(profile)

	return attendee, nil
}

func (s *AttendeeProfileService) UpdateProfile()

func toProfileResponse(profile *model.AttendeeProfiles) *dto.AttendeeProfilesResponse {
	return &dto.AttendeeProfilesResponse{
		ID:           profile.ID,
		UserID:       profile.UserID,
		Email:        profile.User.Email,
		Name:         profile.Name,
		Username:     profile.Username,
		PhotoProfile: profile.PhotoProfile,
		PhoneNumber:  profile.PhoneNumber,
		Country:      profile.Country,
		Address:      profile.Address,
	}

}
