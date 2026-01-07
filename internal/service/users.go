package service

import (
	"errors"
	"ngevent/internal/model"
	"ngevent/internal/repository"
)

type UsersService struct {
	userRepo repository.UsersRepo
}

func NewUsersService(userRepo repository.UsersRepo) *UsersService {
	return &UsersService{userRepo: userRepo}
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
