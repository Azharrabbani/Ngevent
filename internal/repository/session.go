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

func (r *SessionRepository) Create(session *model.Sessions) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) Find(id string) error {
	var session *model.Sessions

	return r.db.Where("id = ?", id).First(&session).Error
}

func (r *SessionRepository) FindByUserID(id string) (*model.Sessions, error) {
	var session *model.Sessions

	if err := r.db.Where("user_id = ?", id).First(&session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) Update(userId, token, userIP, userAgent string, expire time.Time, update time.Time) error {
	return r.db.Where("user_id = ?", userId).Updates(&model.Sessions{
		RefreshToken: token,
		UserID:       userId,
		IPAddress:    userIP,
		UserAgent:    userAgent,
		ExpiredAt:    expire,
		UpdatedAt:    update}).Error
}

func (r *SessionRepository) Delete(id string) error {
	var session *model.Sessions

	return r.db.Where("id = ?", id).Delete(&session).Error
}

func (r *SessionRepository) DeleteByUserID(id string) error {
	var session *model.Sessions

	return r.db.Where("user_id = ?", id).Delete(&session).Error
}
