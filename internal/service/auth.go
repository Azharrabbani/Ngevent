package service

import (
	"errors"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"
	"time"
)

type UsersService struct {
	userRepo    repository.UsersRepo
	sessionRepo repository.SessionRepo
}

func NewUsersService(userRepo repository.UsersRepo, sessionRepo repository.SessionRepo) *UsersService {
	return &UsersService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *UsersService) CreateUser(email, password, role string) (*model.Users, error) {
	user := &model.Users{
		Email:      email,
		Password:   password,
		Role:       role,
		IsVerified: true, // For now
	}

	newUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("email already registred")
	}

	return newUser, nil
}

func (s *UsersService) Login(client *model.Client, email, password string) (*model.Users, string, error) {
	// Login by checking user account
	user, err := s.userRepo.Login(email, password)
	if err != nil {
		return nil, "", errors.New("incorrect password or user not found")
	}

	// is user verified
	if !user.IsVerified {
		return nil, "", errors.New("user not verified")
	}

	// Check user session
	userSession, err := s.sessionRepo.FindByUserID(user.ID)
	if err != nil {
		// Session not found -> create new session
		accessToken, refreshToken, err := helper.GenerateToken(user)
		if err != nil {
			return nil, "", err
		}

		newSession := &model.Sessions{
			UserID:       user.ID,
			RefreshToken: refreshToken,
			IPAddress:    client.IP,
			ExpiredAt:    time.Now().Add(time.Hour * 24 * 7),
			UserAgent:    client.UserAgent,
		}

		// Save the session
		if err := s.sessionRepo.Create(newSession); err != nil {
			return nil, "", err
		}

		return user, accessToken, nil
	}

	// Session exist -> check if session expired
	if time.Now().After(userSession.ExpiredAt) {
		if err := s.sessionRepo.Delete(userSession.ID); err != nil {
			return nil, "", err
		}

		return nil, "", errors.New("session expired, please login again")
	}

	// Session valid -> update token
	accessToken, refreshToken, err := helper.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	// Update session
	updateSession := &model.Sessions{
		ID:           userSession.ID,
		UserID:       userSession.UserID,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 7),
		UpdatedAt:    time.Now(),
	}
	s.sessionRepo.Update(updateSession)

	return user, accessToken, nil
}

func (s *UsersService) Logout(id string) error {
	return s.sessionRepo.DeleteByUserID(id)
}
