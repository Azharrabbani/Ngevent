package repository

import (
	"ngevent/internal/model"
	"time"
)

type SessionRepo interface {
	Create(session *model.Sessions) error
	Find(id string) error
	FindByUserID(id string) (*model.Sessions, error)
	Update(userId, token, userIP, userAgent string, expire time.Time, update time.Time) error
	Delete(id string) error
	DeleteByUserID(id string) error
}
