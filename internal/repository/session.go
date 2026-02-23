package repository

import (
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepo {
	return &SessionRepository{db: db}
}

// GetDB implements SessionRepo.
func (r *SessionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *SessionRepository) Create(session *model.Sessions) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) Find(id string) error {
	var session *model.Sessions

	return r.db.Where("id = ?", id).First(&session).Error
}

func (r *SessionRepository) FindByUserID(id string) ([]*model.Sessions, error) {
	var session []*model.Sessions

	if err := r.db.Where("user_id = ?", id).Find(&session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// FindByJTI implements SessionRepo.
func (r *SessionRepository) FindByJTI(jti string) (*model.Sessions, error) {
	var session *model.Sessions

	if err := r.db.Where("jti = ?", jti).First(&session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) Update(userId, token string, expire time.Time) error {
	return r.db.Where("user_id = ?", userId).Updates(&model.Sessions{RefreshToken: token, ExpiredAt: expire}).Error
}

func (r *SessionRepository) Delete(id string) error {
	var session *model.Sessions

	return r.db.Where("id = ?", id).Delete(&session).Error
}

func (r *SessionRepository) Revoke(jti string) error {
	return r.db.Where("jti = ?", jti).Delete(&model.Sessions{}).Error
}
