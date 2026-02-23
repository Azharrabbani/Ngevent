package repository

import (
	"errors"
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"time"

	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepo {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *UsersRepository) Create(users *model.Users) (*model.Users, error) {
	if err := r.db.Create(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) Login(email, password string) (*model.Users, error) {
	var user *model.Users

	if err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		return nil, err
	}

	// Validate password
	if err := helper.ValidePassword(user.Password, password); err != nil {
		return nil, errors.New("password not valid")
	}

	return user, nil
}

func (r *UsersRepository) FindAll() ([]*model.Users, error) {
	var users []*model.Users

	if err := r.db.Where("is_verified = ?", true).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) FindByID(id string) (*model.Users, error) {
	var user *model.Users

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UsersRepository) FindByRole(role string) ([]*model.Users, error) {
	var users []*model.Users

	if err := r.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) FindByEmail(email string) (*model.Users, error) {
	var user *model.Users

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UsersRepository) Update(users *model.Users) (*model.Users, error) {
	if err := r.db.Save(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Update("deleted_at", time.Now().UTC()).Error
}
