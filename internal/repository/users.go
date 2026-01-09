package repository

import (
	"ngevent/internal/model"
	"ngevent/internal/utils/helper"
	"time"

	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Create(users *model.Users) (*model.Users, error) {
	if err := r.db.Create(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) Login(email, password string) (*model.Users, error) {
	var user *model.Users

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	// Validate password
	if err := helper.ValidePassword(user.Password, password); err != nil {
		return nil, err
	}

	// Update user login time
	r.db.Model(&user).Where("id = ?", user.ID).Update("updated_at = ?", time.Now())

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

func (r *UsersRepository) Update(users *model.Users) (*model.Users, error) {
	if err := r.db.Save(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepository) Delete(id string) error {
	var user *model.Users

	return r.db.Where("id = ?", id).Delete(&user).Error
}
