package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"time"

	"gorm.io/gorm"
)

type OrganizerRepository struct {
	db *gorm.DB
}

func NewOrganizerRepository(db *gorm.DB) OrganizerProfileRepo {
	return &OrganizerRepository{db: db}
}

// Create implements OrganizerProfileRepo.
func (r *OrganizerRepository) Create(profile *model.OrganizerProfiles) error {
	return r.db.Create(profile).Error
}

// Delete implements OrganizerProfileRepo.
func (r *OrganizerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.OrganizerProfiles{}).Error
}

// FindByCountry implements OrganizerProfileRepo.
func (r *OrganizerRepository) FindByCountry(country string, pagination model.Pagination) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var profiles []*model.OrganizerProfiles

	query := r.db.Scopes(getProfileByCountry(country))

	if err := query.Preload("User").
		Scopes(Paginate(profiles, &pagination, query)).
		Find(&profiles).Error; err != nil {
		return nil, err
	}

	// Transform data to response struct
	organizers := toOrganizerResponse(profiles)

	return &model.PaginationRow[*dto.OrganizerProfilesResponse]{
		Pagination: pagination,
		Rows:       organizers,
	}, nil
}

// FindByID implements OrganizerProfileRepo.
func (r *OrganizerRepository) FindByID(id string) (*model.OrganizerProfiles, error) {
	var profile *model.OrganizerProfiles

	if err := r.db.Preload("User").
		Where("id = ?", id).First(&profile).
		Error; err != nil {
		return nil, err
	}

	return profile, nil
}

// FindByUserID implements OrganizerProfileRepo.
func (r *OrganizerRepository) FindByUserID(userID string) (*model.OrganizerProfiles, error) {
	var profile *model.OrganizerProfiles

	if err := r.db.Preload("User").
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

// VerifiedProfile implements OrganizerProfileRepo.
func (r *OrganizerRepository) VerifiedProfile(id string, req *dto.ApprovedReq) error {
	return r.db.
		Model(&model.OrganizerProfiles{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "approved",
			"reviewed_by": req.ReviewedBy,
			"reviewed_at": req.ReviewedAt,
		}).Error
}

// RejectProfile implements OrganizerProfileRepo.
func (r *OrganizerRepository) RejectProfile(id string, req *dto.RejectedReq) error {
	return r.db.
		Model(&model.OrganizerProfiles{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":          "rejected",
			"rejected_reason": req.Reason,
			"reviewed_by":     req.ReviewedBy,
			"reviewed_at":     req.ReviewedAt,
		}).Error
}

// Update implements OrganizerProfileRepo.
func (r *OrganizerRepository) Update(profile *model.OrganizerProfiles) error {
	return r.db.Save(profile).Error
}

// UpdatePhotoProfile implements OrganizerProfileRepo.
func (r *OrganizerRepository) UpdatePhotoProfile(userID string, photo string) error {
	return r.db.Model(&model.AttendeeProfiles{}).
		Where("user_id = ?", userID).
		Updates(&model.OrganizerProfiles{
			PhotoProfile: &photo,
			UpdatedAt:    time.Now().UTC()}).Error
}

func toOrganizerResponse(profiles []*model.OrganizerProfiles) []*dto.OrganizerProfilesResponse {
	var organizers []*dto.OrganizerProfilesResponse

	for _, profile := range profiles {
		reviewedAt := helper.ConvertDatetoUnix(profile.Status.ReviewedAt.Format(time.RFC3339))
		organizers = append(organizers, &dto.OrganizerProfilesResponse{
			ID: profile.ID,
			Status: dto.OrganizerStatusResp{
				Status:         profile.Status.Status,
				RejectedReason: profile.Status.RejectedReason,
				ReviewedBy:     profile.Status.ReviewedBy,
				ReviewedAt:     &reviewedAt,
			},
			Name:         profile.Name,
			Email:        profile.User.Email,
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
				NIB:         profile.CompanyDetail.NIBDocument,
			},
		})
	}

	return organizers
}

func getProfileByCountry(country string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country = ?", country)
	}
}
