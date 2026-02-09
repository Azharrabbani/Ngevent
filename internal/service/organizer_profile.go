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

type OrganizerProfileService struct {
	OrganizerRepo      repository.OrganizerProfileRepo
	UserRepo           repository.UsersRepo
	EmailTaskPublisher NewTaskEmail
}

func NewOrganizerProfileService(
	organizerRepo repository.OrganizerProfileRepo,
	userRepo repository.UsersRepo,
	emailTaskPublisher NewTaskEmail,
) *OrganizerProfileService {
	return &OrganizerProfileService{
		OrganizerRepo:      organizerRepo,
		UserRepo:           userRepo,
		EmailTaskPublisher: emailTaskPublisher,
	}
}

func (s *OrganizerProfileService) CreateProfile(profile *dto.CreateOrganizerProfileReq) error {
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

	country := phonenumber.GetISO3166ByNumber(phoneNumber, true)

	newProfile := &model.OrganizerProfiles{
		UserID:      profile.UserID,
		Name:        profile.Name,
		Address:     profile.Address,
		PhoneNumber: fmt.Sprintf("+%s", phoneNumber),
		Country:     country.CountryName,
		SocialMedias: model.OrganizerSocialMedia{
			Email:     profile.SocialMedia.Email,
			Instagram: profile.SocialMedia.Instagram,
		},
		CompanyDetail: model.OrganizerCompDetail{
			Description: profile.CompanyDetail.Description,
			NPWP:        profile.CompanyDetail.NPWP,
			NIB:         profile.CompanyDetail.NIB,
		},
	}

	// Save photo profile to local storage
	if profile.PhotoProfile != nil {
		_, fileName, err := helper.SaveToLocal(profile.PhotoProfile, profileUploadPath)
		if err != nil {
			return err
		}

		newProfile.PhotoProfile = &fileName

		if err := s.OrganizerRepo.Create(newProfile); err != nil {
			return err
		}

		payload := &model.EmailPayload{
			To:   newProfile.User.Email,
			Name: newProfile.Name,
		}

		// Send email to organizer
		s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfile, payload)

		return nil
	}

	// Save Profile
	if err := s.OrganizerRepo.Create(newProfile); err != nil {
		return err
	}

	payload := &model.EmailPayload{
		To:   newProfile.User.Email,
		Name: newProfile.Name,
	}

	// Send email to organizer
	s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfile, payload)

	return nil
}

func (s *OrganizerProfileService) FindByID(id string) (*dto.OrganizerProfilesResponse, error) {
	profile, err := s.OrganizerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	organizer := toOrganizerProfileResponse(profile)

	return organizer, nil
}

func (s *OrganizerProfileService) FindByUserID(userID string) (*dto.OrganizerProfilesResponse, error) {
	profile, err := s.OrganizerRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	organizer := toOrganizerProfileResponse(profile)

	return organizer, nil
}

func (s *OrganizerProfileService) FindByCountry(
	country string,
	pagination model.Pagination,
) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	profiles, err := s.OrganizerRepo.FindByCountry(country, pagination)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (s *OrganizerProfileService) VerifiedProfile(id string) error {
	// Check user is exist
	profile, err := s.OrganizerRepo.FindByID(id)
	if err != nil {
		return errors.New("profile not found")
	}

	if profile.IsVerified {
		return errors.New("profile already verified")
	}

	if err := s.OrganizerRepo.VerifiedProfile(id); err != nil {
		return errors.New("failed to verified profile")
	}

	// Send email
	payload := &model.EmailPayload{
		To:   profile.User.Email,
		Name: profile.Name,
	}
	s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfileVerified, payload)

	return nil
}

func (s *OrganizerProfileService) RejectProfile(id string) error {
	// Check user is exist
	profile, err := s.OrganizerRepo.FindByID(id)
	if err != nil {
		return errors.New("profile not found")
	}

	if err := s.OrganizerRepo.VerifiedProfile(id); err != nil {
		return errors.New("failed to verified profile")
	}

	// Send email
	payload := &model.EmailPayload{
		To:   profile.User.Email,
		Name: profile.Name,
	}
	s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfileVerified, payload)

	return nil
}

func (s *OrganizerProfileService) UpdatePhotoProfile(file *multipart.FileHeader, userID string) (int, error) {
	// Get old photo profile
	profile, err := s.OrganizerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	// Only validate user can update
	if userID != profile.UserID {
		return fiber.StatusUnauthorized, errors.New("unauthorized action")
	}

	oldPhoto := fmt.Sprintf("%s/%s", profileUploadPath, *profile.PhotoProfile)

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

	if err := s.OrganizerRepo.UpdatePhotoProfile(userID, fileName); err != nil {
		return fiber.StatusBadRequest, err
	}

	return 0, nil
}

func (s *OrganizerProfileService) UpdateProfile(userID string, req *dto.UpdateOrganizerProfileReq) (int, error) {
	// Validate user
	profile, err := s.OrganizerRepo.FindByUserID(userID)
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

	country := phonenumber.GetISO3166ByNumber(phoneNumber, true)

	profile.Name = req.Name
	profile.PhoneNumber = req.PhoneNumber
	profile.Country = country.CountryName
	profile.Address = req.Address
	profile.SocialMedias.Email = req.SocialMedia.Email
	profile.SocialMedias.Instagram = req.SocialMedia.Instagram
	profile.CompanyDetail.Description = req.CompanyDetail.Description
	profile.CompanyDetail.NPWP = req.CompanyDetail.NPWP
	profile.CompanyDetail.NIB = req.CompanyDetail.NIB
	profile.UpdatedAt = time.Now().UTC()

	if err := s.OrganizerRepo.Update(profile); err != nil {
		return fiber.StatusBadRequest, nil
	}

	return 0, nil
}

func toOrganizerProfileResponse(profile *model.OrganizerProfiles) *dto.OrganizerProfilesResponse {
	return &dto.OrganizerProfilesResponse{
		ID:           profile.ID,
		UserID:       profile.UserID,
		Email:        profile.User.Email,
		Name:         profile.Name,
		PhotoProfile: profile.PhotoProfile,
		PhoneNumber:  profile.PhoneNumber,
		Country:      profile.Country,
		Address:      profile.Address,
		SocialMedia: dto.OrganizerSocialMediaReq{
			Email:     profile.SocialMedias.Email,
			Instagram: profile.SocialMedias.Instagram,
		},
		CompanyDetail: dto.OrganizerCompDetailReq{
			Description: profile.CompanyDetail.Description,
			NPWP:        profile.CompanyDetail.NPWP,
			NIB:         profile.CompanyDetail.NIB,
		},
	}
}
