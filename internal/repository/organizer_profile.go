package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
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

// FindAllRaw returns a slice of OrganizerProfiles (not yet mapped to DTOs)
func (r *OrganizerRepository) FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) ([]*model.OrganizerProfiles, model.Pagination, error) {
	var profiles []*model.OrganizerProfiles

	query := r.db.Scopes(filterOrganizer(filter))

	if err := query.Preload("User").
		Scopes(Paginate(profiles, &pagination, query)).
		Find(&profiles).Error; err != nil {
		return nil, pagination, err
	}

	return profiles, pagination, nil
}

// FindAllForPublicRaw is the public-facing equivalent of FindAllRaw.
func (r *OrganizerRepository) FindAllForPublic(pagination model.Pagination, filter *dto.FilterPublicProfileReq,
) ([]*model.OrganizerProfiles, model.Pagination, error) {
	var profiles []*model.OrganizerProfiles

	query := r.db.Scopes(filterOrganizerForPublic(filter))

	if err := query.Preload("User").
		Scopes(Paginate(profiles, &pagination, query)).
		Find(&profiles).Error; err != nil {
		return nil, pagination, err
	}

	return profiles, pagination, nil
}

// CountEventsByProfileIDs returns a map of profileID
func (r *OrganizerRepository) CountEventsByProfileIDs(profileIDs []string) (map[string]int64, error) {
	if len(profileIDs) == 0 {
		return map[string]int64{}, nil
	}

	type row struct {
		ProfileID string `gorm:"column:profile_id"`
		Count     int64  `gorm:"column:count"`
	}

	var rows []row

	err := r.db.
		Model(&model.Events{}).
		Select("profile_id, COUNT(*) AS count").
		Where("profile_id IN ? AND status IN ('active', 'done')", profileIDs).
		Group("profile_id").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64, len(rows))
	for _, r := range rows {
		result[r.ProfileID] = r.Count
	}

	return result, nil
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

// FindBySlug implements [OrganizerProfileRepo].
func (r *OrganizerRepository) FindBySlug(slug string) (*model.OrganizerProfiles, error) {
	var profile *model.OrganizerProfiles

	if err := r.db.Preload("User").
		Where("slug = ?", slug).
		First(&profile).Error; err != nil {
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

func getProfileByCountry(country string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country = ?", country)
	}
}
