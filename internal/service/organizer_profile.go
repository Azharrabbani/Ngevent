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
	EventsRepo          repository.EventsRepo
	EventsUpdateRepo    repository.EventsUpdateRepo
	EmailTaskPublisher  NewTaskEmail
	rdb                 *redis.Client
}

func NewOrganizerProfileService(
	organizerRepo repository.OrganizerProfileRepo,
	userRepo repository.UsersRepo,
	organizerUpdateRepo repository.OrganizerProfileUpdateRepo,
	emailTaskPublisher NewTaskEmail,
	eventsRepo repository.EventsRepo,
	eventsUpdateRepo repository.EventsUpdateRepo,
	rdb *redis.Client,
) *OrganizerProfileService {
	return &OrganizerProfileService{
		OrganizerRepo:       organizerRepo,
		UserRepo:            userRepo,
		OrganizerUpdateRepo: organizerUpdateRepo,
		EventsRepo:          eventsRepo,
		EventsUpdateRepo:    eventsUpdateRepo,
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
	"organizer:public:*",
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
		Slug:        utils.GenerateEventSlug(profile.Name),
		Address:     &profile.Address,
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

func (s *OrganizerProfileService) FindBySlug(slug string) (*dto.OrganizerProfilesResponse, error) {
	profile, err := s.OrganizerRepo.FindBySlug(slug)
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

func (s *OrganizerProfileService) FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var organizers *model.PaginationRow[*dto.OrganizerProfilesResponse]

	// Genereate cache key
	cacheKey := fmt.Sprintf(
		"organizer:all:%d:%d:%s:%s:%s:%s",
		pagination.Limit,
		pagination.Page,
		pagination.Sort,
		helper.StringValue(filter.Filter),
		helper.StringValue(filter.Status),
		helper.BoolValue(filter.RequestUpdates),
	)

	// Tru get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &organizers)
	}

	if organizers == nil {
		// if cache miss, get from db
		profiles, paginationMeta, err := s.OrganizerRepo.FindAll(pagination, filter)
		if err != nil {
			return nil, err
		}

		counts := s.attachEventCounts(profiles)
		rows := toOrganizerListResponse(profiles, counts)

		organizers = &model.PaginationRow[*dto.OrganizerProfilesResponse]{
			Pagination: paginationMeta,
			Rows:       rows,
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(organizers); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return organizers, nil
}

func (s *OrganizerProfileService) FindAllForPublic(pagination model.Pagination, filter *dto.FilterPublicProfileReq) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var organizers *model.PaginationRow[*dto.OrganizerProfilesResponse]

	filterValue := ""

	if filter != nil && filter.Filter != nil {
		filterValue = strings.TrimSpace(
			strings.ToLower(*filter.Filter),
		)
	}

	// Genereate cache key
	cacheKey := fmt.Sprintf("organizer:public:%d:%d:%s:%s", pagination.Limit, pagination.Page, pagination.Sort, filterValue)

	// Tru get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &organizers)
	}

	if organizers == nil {
		// if cache miss, get from db
		profiles, paginationMeta, err := s.OrganizerRepo.FindAllForPublic(pagination, filter)
		if err != nil {
			return nil, err
		}

		counts := s.attachEventCounts(profiles)
		rows := toOrganizerListResponse(profiles, counts)

		organizers = &model.PaginationRow[*dto.OrganizerProfilesResponse]{
			Pagination: paginationMeta,
			Rows:       rows,
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(organizers); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return organizers, nil
}

func (s *OrganizerProfileService) attachEventCounts(profiles []*model.OrganizerProfiles) map[string]int64 {
	ids := make([]string, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.ID)
	}

	counts, err := s.OrganizerRepo.CountEventsByProfileIDs(ids)
	if err != nil {
		return map[string]int64{}
	}

	return counts
}

func (s *OrganizerProfileService) CloseAccount(userID string) (int, error) {
	// 1. Find organizer profile by userID
	profile, err := s.OrganizerRepo.FindByUserID(userID)
	if err != nil {
		return fiber.StatusNotFound, errors.New("profile not found")
	}

	// 2. Check for blocking events (pending or active)
	hasBlocking, err := s.EventsRepo.HasBlockingEvents(profile.ID)
	if err != nil {
		return fiber.StatusInternalServerError, errors.New("failed to check events status")
	}

	if hasBlocking {
		return fiber.StatusConflict, errors.New("account cannot be closed while you have active or pending events, please cancel or wait for them to complete first")
	}

	// 3. Begin shared transaction
	tx := s.OrganizerRepo.GetDB().Begin()
	if tx.Error != nil {
		return fiber.StatusInternalServerError, errors.New("failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[PANIC] CloseAccount transaction rolled back: %v", r)
		}
	}()

	// 4. Soft delete event update categories + event updates
	if err := s.EventsUpdateRepo.SoftDeleteEventUpdates(tx, profile.ID); err != nil {
		tx.Rollback()
		return fiber.StatusInternalServerError, errors.New("failed to remove event updates")
	}

	// 5. Soft delete event categories + events
	if err := s.EventsRepo.SoftDeleteEvents(tx, profile.ID); err != nil {
		tx.Rollback()
		return fiber.StatusInternalServerError, errors.New("failed to remove events")
	}

	// 6. Soft delete organizer profile updates
	if err := s.OrganizerUpdateRepo.SoftDeleteProfileUpdates(tx, profile.ID); err != nil {
		tx.Rollback()
		return fiber.StatusInternalServerError, errors.New("failed to remove profile update requests")
	}

	// 7. Soft delete organizer profile (status → deactivated)
	if err := s.OrganizerRepo.SoftDeleteProfile(tx, profile.ID); err != nil {
		tx.Rollback()
		return fiber.StatusInternalServerError, errors.New("failed to deactivate organizer profile")
	}

	// 8. Soft delete user
	if err := s.UserRepo.SoftDeleteUser(tx, userID); err != nil {
		tx.Rollback()
		return fiber.StatusInternalServerError, errors.New("failed to close user account")
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.StatusInternalServerError, errors.New("failed to commit account closure")
	}

	utils.InvalidateCache(s.rdb, organizerCache)
	utils.InvalidateCache(s.rdb, organizerUpdateCache)
	utils.InvalidateCache(s.rdb, eventCache)
	utils.InvalidateCache(s.rdb, updatedEventCache)
	utils.InvalidateCache(s.rdb, userCache)

	return fiber.StatusOK, nil
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

	if profile.Status.Status == string(model.Pending) || profile.RequestUpdates {
		return fiber.StatusBadRequest, errors.New("Profile still under review. Please wait for the review process to complete before updating.")
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

	if profile.Status.Status == string(model.Pending) {
		return fiber.StatusBadRequest, false, errors.New("Profile still under review. Please wait for the review process to complete before updating.")
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
		profile.PhoneNumber != req.PhoneNumber ||
		profile.Country != country ||
		profile.CompanyDetail.NPWPNumber != req.CompanyDetail.NPWP ||
		profile.CompanyDetail.NIBNumber != req.CompanyDetail.NIB {
		criticalChanged = true
	}

	// =============================
	// HANDLE FILE UPLOAD (STAGING)
	// =============================
	var npwpFile, nibFile string

	hasNPWP := req.CompanyDetail.NPWPFile != nil
	hasNIB := req.CompanyDetail.NIBFile != nil

	if hasNPWP || hasNIB {
		if hasNPWP {
			npwpPath, npwpFileName, err := helper.SaveToLocal(req.CompanyDetail.NPWPFile, npwpStagePath)
			if err != nil {
				return fiber.StatusBadRequest, false, err
			}
			extNPWP := filepath.Ext(npwpPath)
			baseNPWP := strings.TrimSuffix(npwpPath, extNPWP)
			tempNPWP := baseNPWP + "_tmp" + extNPWP
			if err := helper.CompressPDF(npwpPath, tempNPWP); err != nil {
				return fiber.StatusBadRequest, false, err
			}
			if err := helper.OptimizePDF(tempNPWP); err != nil {
				return fiber.StatusBadRequest, false, err
			}
			if err := os.Rename(tempNPWP, npwpPath); err != nil {
				return fiber.StatusBadRequest, false, err
			}
			npwpFile = npwpFileName
		}

		if hasNIB {
			nibPath, nibFileName, err := helper.SaveToLocal(req.CompanyDetail.NIBFile, nibStagePath)
			if err != nil {
				return fiber.StatusBadRequest, false, err
			}
			extNIB := filepath.Ext(nibPath)
			baseNIB := strings.TrimSuffix(nibPath, extNIB)
			tempNIB := baseNIB + "_tmp" + extNIB
			if err := helper.CompressPDF(nibPath, tempNIB); err != nil {
				return fiber.StatusBadRequest, false, err
			}
			if err := helper.OptimizePDF(tempNIB); err != nil {
				return fiber.StatusBadRequest, false, err
			}
			if err := os.Rename(tempNIB, nibPath); err != nil {
				return fiber.StatusBadRequest, false, err
			}
			nibFile = nibFileName
		}

		criticalChanged = true
	}

	// =============================
	// IF CRITICAL → SAVE TO STAGING
	// =============================
	if criticalChanged {

		profileUpdate := &model.OrganizerProfilesUpdates{
			ProfileID:    profile.ID,
			Name:         req.Name,
			Slug:         utils.GenerateEventSlug(req.Name),
			PhoneNumber:  fmt.Sprintf("+%s", phonenumber),
			Status:       "pending",
			Country:      country,
			Email:        req.SocialMedia.Email,
			Instagram:    req.SocialMedia.Instagram,
			Address:      req.Address,
			Description:  req.CompanyDetail.Description,
			NPWPNumber:   req.CompanyDetail.NPWP,
			NIBNumber:    req.CompanyDetail.NIB,
			NPWPDocument: npwpFile,
			NIBDocument:  nibFile,
		}

		if err := s.OrganizerUpdateRepo.Create(profileUpdate); err != nil {
			return fiber.StatusBadRequest, false, err
		}

		profile.RequestUpdates = true
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
	utils.InvalidateCache(s.rdb, organizerUpdateCache)

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
	if req.NIB == nil || req.NPWP == nil {
		return "", "", errors.New("NPWP and NIB file required")
	}

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
			ReviewedBy:     &profile.User.Email,
			ReviewedAt:     &reviewedAt,
		},
		RequestUpdates: profile.RequestUpdates,
		Email:          profile.User.Email,
		Name:           profile.Name,
		Slug:           profile.Slug,
		PhotoProfile:   fmt.Sprintf("http://localhost:8080/api/v1/organizer/photo/%s", helper.StringValue(profile.PhotoProfile)),
		PhoneNumber:    profile.PhoneNumber,
		Country:        profile.Country,
		Address:        profile.Address,
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

func toOrganizerListResponse(profiles []*model.OrganizerProfiles, counts map[string]int64) []*dto.OrganizerProfilesResponse {
	out := make([]*dto.OrganizerProfilesResponse, 0, len(profiles))

	for _, profile := range profiles {
		var reviewedAt int64
		if profile.Status.ReviewedAt != nil {
			reviewedAt = helper.ConvertDatetoUnix(profile.Status.ReviewedAt.Format(time.RFC3339))
		}

		createdAt := helper.ConvertDatetoUnix(profile.CreatedAt.Format(time.RFC3339))
		updatedAt := helper.ConvertDatetoUnix(profile.UpdatedAt.Format(time.RFC3339))

		count := counts[profile.ID]

		out = append(out, &dto.OrganizerProfilesResponse{
			ID:     profile.ID,
			UserID: profile.UserID,
			Status: dto.OrganizerStatusResp{
				Status:         profile.Status.Status,
				RejectedReason: profile.Status.RejectedReason,
				ReviewedBy:     profile.Status.ReviewedBy,
				ReviewedAt:     &reviewedAt,
			},
			RequestUpdates: profile.RequestUpdates,
			Name:           profile.Name,
			Slug:           profile.Slug,
			Email:          profile.User.Email,
			PhotoProfile:   fmt.Sprintf("http://localhost:8080/api/v1/organizer/photo/%s", helper.StringValue(profile.PhotoProfile)),
			PhoneNumber:    profile.PhoneNumber,
			Country:        profile.Country,
			Address:        profile.Address,
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
			EventCount: &count,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	return out
}
