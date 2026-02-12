package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
)

type AttendeeProfilesRepo interface {
	Create(profile *model.AttendeeProfiles) error
	FindByID(id string) (*model.AttendeeProfiles, error)
	FindByUserID(id string) (*model.AttendeeProfiles, error)
	Update(profile *model.AttendeeProfiles) error
	UpdatePhotoProfile(userID, photo string) error
	Delete(id string) error
}

type OrganizerProfileRepo interface {
	Create(profile *model.OrganizerProfiles) error
	FindByID(id string) (*model.OrganizerProfiles, error)
	FindByUserID(userID string) (*model.OrganizerProfiles, error)
	FindByCountry(country string, pagination model.Pagination) (*model.PaginationRow[*dto.OrganizerProfilesResponse], error)
	VerifiedProfile(id string) error
	RejectProfile(id string) error
	Update(profile *model.OrganizerProfiles) error
	UpdatePhotoProfile(userID, photo string) error
	Delete(id string) error
}
