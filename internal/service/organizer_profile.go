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
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type OrganizerProfileService struct {
	OrganizerRepo       repository.OrganizerProfileRepo
	UserRepo            repository.UsersRepo
	OrganizerUpdateRepo repository.OrganizerProfileUpdateRepo
	EmailTaskPublisher  NewTaskEmail
	rdb                 *redis.Client
}

func NewOrganizerProfileService(
	organizerRepo repository.OrganizerProfileRepo,
	userRepo repository.UsersRepo,
	organizerUpdateRepo repository.OrganizerProfileUpdateRepo,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *OrganizerProfileService {
	return &OrganizerProfileService{
		OrganizerRepo:       organizerRepo,
		UserRepo:            userRepo,
		OrganizerUpdateRepo: organizerUpdateRepo,
		EmailTaskPublisher:  emailTaskPublisher,
		rdb:                 rdb,
	}
}

var (
	nibFilePath  = "./storage/nib"
	npwpFilePath = "./storage/npwp"
)

var organizerCache []string = []string{
	"organizer:all:*",
}

func (s *OrganizerProfileService) CreateProfile(profile *dto.CreateOrganizerProfileReq) error {
	validateReq := &dto.ValidateFileReq{
		Photo: profile.PhotoProfile,
		NPWP:  *profile.CompanyDetail.NPWPFile,
		NIB:   *profile.CompanyDetail.NIBFile,
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

	// Get corresponded user
	user, err := s.UserRepo.FindByID(newProfile.UserID)
	if err != nil {
		return err
	}

	// Get admins
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

	// Save npwp & nib file
	fileReq := &dto.SaveNPWPAndNIBFileReq{
		NPWP:     profile.CompanyDetail.NPWPFile,
		NIB:      profile.CompanyDetail.NIBFile,
		NPWPPath: npwpFilePath,
		NIBPath:  nibFilePath,
	}

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

	// Send email to organizer
	organizerpayload := &model.EmailPayload{
		To:   user.Email,
		Name: newProfile.Name,
	}

	if err := s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfile, organizerpayload); err != nil {
		log.Printf("[email] failed to send email: %v\n", err)
	}

	// Send email to admin
	for _, admin := range admins {
		adminPayload := &model.EmailPayload{
			To:        admin.Email,
			Name:      newProfile.Name,
			UserEmail: user.Email,
			Action:    "registered",
		}

		if err := s.EmailTaskPublisher.Enqueue(model.TypeEmailAdminVerification, adminPayload); err != nil {
			log.Printf("[email] failed to send email: %v\n", err)
		}
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, organizerCache)

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

func (s *OrganizerProfileService) FindAll(pagination model.Pagination) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var organizers *model.PaginationRow[*dto.OrganizerProfilesResponse]

	// Genereate cache key
	cacheKey := fmt.Sprintf("organizer:all:%d:%d:%s", pagination.Limit, pagination.Page, pagination.Sort)

	// Tru get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &organizers)
	}

	if organizers == nil {
		// if cache miss, get from db
		organizers, err = s.OrganizerRepo.FindAll(pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(organizers); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return organizers, nil
}

func (s *OrganizerProfileService) FindByCountry(
	country string,
	pagination model.Pagination,
) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var organizers *model.PaginationRow[*dto.OrganizerProfilesResponse]

	// Genereate cache key
	cacheKey := fmt.Sprintf("organizer:all:%s:%d:%d:%s", country, pagination.Limit, pagination.Page, pagination.Sort)

	// Tru get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &organizers)
	}

	if organizers == nil {
		// if cache miss, get from db
		organizers, err = s.OrganizerRepo.FindByCountry(country, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(organizers); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return organizers, nil
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

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, organizerCache)

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

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, organizerCache)

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

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, organizerCache)

	return 0, nil
}

