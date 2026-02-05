package service

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"
	"os"
	"time"

	"github.com/dongri/phonenumber"
	"github.com/gofiber/fiber/v2"
)

type AttendeeProfileService struct {
	AttendeeRepo repository.AttendeeProfilesRepo
}

func NewAttendeeProfileService(attendeeRepo repository.AttendeeProfilesRepo) *AttendeeProfileService {
	return &AttendeeProfileService{AttendeeRepo: attendeeRepo}
}

var (
	profileUploadPath = "./storage/profile"
)

func (s *AttendeeProfileService) Create(file *multipart.FileHeader, profile *dto.CreateProfileReq) error {
	// Validate image
	if err := helper.ValidateImage(file); err != nil {
		return err
	}

	// Parse the phone code
	phoneNumber := phonenumber.ParseWithLandLine(profile.PhoneNumber, profile.ISO)
	if phoneNumber == "" {
		return errors.New("invalid phone number")
	}

	country := phonenumber.GetISO3166ByNumber(phoneNumber, true)

	newProfile := &model.AttendeeProfiles{
		UserID:      profile.UserID,
		Name:        profile.Name,
		Username:    profile.Username,
		Address:     profile.Address,
		PhoneNumber: fmt.Sprintf("+%s", phoneNumber),
		Country:     country.CountryName,
	}

	// Save photo profile to local storage
	if file != nil {
		photoPath, _, err := helper.SaveToLocal(file, profileUploadPath)
		if err != nil {
			return err
		}

		newProfile.PhotoProfile = &photoPath
	}

	// Save Profile
	return s.AttendeeRepo.Create(newProfile)
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

func (s *AttendeeProfileService) UpdatePhotoProfile(file *multipart.FileHeader, userID string) (int, error) {
	// Get old photo profile
	profile, err := s.AttendeeRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	oldPhoto := profile.PhotoProfile

	// Validate image
	if err := helper.ValidateImage(file); err != nil {
		return fiber.StatusBadRequest, err
	}

	// Save to local
	newPath, _, err := helper.SaveToLocal(file, profileUploadPath)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Remove old photo
	if err := os.Remove(*oldPhoto); err != nil {
		log.Printf("failed to remove file from local %v\n", err)
	}

	if err := s.AttendeeRepo.UpdatePhotoProfile(userID, newPath); err != nil {
		return fiber.StatusBadRequest, err
	}

	return 0, nil
}

func (s *AttendeeProfileService) UpdateProfile(req *dto.UpdateProfileReq) (int, error) {
	// Validate user
	profile, err := s.AttendeeRepo.FindByUserID(req.UserID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	// Validate phone number
	// Parse the phone code
	phoneNumber := phonenumber.ParseWithLandLine(req.PhoneNumber, req.ISO)
	if phoneNumber == "" {
		return fiber.StatusBadRequest, errors.New("invalid phone number")
	}

	country := phonenumber.GetISO3166ByNumber(phoneNumber, true)

	profile.Name = req.Name
	profile.Username = req.Username
	profile.PhoneNumber = req.PhoneNumber
	profile.Country = country.CountryName
	profile.Address = req.Address
	profile.UpdatedAt = time.Now().UTC()

	if err := s.AttendeeRepo.Update(profile); err != nil {
		return fiber.StatusBadRequest, nil
	}

	return 0, nil
}

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
