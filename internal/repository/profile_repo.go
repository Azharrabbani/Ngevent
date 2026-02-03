package repository

import (
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
