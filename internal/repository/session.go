package repository

import (
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
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

func (r *SessionRepository) Update(session *model.Sessions) (*model.Sessions, error) {
	if err := r.db.Save(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) Delete(id string) error {
	var session *model.Sessions

	return r.db.Where("id = ?", id).Delete(&session).Error
}

func (r *SessionRepository) DeleteByUserID(id string) error {
	var session *model.Sessions

	return r.db.Where("user_id = ?", id).Delete(&session).Error
}
