package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AttendeeProfileService struct {
	AttendeeRepo repository.AttendeeProfilesRepo
	rdb          *redis.Client
}

func NewAttendeeProfileService(attendeeRepo repository.AttendeeProfilesRepo, rdb *redis.Client) *AttendeeProfileService {
	return &AttendeeProfileService{
		AttendeeRepo: attendeeRepo,
		rdb:          rdb,
	}
}

var (
	profileUploadPath = "./storage/profiles"
)

var attendeeCache []string = []string{
	"attendee:all:*",
}

func (s *AttendeeProfileService) Create(profile *dto.CreateAttendeeProfileReq) error {
	// Validate image
	if err := helper.ValidateImage(profile.PhotoProfile); err != nil {
		return err
	}

	// Validate the phone number
	phonenumber, country, err := utils.ValidatePhoneCode(profile.PhoneNumber, profile.ISO)
	if err != nil {
		return err
	}

	newProfile := &model.AttendeeProfiles{
		UserID:      profile.UserID,
		Name:        profile.Name,
		Username:    profile.Username,
		Address:     profile.Address,
		PhoneNumber: fmt.Sprintf("+%s", phonenumber),
		Country:     country,
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
	if err := s.AttendeeRepo.Create(newProfile); err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, attendeeCache)

	return nil
}

func (s *AttendeeProfileService) FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) (*model.PaginationRow[*dto.AttendeeProfilesResponse], error) {
	var attendess *model.PaginationRow[*dto.AttendeeProfilesResponse]

	cacheKey := fmt.Sprintf("attendee:all:%d:%d:%s:%s", pagination.Limit, pagination.Page, pagination.Sort, filter.Filter)

	// Tru get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &attendess)
	}

	if attendess == nil {
		// if cache miss, get from db
		attendess, err = s.AttendeeRepo.FindAll(pagination, filter)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(attendess); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return attendess, nil
}

func (s *AttendeeProfileService) HasProfile(userID string) (bool, error) {
	hasProfile, err := s.AttendeeRepo.HasProfile(userID)
	if err != nil {
		return false, err
	}

	return hasProfile, nil
}

func (s *AttendeeProfileService) FindByID(id string) (*dto.AttendeeProfilesResponse, error) {
	profile, err := s.AttendeeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	attendee := toAttendeeProfileResponse(profile)

	return attendee, nil
}

func (s *AttendeeProfileService) FindByUserID(id string) (*dto.AttendeeProfilesResponse, error) {
	profile, err := s.AttendeeRepo.FindByUserID(id)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	attendee := toAttendeeProfileResponse(profile)

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

	// Validate image
	if err := helper.ValidateImage(file); err != nil {
		return fiber.StatusBadRequest, err
	}

	// Save to local
	_, fileName, err := helper.SaveToLocal(file, profileUploadPath)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	if err := s.AttendeeRepo.UpdatePhotoProfile(userID, fileName); err != nil {
		return fiber.StatusBadRequest, err
	}

	// Remove old photo
	if err := os.Remove(oldPhoto); err != nil {
		log.Printf("failed to remove file from local %v\n", err)
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, attendeeCache)

	return 0, nil
}

func (s *AttendeeProfileService) UpdateProfile(userID string, req *dto.UpdateAttendeeProfileReq) (int, error) {
	// Validate user
	profile, err := s.AttendeeRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	// Only validate user can update
	if userID != profile.UserID {
		return fiber.StatusUnauthorized, errors.New("unauthorized action")
	}

	// Validate the phone number
	phonenumber, country, err := utils.ValidatePhoneCode(req.PhoneNumber, req.ISO)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	profile.Name = req.Name
	profile.Username = req.Username
	profile.PhoneNumber = fmt.Sprintf("+%s", phonenumber)
	profile.Country = country
	profile.Address = req.Address
	profile.UpdatedAt = time.Now().UTC()

	if err := s.AttendeeRepo.Update(profile); err != nil {
		return fiber.StatusBadRequest, nil
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, attendeeCache)

	return 0, nil
}

func toAttendeeProfileResponse(profile *model.AttendeeProfiles) *dto.AttendeeProfilesResponse {
	return &dto.AttendeeProfilesResponse{
		ID:           profile.ID,
		UserID:       profile.UserID,
		Email:        profile.User.Email,
		Name:         profile.Name,
		Username:     profile.Username,
		PhotoProfile: fmt.Sprintf("http://localhost:8080/api/v1/attendee/photo/%s", helper.StringValue(profile.PhotoProfile)),
		PhoneNumber:  profile.PhoneNumber,
		Country:      profile.Country,
		Address:      profile.Address,
	}
}
