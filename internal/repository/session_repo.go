package repository

import (
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

type SessionRepo interface {
	GetDB() *gorm.DB
	Create(session *model.Sessions) error
	Find(id string) error
	FindByUserID(id string) ([]*model.Sessions, error)
	FindByJTI(jti string) (*model.Sessions, error)
	Update(userId, token string, expire time.Time) error
	Delete(id string) error
	Revoke(jti string) error
}
