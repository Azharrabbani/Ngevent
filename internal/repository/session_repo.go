package repository

import "ngevent/internal/model"

type SessionRepo interface {
	Create(session *model.Sessions) error
	Find(id string) error
	FindByUserID(id string) (*model.Sessions, error)
	Update(session *model.Sessions) (*model.Sessions, error)
	Delete(id string) error
}