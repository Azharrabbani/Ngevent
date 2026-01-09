package service

import (
	"ngevent/internal/model"
	"ngevent/internal/repository"
)

type SessionService struct {
	sessionRepo repository.SessionRepo
}

func NewSessionService(sessionRepo repository.SessionRepo) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) Create(session *model.Sessions)  error {
	return s.sessionRepo.Create(session)
}

func (s *SessionService) Find(id string) error {
	return s.sessionRepo.Find(id)
}

func (s *SessionService) FindByUserID(id string) (*model.Sessions, error) {
	return s.sessionRepo.FindByUserID(id)
}

func (s *SessionService) Update(session *model.Sessions) (*model.Sessions, error) {
	return s.sessionRepo.Update(session)
}

func (s *SessionService) Delete(id string) error {
	return s.sessionRepo.Delete(id)
}
