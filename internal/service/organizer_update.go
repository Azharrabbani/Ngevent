package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
)

type OrganizerUpdateService struct {
	UserRepo             repository.UsersRepo
	OrganizerProfileRepo repository.OrganizerProfileRepo
	OrganizerUpdateRepo  repository.OrganizerProfileUpdateRepo
	EmailTaskPublisher   NewTaskEmail
	rdb                  *redis.Client
}

func NewOrganizerUpdateService(
	userRepo repository.UsersRepo,
	organizerProfileRepo repository.OrganizerProfileRepo,
	organizerUpdateRepo repository.OrganizerProfileUpdateRepo,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *OrganizerUpdateService {
	return &OrganizerUpdateService{
		UserRepo:             userRepo,
		OrganizerProfileRepo: organizerProfileRepo,
		OrganizerUpdateRepo:  organizerUpdateRepo,
		EmailTaskPublisher:   emailTaskPublisher,
		rdb:                  rdb,
	}
}

var (
	npwpStagePath = "./storage/npwp/stage"
	nibStagePath  = "./storage/nib/stage"
)

var organizerUpdateCache []string = []string{
	"organizer:update:all:*",
}

func (s *OrganizerUpdateService) Validate(req *dto.ValidateUpdateReq) error {
	// Start transaction
	profileX := s.OrganizerProfileRepo.GetDB().Begin()
	updateX := s.OrganizerUpdateRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			profileX.Rollback()
			updateX.Rollback()
		}
	}()

	// Get the update data
	updateData, err := s.OrganizerUpdateRepo.FindByID(req.UpdateID)
	if err != nil {
		return errors.New("update data not found")
	}

	// Get profile
	profile, err := s.OrganizerProfileRepo.FindByID(updateData.ProfileID)
	if err != nil {
		return errors.New("profile not found")
	}

	isNPWPUpdated := updateData.NPWPDocument != ""
	isNIBUpdated := updateData.NIBDocument != ""
	isFilesUpdates := isNPWPUpdated || isNIBUpdated

	oldNPWPFile := fmt.Sprintf("%s/%s", npwpFilePath, profile.CompanyDetail.NPWPDocument)
	oldNIBFile := fmt.Sprintf("%s/%s", nibFilePath, profile.CompanyDetail.NIBDocument)

	// Determine the status
	switch status := req.Status; status {
	case "approved":
		// 1. When approves update the profile with the update data in staging
		profile.Name = updateData.Name
		profile.Slug = updateData.Slug
		profile.Status.Status = status
		profile.PhoneNumber = updateData.PhoneNumber
		profile.Country = updateData.Country
		profile.CompanyDetail.Description = updateData.Description
		profile.Address = updateData.Address
		profile.SocialMedias.Email = updateData.Email
		profile.SocialMedias.Instagram = updateData.Instagram
		profile.CompanyDetail.NPWPNumber = updateData.NPWPNumber
		profile.CompanyDetail.NIBNumber = updateData.NIBNumber
		profile.RequestUpdates = false

		// 2. [Opt] Copy file from staging to destination path
		if isFilesUpdates {
			if isNPWPUpdated {
				npwpSrcPath := fmt.Sprintf("%s/%s", npwpStagePath, updateData.NPWPDocument)
				dstPath := fmt.Sprintf("%s/%s", npwpFilePath, filepath.Base(updateData.NPWPDocument))
				npwpFile, err := helper.CopyFile(npwpSrcPath, dstPath)
				if err != nil {
					return err
				}
				profile.CompanyDetail.NPWPDocument = npwpFile
			}

			if isNIBUpdated {
				nibSrcPath := fmt.Sprintf("%s/%s", nibStagePath, updateData.NIBDocument)
				nibDstPath := fmt.Sprintf("%s/%s", nibFilePath, filepath.Base(updateData.NIBDocument))
				nibFile, err := helper.CopyFile(nibSrcPath, nibDstPath)
				if err != nil {
					return err
				}
				profile.CompanyDetail.NIBDocument = nibFile
			}
		}

		// 3. Update the data
		if err := profileX.Save(profile).Error; err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return err
		}

		// 4. Update staging status
		status := "approved"
		if err := s.OrganizerUpdateRepo.Validate(updateData.ID, status); err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return errors.New("validate failed")
		}

		// 5. [Opt] Remove old files
		if isFilesUpdates {
			if isNPWPUpdated {
				if err := os.Remove(oldNPWPFile); err != nil {
					profileX.Rollback()
					updateX.Rollback()
					return err
				}
			}
			if isNIBUpdated {
				if err := os.Remove(oldNIBFile); err != nil {
					profileX.Rollback()
					updateX.Rollback()
					return err
				}
			}
		}

		// 6. Send email to user
		EmailPayload := &model.EmailPayload{
			To:   profile.User.Email,
			Name: profile.Name,
		}
		s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfileVerified, EmailPayload)
	case "rejected":
		if req.Reason == "" {
			return errors.New("reason required")
		}

		// Update staging status
		status := "rejected"
		updateData.UpdatedAt = time.Now().UTC()

		if err := s.OrganizerUpdateRepo.Validate(updateData.ID, status); err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return errors.New("validate failed")
		}

		profile.RequestUpdates = false
		if err := profileX.Save(profile).Error; err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return err
		}

		// Send email to notify user
		emailPayload := &model.RejectedEmailPayload{
			To:     profile.User.Email,
			Name:   profile.Name,
			Reason: req.Reason,
		}
		s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfileRejected, emailPayload)
	default:
		return errors.New("validate must be one of approved or rejected")
	}

	if err := profileX.Commit().Error; err != nil {
		return err
	}

	if err := updateX.Commit().Error; err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, organizerUpdateCache)
	utils.InvalidateCache(s.rdb, organizerCache)

	return nil
}

