package service

import (
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
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

var (
	nibFilePath  = "./storage/nib"
	npwpFilePath = "./storage/npwp"
)

func (s *OrganizerProfileService) CreateProfile(profile *dto.CreateOrganizerProfileReq) error {
	validateReq := &dto.ValidateFileReq{
		Photo: profile.PhotoProfile,
		NPWP:  profile.CompanyDetail.NPWPFile,
		NIB:   profile.CompanyDetail.NIBFile,
	}

	// Validate files
	if err := validateFile(validateReq); err != nil {
		return err
	}

	// Validate phone code
	phonenumber, country, err := utils.ValidatePhoneCode(profile.PhoneNumber, profile.ISO)
	if err != nil {
		return err
	}

	newProfile := &model.OrganizerProfiles{
		UserID:      profile.UserID,
		Name:        profile.Name,
		Address:     profile.Address,
		PhoneNumber: fmt.Sprintf("+%s", phonenumber),
		Country:     country,
		SocialMedias: model.OrganizerSocialMedia{
			Email:     profile.SocialMedia.Email,
			Instagram: profile.SocialMedia.Instagram,
		},
		CompanyDetail: model.OrganizerCompDetail{
			Description: profile.CompanyDetail.Description,
			NPWPNumber:  profile.CompanyDetail.NPWP,
			NIBNumber:   profile.CompanyDetail.NIB,
		},
	}

	admins, err := s.UserRepo.FindByRole("admin")
	if err != nil {
		return err
	}

	if profile.PhotoProfile != nil {
		// Save photo profile to local storage
		_, fileName, err := helper.SaveToLocal(profile.PhotoProfile, profileUploadPath)
		if err != nil {
			return err
		}

		newProfile.PhotoProfile = &fileName
	}

	fileReq := &dto.SaveNPWPAndNIBFileReq{
		NPWP: &profile.CompanyDetail.NPWPFile,
		NIB:  &profile.CompanyDetail.NIBFile,
	}
	
	// Save npwp & nib file
	npwpFile, nibFile, err := saveNPWPAndNIBFile(fileReq)
	if err != nil {
		return err
	}

	newProfile.CompanyDetail.NPWPDocument = npwpFile
	newProfile.CompanyDetail.NIBDocument = nibFile

	// Save Profile
	if err := s.OrganizerRepo.Create(newProfile); err != nil {
		return err
	}

	organizerpayload := &model.EmailPayload{
		To:   newProfile.User.Email,
		Name: newProfile.Name,
	}
	
	// Send email to organizer
	s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfile, organizerpayload)

	for _, admin := range admins {
		adminPayload := &model.EmailPayload{
			To:        admin.Email,
			Name:      newProfile.Name,
			UserEmail: newProfile.User.Email,
			Action:    "registered",
		}
		
		// Send email to admin
		s.EmailTaskPublisher.Enqueue(model.TypeEmailAdminVerification, adminPayload)
	}

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

func (s *OrganizerProfileService) VerifiedProfile(id string, req *dto.ApprovedReq) error {
	// Check user is exist
	profile, err := s.OrganizerRepo.FindByID(id)
	if err != nil {
		return errors.New("profile not found")
	}

	if profile.Status.Status == "approved" {
		return errors.New("profile already approved")
	}

	if err := s.OrganizerRepo.VerifiedProfile(id, req); err != nil {
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

func (s *OrganizerProfileService) RejectProfile(id string, req *dto.RejectedReq) error {
	// Check user is exist
	profile, err := s.OrganizerRepo.FindByID(id)
	if err != nil {
		return errors.New("profile not found")
	}

	if profile.Status.Status == "rejected" {
		return errors.New("profile already rejected")
	}

	if err := s.OrganizerRepo.RejectProfile(id, req); err != nil {
		return errors.New("failed to reject profile")
	}

	// Send email
	payload := &model.RejectedEmailPayload{
		To:     profile.User.Email,
		Name:   profile.Name,
		Reason: req.Reason,
	}
	s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfileRejected, payload)

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

	profile, err := s.OrganizerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	admins, err := s.UserRepo.FindByRole("admin")
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	if userID != profile.UserID {
		return fiber.StatusUnauthorized, errors.New("unauthorized action")
	}

	// Validate phone number
	phonenumber, country, err := utils.ValidatePhoneCode(req.PhoneNumber, req.ISO)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Track whether critical field changed
	criticalChanged := false

	// Only check critical fields if currently approved
	if profile.Status.Status == "approved" {

		if profile.Name != req.Name {
			criticalChanged = true
		}

		if profile.Country != country {
			criticalChanged = true
		}

		if criticalChanged {
			profile.Status.Status = "pending"
		}
	}

	// Update allowed fields
	profile.Name = req.Name
	profile.PhoneNumber = fmt.Sprintf("+%s", phonenumber)
	profile.Country = country
	profile.Address = req.Address
	profile.SocialMedias.Email = req.SocialMedia.Email
	profile.SocialMedias.Instagram = req.SocialMedia.Instagram
	profile.CompanyDetail.Description = req.CompanyDetail.Description

	// Only allow NPWP & NIB change if not approved
	if profile.Status.Status != "approved" {
		profile.CompanyDetail.NPWPNumber = req.CompanyDetail.NPWP
		profile.CompanyDetail.NIBNumber = req.CompanyDetail.NIB

		oldNPWP := fmt.Sprintf("%s/%s", npwpFilePath, profile.CompanyDetail.NPWPDocument)
		oldNIB := fmt.Sprintf("%s/%s", nibFilePath, profile.CompanyDetail.NIBDocument)

		// Save new NPWP & NIB file
		fileReq := &dto.SaveNPWPAndNIBFileReq{
			NPWP: &req.CompanyDetail.NPWPFile,
			NIB:  &req.CompanyDetail.NIBFile,
		}

		npwpFile, nibFile, err := saveNPWPAndNIBFile(fileReq)
		if err != nil {
			return fiber.StatusBadRequest, err
		}

		profile.CompanyDetail.NPWPDocument = npwpFile
		profile.CompanyDetail.NIBDocument = nibFile

		// Remove old files
		if err := os.Remove(oldNPWP); err != nil {
			fmt.Errorf("failed to remove file: %v", err)
		}

		if err := os.Remove(oldNIB); err != nil {
			fmt.Errorf("failed to remove file: %v", err)
		}
	}

	profile.UpdatedAt = time.Now().UTC()

	if err := s.OrganizerRepo.Update(profile); err != nil {
		return fiber.StatusBadRequest, err
	}

	if criticalChanged {
		// Send email to organizer
		organizerpayload := &model.EmailPayload{
			To:   profile.User.Email,
			Name: profile.Name,
		}

		s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfile, organizerpayload)

		// Send email to admin
		for _, admin := range admins {
			adminPayload := &model.EmailPayload{
				To:        admin.Email,
				Name:      profile.Name,
				UserEmail: profile.User.Email,
				Action:    "updated",
			}

			s.EmailTaskPublisher.Enqueue(model.TypeEmailAdminVerification, adminPayload)
		}
	}

	return fiber.StatusOK, nil
}

func validateFile(req *dto.ValidateFileReq) error {
	// Validate image
	if err := helper.ValidateImage(req.Photo); err != nil {
		return err
	}

	// Validate npwp file
	if err := helper.ValidatePDF(req.NPWP); err != nil {
		return err
	}

	// Validate nib file
	if err := helper.ValidatePDF(req.NIB); err != nil {
		return err
	}

	return nil
}

func saveNPWPAndNIBFile(req *dto.SaveNPWPAndNIBFileReq) (string, string, error) {
	npwpPath, npwpFile, err := helper.SaveToLocal(req.NPWP, npwpFilePath)
	if err != nil {
		return "", "", err
	}

	nipPath, nibFile, err := helper.SaveToLocal(req.NIB, nibFilePath)
	if err != nil {
		return "", "", err
	}

	// temp compressed file
	extNPWP := filepath.Ext(npwpPath)
	baseNPWP := strings.TrimSuffix(npwpPath, extNPWP)
	tempPathNPWP := baseNPWP + "_tmp" + extNPWP

	extNIB := filepath.Ext(nipPath)
	baseNIB := strings.TrimSuffix(nipPath, extNIB)
	tempPathNIB := baseNIB + "_tmp" + extNIB

	// Compress the PDF file
	if err := helper.CompressPDF(npwpPath, tempPathNPWP); err != nil {
		return "", "", err
	}

	if err := helper.CompressPDF(nipPath, tempPathNIB); err != nil {
		return "", "", err
	}

	// Optimize the file
	if err := helper.OptimizePDF(tempPathNPWP); err != nil {
		return "", "", err
	}

	if err := helper.OptimizePDF(tempPathNIB); err != nil {
		return "", "", err
	}

	// replace original
	if err := os.Rename(tempPathNPWP, npwpPath); err != nil {
		return "", "", err
	}

	if err := os.Rename(tempPathNIB, nipPath); err != nil {
		return "", "", err
	}

	// return final path
	return npwpFile, nibFile, nil
}

func toOrganizerProfileResponse(profile *model.OrganizerProfiles) *dto.OrganizerProfilesResponse {
	reviewedAt := helper.ConvertDatetoUnix(profile.Status.ReviewedAt.Format(time.RFC3339))
	return &dto.OrganizerProfilesResponse{
		ID:     profile.ID,
		UserID: profile.UserID,
		Status: dto.OrganizerStatusResp{
			Status:         profile.Status.Status,
			RejectedReason: profile.Status.RejectedReason,
			ReviewedBy:     profile.Status.ReviewedBy,
			ReviewedAt:     &reviewedAt,
		},
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
		CompanyDetail: dto.OrganizerCompDetailRes{
			Description: profile.CompanyDetail.Description,
			NPWP:        profile.CompanyDetail.NPWPNumber,
			NPWPFile:    profile.CompanyDetail.NPWPDocument,
			NIB:         profile.CompanyDetail.NIBNumber,
			NIBFile:     profile.CompanyDetail.NIBDocument,
		},
	}
}
