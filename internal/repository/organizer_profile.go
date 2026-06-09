package repository

import (
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrganizerRepository struct {
	db *gorm.DB
}

func NewOrganizerRepository(db *gorm.DB) OrganizerProfileRepo {
	return &OrganizerRepository{db: db}
}

// GetDB implements OrganizerProfileRepo.
func (r *OrganizerRepository) GetDB() *gorm.DB {
	return r.db
}

// Create implements OrganizerProfileRepo.
func (r *OrganizerRepository) Create(profile *model.OrganizerProfiles) error {
	return r.db.Create(profile).Error
}

// HasProfile implements OrganizerProfileRepo.
func (r *OrganizerRepository) HasProfile(userID string) (bool, error) {
	if err := r.db.
		Where("user_id = ?", userID).
		First(&model.OrganizerProfiles{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

// Delete implements OrganizerProfileRepo.
func (r *OrganizerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.OrganizerProfiles{}).Error
}

// FindAll implements OrganizerProfileRepo.
func (r *OrganizerRepository) FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var profiles []*model.OrganizerProfiles

	query := r.db.Scopes(filterOrganizer(filter))

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

// FindAllForPublic implements [OrganizerProfileRepo].
func (r *OrganizerRepository) FindAllForPublic(pagination model.Pagination, filter *dto.FilterPublicProfileReq) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error) {
	var profiles []*model.OrganizerProfiles

	query := r.db.Scopes(filterOrganizerForPublic(filter))

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
			"status":          "approved",
			"rejected_reason": nil,
			"reviewed_by":     req.ReviewedBy,
			"reviewed_at":     req.ReviewedAt,
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
	return r.db.Model(&model.OrganizerProfiles{}).
		Where("user_id = ?", userID).
		Updates(&model.OrganizerProfiles{
			PhotoProfile: &photo,
			UpdatedAt:    time.Now().UTC()}).Error
}

func (r *OrganizerRepository) SoftDeleteProfile(tx *gorm.DB, profileID string) error {
	now := time.Now().UTC()

	return tx.Model(&model.OrganizerProfiles{}).
		Where("id = ?", profileID).
		Updates(map[string]interface{}{
			"status":     "deactivated",
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func filterOrganizerForPublic(filter *dto.FilterPublicProfileReq) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(
			"organizer_profiles.deleted_at IS NULL",
		).Where(
			"organizer_profiles.status = ?",
			model.Approved,
		)

		if filter != nil && filter.Filter != nil {
			query := "%" + strings.ToLower(*filter.Filter) + "%"

			db = db.Joins(
				"JOIN users ON users.id = organizer_profiles.user_id",
			).Where(`
                LOWER(users.email) LIKE ?
                OR LOWER(organizer_profiles.name) LIKE ?
                OR LOWER(organizer_profiles.country) LIKE ?
            `, query, query, query)
		}

		return db
	}
}

func filterOrganizer(filter *dto.FilterProfileReq) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if filter.Status != nil && *filter.Status == "deactivated" {
			db = db.Where("organizer_profiles.deleted_at IS NOT NULL")
		} else {
			db = db.Where("organizer_profiles.deleted_at IS NULL")
		}

		if filter.Filter != nil {
			query := "%" + strings.ToLower(*filter.Filter) + "%"

			db = db.Joins("JOIN users ON users.id = organizer_profiles.user_id").
				Where(
					db.Where("LOWER(users.email) LIKE ?", query).
						Or("LOWER(organizer_profiles.name) LIKE ?", query).
						Or("LOWER(organizer_profiles.country) LIKE ?", query),
				)
		}

		if filter.Status != nil && *filter.Status != "deactivated" {
			db = db.Where("organizer_profiles.status = ?", *filter.Status)
		}

		return db
	}
}

func toOrganizerResponse(profiles []*model.OrganizerProfiles) []*dto.OrganizerProfilesResponse {
	var organizers []*dto.OrganizerProfilesResponse

	for _, profile := range profiles {
		var reviewedAt int64

		if profile.Status.ReviewedAt != nil {
			reviewedAt = helper.ConvertDatetoUnix(profile.Status.ReviewedAt.Format(time.RFC3339))
		}

		createdAt := helper.ConvertDatetoUnix(profile.CreatedAt.Format(time.RFC3339))
		updatedAt := helper.ConvertDatetoUnix(profile.UpdatedAt.Format(time.RFC3339))

		organizers = append(organizers, &dto.OrganizerProfilesResponse{
			ID:     profile.ID,
			UserID: profile.UserID,
			Status: dto.OrganizerStatusResp{
				Status:         profile.Status.Status,
				RejectedReason: profile.Status.RejectedReason,
				ReviewedBy:     profile.Status.ReviewedBy,
				ReviewedAt:     &reviewedAt,
			},
			Name:         profile.Name,
			Email:        profile.User.Email,
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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return organizers
}

func getProfileByCountry(country string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country = ?", country)
	}
}