func (s *OrganizerUpdateService) FindByID(id string) (*model.OrganizerProfilesUpdates, error) {
	return s.OrganizerUpdateRepo.FindByID(id)
}

func (s *OrganizerUpdateService) FindByProfileID(id string) (*dto.OrganizerUpdatesResponse, error) {
	organizerUpdate, err := s.OrganizerUpdateRepo.FindByProfileID(id)
	if err != nil {
		return nil, err
	}

	resp := toOrganizerUpdateResponse(organizerUpdate)

	return resp, nil
}

func (s *OrganizerUpdateService) FindUpdatesByProfileID(id string, pagination model.Pagination) (*model.PaginationRow[*model.OrganizerProfilesUpdates], error) {
	var organizerUpdate *model.PaginationRow[*model.OrganizerProfilesUpdates]

	// Generate cache key
	cacheKey := fmt.Sprintf("organizer:update:%d:%d:%s", pagination.Page, pagination.Limit, pagination.Sort)

	// Try get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &organizerUpdate)
	}

	if organizerUpdate == nil {
		// If cache miss, get from db
		organizerUpdate, err = s.OrganizerUpdateRepo.FindUpdatesByProfileID(pagination, id)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(organizerUpdate); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return organizerUpdate, nil
}

func toOrganizerUpdateResponse(organizerUpdate *model.OrganizerProfilesUpdates) *dto.OrganizerUpdatesResponse {
	var npwp string
	var nib string

	if organizerUpdate.NPWPDocument != "" {
		npwp = fmt.Sprintf("http://localhost:8080/api/v1/staging-organizer/npwp/%s", organizerUpdate.NPWPDocument)
	}

	if organizerUpdate.NIBDocument != "" {
		nib = fmt.Sprintf("http://localhost:8080/api/v1/staging-organizer/nib/%s", organizerUpdate.NIBDocument)
	}

	return &dto.OrganizerUpdatesResponse{
		ID:           organizerUpdate.ID,
		ProfileID:    organizerUpdate.ProfileID,
		Status:       organizerUpdate.Status,
		Name:         organizerUpdate.Name,
		Slug:         organizerUpdate.Slug,
		Email:        helper.StringValue(organizerUpdate.Email),
		Instagram:    helper.StringValue(organizerUpdate.Instagram),
		Description:  helper.StringValue(organizerUpdate.Description),
		Address:      helper.StringValue(organizerUpdate.Address),
		PhoneNumber:  organizerUpdate.PhoneNumber,
		Country:      organizerUpdate.Country,
		NPWPNumber:   organizerUpdate.NPWPNumber,
		NPWPDocument: npwp,
		NIBNumber:    organizerUpdate.NIBNumber,
		NIBDocument:  nib,
		CreatedAt:    helper.ConvertDatetoUnix(organizerUpdate.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:    helper.ConvertDatetoUnix(organizerUpdate.UpdatedAt.Format(time.RFC3339)),
	}
}
