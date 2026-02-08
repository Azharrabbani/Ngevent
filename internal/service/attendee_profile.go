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
	"github.com/nyaruka/phonenumbers"
)

type AttendeeProfileService struct {
	AttendeeRepo repository.AttendeeProfilesRepo
}

func NewAttendeeProfileService(attendeeRepo repository.AttendeeProfilesRepo) *AttendeeProfileService {
	return &AttendeeProfileService{AttendeeRepo: attendeeRepo}
}

var (
	profileUploadPath = "./storage/profiles"
)

func (s *AttendeeProfileService) Create(profile *dto.CreateProfileReq) error {
	// Validate image
	if err := helper.ValidateImage(profile.PhotoProfile); err != nil {
		return err
	}

	// Parse the phone code
	phoneNumber := phonenumber.ParseWithLandLine(profile.PhoneNumber, profile.ISO)
	if phoneNumber == "" {
		return errors.New("invalid phone number")
	}

	// Validate the phone number
	num, err := phonenumbers.Parse(profile.PhoneNumber, profile.ISO)
	if err != nil {
		return errors.New("invalid phone number format")
	}

	if !phonenumbers.IsValidNumber(num) {
		return errors.New("invalid phone number")
	}

	// Make sure the phone number is mobile phone number
	numberType := phonenumbers.GetNumberType(num)
	if numberType != phonenumbers.MOBILE && numberType != phonenumbers.FIXED_LINE_OR_MOBILE {
		return errors.New("phone number must be a mobile number")
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
	if profile.PhotoProfile != nil {
		_, fileName, err := helper.SaveToLocal(profile.PhotoProfile, profileUploadPath)
		if err != nil {
			return err
		}

		newProfile.PhotoProfile = &fileName

		return s.AttendeeRepo.Create(newProfile)
	}

	// Save Profile
	return s.AttendeeRepo.Create(newProfile)
}

func (s *AttendeeProfileService) FindByID(id string) (*dto.AttendeeProfilesResponse, error) {
	profile, err := s.AttendeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	attendee := toProfileResponse(profile)

	return attendee, nil
}

func (s *AttendeeProfileService) FindByUserID(id string) (*dto.AttendeeProfilesResponse, error) {
	profile, err := s.AttendeeRepo.FindByUserID(id)
	if err != nil {
		return nil, errors.New("profile not found")
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

	// Only validate user can update
	if userID != profile.UserID {
		return fiber.StatusUnauthorized, errors.New("unauthorized action")
	}

	oldPhoto := fmt.Sprintf("%s/%s", profileUploadPath, *profile.PhotoProfile)

	fmt.Println("old photo", oldPhoto)

	// Validate image
	if err := helper.ValidateImage(file); err != nil {
		return fiber.StatusBadRequest, err
	}

	// Save to local
	_, fileName, err := helper.SaveToLocal(file, profileUploadPath)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Remove old photo
	if err := os.Remove(oldPhoto); err != nil {
		log.Printf("failed to remove file from local %v\n", err)
	}

	if err := s.AttendeeRepo.UpdatePhotoProfile(userID, fileName); err != nil {
		return fiber.StatusBadRequest, err
	}

	return 0, nil
}

func (s *AttendeeProfileService) UpdateProfile(userID string, req *dto.UpdateProfileReq) (int, error) {
	// Validate user
	profile, err := s.AttendeeRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	// Only validate user can update
	if userID != profile.UserID {
		return fiber.StatusUnauthorized, errors.New("unauthorized action")
	}

	// Parse the phone code
	phoneNumber := phonenumber.ParseWithLandLine(req.PhoneNumber, req.ISO)

	// Validate the phone number
	num, err := phonenumbers.Parse(req.PhoneNumber, req.ISO)
	if err != nil {
		return fiber.StatusBadRequest, errors.New("invalid phone number format")
	}

	if !phonenumbers.IsValidNumber(num) {
		return fiber.StatusBadRequest, errors.New("invalid phone number")
	}

	// Make sure the phone number is mobile phone number
	numberType := phonenumbers.GetNumberType(num)
	if numberType != phonenumbers.MOBILE && numberType != phonenumbers.FIXED_LINE_OR_MOBILE {
		return fiber.StatusBadRequest, errors.New("phone number must be a mobile number")
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
