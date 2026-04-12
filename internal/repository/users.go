package repository

import (
	"errors"
	"ngevent/internal/dto"
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

// UpdateRole implements UsersRepo.
func (r *UsersRepository) UpdateRole(user *model.Users) (*model.Users, error) {
	if err := r.db.Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
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

func (r *UsersRepository) FindAll(filter dto.ListUsersReq, pagination model.Pagination) (*model.PaginationRow[*dto.UsersResponse], error) {
	var users []*model.Users

	query := r.db.Scopes(filterUserList(&filter))

	if err := query.Scopes(Paginate(users, &pagination, query)).Find(&users).Error; err != nil {
		return nil, err
	}

	usersResp, err := toUsersResponse(users)
	if err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.UsersResponse]{
		Pagination: pagination,
		Rows:       usersResp,
	}, nil
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
	return r.db.Model(&model.Users{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now().UTC()).Error
}

func filterUserList(filter *dto.ListUsersReq) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Default list
		db = db.Where("is_verified = ?", true)

		// Check optional filters
		if filter.IsVerified != nil {
			db = db.Where("is_verified = ?", false)
		}

		if filter.Email != nil {
			db = db.Where("LOWER(email) LIKE LOWER(?)", "%"+*filter.Email+"%")
		}

		if filter.Role != nil {
			db = db.Where("role = ?", filter.Role)
		}

		return db
	}
}

func toUsersResponse(users []*model.Users) ([]*dto.UsersResponse, error) {
	var usersResp []*dto.UsersResponse

	if len(usersResp) < 0 {
		return nil, errors.New("no data found")
	}

	for _, user := range users {
		usersResp = append(usersResp, &dto.UsersResponse{
			ID:         user.ID,
			Email:      user.Email,
			Role:       user.Role,
			IsVerified: user.IsVerified,
			CreatedAt:  helper.ConvertDatetoUnix(user.CreatedAt.Format(time.RFC3339)),
			UpdatedAt:  helper.ConvertDatetoUnix(user.UpdatedAt.Format(time.RFC3339)),
			DeletedAt:  helper.ConvertDatetoUnix(user.DeletedAt.Format(time.RFC3339)),
		})
	}

	return usersResp, nil
}
