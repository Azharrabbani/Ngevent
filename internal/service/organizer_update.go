package service

import (
	"errors"
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"
	"os"
	"time"
)

type OrganizerUpdateService struct {
	UserRepo             repository.UsersRepo
	OrganizerProfileRepo repository.OrganizerProfileRepo
	OrganizerUpdateRepo  repository.OrganizerProfileUpdateRepo
	EmailTaskPublisher   NewTaskEmail
}

func NewOrganizerUpdateService(
	userRepo repository.UsersRepo,
	organizerProfileRepo repository.OrganizerProfileRepo,
	organizerUpdateRepo repository.OrganizerProfileUpdateRepo,
	emailTaskPublisher NewTaskEmail,
) *OrganizerUpdateService {
	return &OrganizerUpdateService{
		UserRepo:             userRepo,
		OrganizerProfileRepo: organizerProfileRepo,
		OrganizerUpdateRepo:  organizerUpdateRepo,
		EmailTaskPublisher:   emailTaskPublisher,
	}
}

var (
	npwpStagePath = "./storage/npwp/stage"
	nibStagePath  = "./storage/nib/stage"
)

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
	profile, err := s.OrganizerProfileRepo.FindByUserID(updateData.ProfileID)
	if err != nil {
		return errors.New("profile not found")
	}

	// Determine the status
	switch status := req.Status; status {
	case "approved":
		// 1. When approves update the profile with the update data in staging
		profile.Name = updateData.Name
		profile.PhoneNumber = updateData.PhoneNumber
		profile.Country = updateData.Country
		profile.CompanyDetail.NPWPNumber = updateData.NPWPNumber
		profile.CompanyDetail.NIBNumber = updateData.NIBNumber

		// 2. Copy file from staging to destination path
		npwpSrcPath := fmt.Sprintf("%s/%s", npwpStagePath, updateData.NPWPDocument)
		nibSrcPath := fmt.Sprintf("%s/%s", nibStagePath, updateData.NIBDocument)

		npwpFile, err := helper.CopyFile(npwpSrcPath, npwpFilePath)
		if err != nil {
			return err
		}

		nibFile, err := helper.CopyFile(nibSrcPath, nibFilePath)
		if err != nil {
			return err
		}

		profile.CompanyDetail.NPWPDocument = npwpFile
		profile.CompanyDetail.NIBDocument = nibFile

		// 3. Update the data
		if err := profileX.Save(profile).Error; err != nil {
			profileX.Rollback()
			return err
		}

		// 4. Update staging status
		status := "approved"
		if err := s.OrganizerUpdateRepo.Validate(updateData.ID, status); err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return errors.New("validate failed")
		}

		// 5. Remove old files
		oldNPWP := fmt.Sprintf("%s/%s", npwpFilePath, profile.CompanyDetail.NPWPDocument)
		oldNIB := fmt.Sprintf("%s/%s", nibFilePath, profile.CompanyDetail.NIBDocument)

		if err := os.Remove(oldNPWP); err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return err
		}

		if err := os.Remove(oldNIB); err != nil {
			profileX.Rollback()
			updateX.Rollback()
			return err
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

	return nil
}

func (s *OrganizerUpdateService) FindByID(id string) (*model.OrganizerProfilesUpdates, error) {
	return s.OrganizerUpdateRepo.FindByID(id)
}

func (s *OrganizerUpdateService) FindByProfileID(id string, pagination model.Pagination) (*model.PaginationRow[*model.OrganizerProfilesUpdates], error) {
	return s.OrganizerUpdateRepo.FindByProfileID(pagination, id)
}
