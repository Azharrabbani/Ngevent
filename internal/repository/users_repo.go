package repository

import "ngevent/internal/model"

type UsersRepo interface {
	Create(users *model.Users) (*model.Users, error)
	Login(email, password string) (*model.Users, error)
	FindAll() ([]*model.Users, error)
	FindByID(id string) (*model.Users, error)
	FindByRole(role string) ([]*model.Users, error)
	FindByEmail(email string) (*model.Users, error)
	Update(users *model.Users) (*model.Users, error)
	Delete(id string) error
}