func (s *OrganizerProfileService) UpdateProfile(userID string, req *dto.UpdateOrganizerProfileReq) (int, bool, error) {
	profile, err := s.OrganizerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, false, errors.New("profile not found")
	}

	// Authorization
	if userID != profile.UserID && profile.User.Role != helper.StrPointerIfNotEmpty(string(model.Admin)) {
		return fiber.StatusUnauthorized, false, errors.New("unauthorized action")
	}

	admins, err := s.UserRepo.FindByRole("admin")
	if err != nil {
		return fiber.StatusBadRequest, false, err
	}

	// Validate phone
	phonenumber, country, err := utils.ValidatePhoneCode(req.PhoneNumber, req.ISO)
	if err != nil {
		return fiber.StatusBadRequest, false, err
	}

	// Detect critical changes
	criticalChanged := false

	if profile.Name != req.Name ||
		profile.Country != country ||
		profile.CompanyDetail.NPWPNumber != req.CompanyDetail.NPWP ||
		profile.CompanyDetail.NIBNumber != req.CompanyDetail.NIB {
		criticalChanged = true
	}

	// =============================
	// HANDLE FILE UPLOAD (STAGING)
	// =============================
	var npwpFile, nibFile string

	if req.CompanyDetail.NPWPFile != nil && req.CompanyDetail.NIBFile != nil {
		fileReq := &dto.SaveNPWPAndNIBFileReq{
			NPWP:     req.CompanyDetail.NPWPFile,
			NIB:      req.CompanyDetail.NIBFile,
			NPWPPath: npwpStagePath,
			NIBPath:  nibStagePath,
		}

		npwpFile, nibFile, err = saveNPWPAndNIBFile(fileReq)
		if err != nil {
			return fiber.StatusBadRequest, false, err
		}

		criticalChanged = true
	}

	// Only update if you already approve
	if profile.Status.Status == string(model.UpdatePending) && criticalChanged {
		return fiber.StatusBadRequest, false, errors.New("Profile still under verification from admin")
	}

	// =============================
	// IF CRITICAL → SAVE TO STAGING
	// =============================
	if criticalChanged {

		profileUpdate := &model.OrganizerProfilesUpdates{
			ProfileID:    profile.ID,
			Name:         req.Name,
			PhoneNumber:  fmt.Sprintf("+%s", phonenumber),
			Status:       string(model.UpdatePending),
			Country:      country,
			NPWPNumber:   req.CompanyDetail.NPWP,
			NIBNumber:    req.CompanyDetail.NIB,
			NPWPDocument: npwpFile,
			NIBDocument:  nibFile,
		}

		if err := s.OrganizerUpdateRepo.Create(profileUpdate); err != nil {
			return fiber.StatusBadRequest, false, err
		}

		// Update status → pending
		profile.Status.Status = "pending"

		if err := s.OrganizerRepo.Update(profile); err != nil {
			return fiber.StatusBadRequest, false, err
		}

		// Send email async
		go func() {
			// Organizer email
			organizerPayload := &model.EmailPayload{
				To:   profile.User.Email,
				Name: profile.Name,
			}
			s.EmailTaskPublisher.Enqueue(model.TypeEmailOrganizerProfile, organizerPayload)

			// Admin email
			for _, admin := range admins {
				adminPayload := &model.EmailPayload{
					To:        admin.Email,
					Name:      profile.Name,
					UserEmail: profile.User.Email,
					Action:    "updated",
				}
				s.EmailTaskPublisher.Enqueue(model.TypeEmailAdminVerification, adminPayload)
			}
		}()

	} else {
		// =============================
		// NON-CRITICAL → DIRECT UPDATE
		// =============================
		profile.Name = req.Name
		profile.PhoneNumber = fmt.Sprintf("+%s", phonenumber)
		profile.Country = country
		profile.Address = req.Address
		profile.SocialMedias.Email = req.SocialMedia.Email
		profile.SocialMedias.Instagram = req.SocialMedia.Instagram
		profile.CompanyDetail.Description = req.CompanyDetail.Description
		profile.UpdatedAt = time.Now().UTC()

		if err := s.OrganizerRepo.Update(profile); err != nil {
			return fiber.StatusBadRequest, false, err
		}
	}

	// Invalidate cache
	utils.InvalidateCache(s.rdb, organizerCache)

	return fiber.StatusOK, criticalChanged, nil
}

func validateFile(req *dto.ValidateFileReq) error {
	// Validate image
	if err := helper.ValidateImage(req.Photo); err != nil {
		return err
	}

	// Validate npwp file
	if err := helper.ValidatePDF(&req.NPWP); err != nil {
		return err
	}

	// Validate nib file
	if err := helper.ValidatePDF(&req.NIB); err != nil {
		return err
	}

	return nil
}

func saveNPWPAndNIBFile(req *dto.SaveNPWPAndNIBFileReq) (string, string, error) {
	npwpPath, npwpFile, err := helper.SaveToLocal(req.NPWP, req.NPWPPath)
	if err != nil {
		return "", "", err
	}

	nipPath, nibFile, err := helper.SaveToLocal(req.NIB, req.NIBPath)
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
	var reviewedAt int64

	if profile.Status.ReviewedAt != nil {
		reviewedAt = helper.ConvertDatetoUnix(profile.Status.ReviewedAt.Format(time.RFC3339))
	}

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
		PhotoProfile: fmt.Sprintf("http://localhost:8080/api/v1/organizer/photo/%s", helper.StringValue(profile.PhotoProfile)),
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
			NPWPFile:    fmt.Sprintf("http://localhost:8080/api/v1/organizer/npwp/%s", profile.CompanyDetail.NPWPDocument),
			NIB:         profile.CompanyDetail.NIBNumber,
			NIBFile:     fmt.Sprintf("http://localhost:8080/api/v1/organizer/nib/%s", profile.CompanyDetail.NIBDocument),
		},
		CreatedAt: profile.CreatedAt.Unix(),
		UpdatedAt: profile.UpdatedAt.Unix(),
	}
}
