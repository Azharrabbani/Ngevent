package repository

import (
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type UsersRepo interface {
	GetDB() *gorm.DB
	Create(users *model.Users) (*model.Users, error)
	Login(email, password string) (*model.Users, error)
	FindAll() ([]*model.Users, error)
	FindByID(id string) (*model.Users, error)
	FindByRole(role string) ([]*model.Users, error)
	FindByEmail(email string) (*model.Users, error)
	Update(users *model.Users) (*model.Users, error)
	Delete(id string) error
}
