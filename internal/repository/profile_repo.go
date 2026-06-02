package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type AttendeeProfilesRepo interface {
	Create(profile *model.AttendeeProfiles) error
	FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) (*model.PaginationRow[*dto.AttendeeProfilesResponse], error)
	HasProfile(userID string) (bool, error)
	FindByID(id string) (*model.AttendeeProfiles, error)
	FindByUserID(id string) (*model.AttendeeProfiles, error)
	Update(profile *model.AttendeeProfiles) error
	UpdatePhotoProfile(userID, photo string) error
	Delete(id string) error
}

type OrganizerProfileRepo interface {
	GetDB() *gorm.DB
	Create(profile *model.OrganizerProfiles) error
	HasProfile(userID string) (bool, error)
	FindAll(pagination model.Pagination, filter *dto.FilterProfileReq) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error)
	FindAllForPublic(pagination model.Pagination, filter *dto.FilterPublicProfileReq) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error)
	FindByID(id string) (*model.OrganizerProfiles, error)
	FindByUserID(userID string) (*model.OrganizerProfiles, error)
	FindByCountry(country string, pagination model.Pagination) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error)
	VerifiedProfile(id string, req *dto.ApprovedReq) error
	RejectProfile(id string, req *dto.RejectedReq) error
	Update(profile *model.OrganizerProfiles) error
	UpdatePhotoProfile(userID, photo string) error
	Delete(id string) error
	SoftDeleteProfile(tx *gorm.DB, profileID string) error
}

type OrganizerProfileUpdateRepo interface {
	GetDB() *gorm.DB
	Create(profile *model.OrganizerProfilesUpdates) error
	FindByID(id string) (*model.OrganizerProfilesUpdates, error)
	FindByProfileID(id string) (*model.OrganizerProfilesUpdates, error)
	FindUpdatesByProfileID(pagination model.Pagination, profileID string) (*model.PaginationRow[*model.OrganizerProfilesUpdates], error)
	Validate(id, status string) error
	Delete(id string) error
	SoftDeleteProfileUpdates(tx *gorm.DB, profileID string) error
}
