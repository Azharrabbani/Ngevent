package repository

import (
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

type OrganizerProfileUpdateRepository struct {
	db *gorm.DB
}

func NewOrganizerProfileUpdateRepository(db *gorm.DB) OrganizerProfileUpdateRepo {
	return &OrganizerProfileUpdateRepository{db: db}
}

// GetDB implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) GetDB() *gorm.DB {
	return r.db
}

// Create implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) Create(profile *model.OrganizerProfilesUpdates) error {
	return r.db.Create(profile).Error
}

// Delete implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) Delete(id string) error {
	return r.db.Delete(&model.OrganizerProfilesUpdates{}).Error
}

// FindByID implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) FindByID(id string) (*model.OrganizerProfilesUpdates, error) {
	var profileUpdate *model.OrganizerProfilesUpdates

	if err := r.db.Where("id = ?", id).First(&profileUpdate).Error; err != nil {
		return nil, err
	}

	return profileUpdate, nil
}

// FindByProfileID implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) FindByProfileID(id string) (*model.OrganizerProfilesUpdates, error) {
	var update *model.OrganizerProfilesUpdates

	if err := r.db.Where("profile_id = ? AND status = ?", id, "pending").First(&update).Error; err != nil {
		return nil, err
	}

	return update, nil
}

// FindByProfileID implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) FindUpdatesByProfileID(pagination model.Pagination, profileID string) (*model.PaginationRow[*model.OrganizerProfilesUpdates], error) {
	var profileUpdates []*model.OrganizerProfilesUpdates

	if err := r.db.Scopes(Paginate(profileID, &pagination, r.db)).
		Where("profile_id = ?", profileID).
		Find(&profileUpdates).Error; err != nil {
		return nil, err
	}

	return &model.PaginationRow[*model.OrganizerProfilesUpdates]{
		Pagination: pagination,
		Rows:       profileUpdates,
	}, nil
}

// RejectUpdate implements OrganizerProfileUpdateRepo.
func (r *OrganizerProfileUpdateRepository) Validate(id, status string) error {
	return r.db.Model(&model.OrganizerProfilesUpdates{}).
		Where("id = ?", id).
		Updates(&model.OrganizerProfilesUpdates{
			Status:    status,
			UpdatedAt: time.Now().UTC(),
		}).Error
}
