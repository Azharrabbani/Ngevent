package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"time"

	"github.com/redis/go-redis/v9"
)

type NewTaskUnverifiedUser interface {
	EnqueueUnverifiedUser(taskType string, payload *model.UnverifiedUserPayload) error
	CancelUnverifiedUser(id string) error
}

type UserService struct {
	UserRepo           repository.UsersRepo
	AttendeeRepo       repository.AttendeeProfilesRepo
	OrganizerRepo      repository.OrganizerProfileRepo
	OtpRepo            repository.OtpRepo
	UserTaskPublisher  NewTaskUnverifiedUser
	OtpTaskPublisher   NewTaskOTP
	EmailTaskPublisher NewTaskEmail
	rdb                *redis.Client
}

func NewUserService(
	userRepo repository.UsersRepo,
	attendeeRepo repository.AttendeeProfilesRepo,
	organizerRepo repository.OrganizerProfileRepo,
	otpRepo repository.OtpRepo,
	userTaskPublisher NewTaskUnverifiedUser,
	otpTaskPublisher NewTaskOTP,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *UserService {
	return &UserService{
		UserRepo:           userRepo,
		AttendeeRepo:       attendeeRepo,
		OrganizerRepo:      organizerRepo,
		OtpRepo:            otpRepo,
		UserTaskPublisher:  userTaskPublisher,
		OtpTaskPublisher:   otpTaskPublisher,
		EmailTaskPublisher: emailTaskPublisher,
		rdb:                rdb,
	}
}

var userCache []string = []string{
	"users:all:*",
}

func (s *UserService) CreateUser(email, password, confirmPassword string) (*dto.UsersResponse, error) {
	userX := s.UserRepo.GetDB().Begin()
	otpX := s.OtpRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			userX.Rollback()
			otpX.Rollback()
		}
	}()

	// Check password
	if password != confirmPassword {
		return nil, errors.New("password not match")
	}

	// Hash the password
	HashPassword, err := helper.HashPassword(password)
	if err != nil {
		return nil, err
	}

	role := string(model.Organizer)

	user := &model.Users{
		Email:    email,
		Role:     &role,
		Password: HashPassword,
	}

	newUser, err := s.UserRepo.Create(user)
	if err != nil {
		userX.Rollback()
		return nil, errors.New("email already registred")
	}

	userResp, err := toUserResponse(newUser)
	if err != nil {
		userX.Rollback()
		return nil, err
	}

	// Generate OTP
	otpCode, err := helper.GenerateOTP()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	otp := helper.NewOTP(
		otpCode,
		newUser.ID,
		"verified_email",
		now.Add(3*time.Minute),
	)

	// Save OTP
	newOTP, err := s.OtpRepo.Create(otp)
	if err != nil {
		userX.Rollback()
		return nil, err
	}

	// Create unverified user task
	// This task function is to delete unverified user
	userPayload := &model.UnverifiedUserPayload{UserID: newUser.ID}
	if err := s.UserTaskPublisher.EnqueueUnverifiedUser(model.TypeVerifiedUser, userPayload); err != nil {
		userX.Rollback()
		return nil, err
	}

	// Create otp task
	// This task function is to delete unused otp
	otpPayload := &model.OTPPayload{OTPID: newOTP.ID}
	if err := s.OtpTaskPublisher.EnqueueOTPVerification(model.TypeVerifiedOTP, otpPayload); err != nil {
		userX.Rollback()
		otpX.Rollback()
		return nil, err
	}

	// Commit all changes
	if err := userX.Commit().Error; err != nil {
		return nil, err
	}

	if err := otpX.Commit().Error; err != nil {
		return nil, err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, userCache)

	// Send to email
	emailPayload := &model.EmailPayload{
		To:  newUser.Email,
		OTP: newOTP.OTP,
	}

	s.EmailTaskPublisher.Enqueue(model.TypeEMailVerify, emailPayload)

	return userResp, nil
}

func (s *UserService) RegisterAdmin(email, password, confirmPassword string) (*dto.UsersResponse, error) {
	// Check password
	if password != confirmPassword {
		return nil, errors.New("password not match")
	}

	// Hash the password
	HashPassword, err := helper.HashPassword(password)
	if err != nil {
		return nil, err
	}

	role := string(model.Admin)

	user := &model.Users{
		Email:      email,
		IsVerified: true,
		Role:       &role,
		Password:   HashPassword,
	}

	newUser, err := s.UserRepo.Create(user)
	if err != nil {
		return nil, errors.New("email already registered")
	}

	userResp, err := toUserResponse(newUser)
	if err != nil {
		return nil, err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, userCache)

	return userResp, nil

}

func (s *UserService) UpdateRole(id, role string) error {
	// Validate user
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role != nil {
		return errors.New("you have already selected a role")
	}

	// Update the role
	user.Role = &role
	return s.UserRepo.UpdateRole(user)
}

func (s *UserService) FindAllUsers(filter *dto.ListUsersReq, pagination model.Pagination) (*model.PaginationRow[*dto.UsersResponse], error) {
	var users *model.PaginationRow[*dto.UsersResponse]

	// Generate cache key
	cacheKey := fmt.Sprintf(
		"users:all:%s:%s:%s:%d:%d:%s",
		helper.StringValue(filter.Role),
		helper.StringValue(filter.Email),
		helper.BoolValue(filter.IsVerified),
		pagination.Limit,
		pagination.Page,
		pagination.Sort,
	)

	// Try get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &users)
	}

	if users == nil {
		// if cache miss, get from db
		users, err = s.UserRepo.FindAll(*filter, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(users); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return users, nil
}

func (s *UserService) FindUserByID(id string) (*dto.UsersResponse, error) {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userResp, err := toUserResponse(user)
	if err != nil {
		return nil, err
	}

	role := helper.StringValue(user.Role)

	if role == string(model.Attendee) {
		// Check if user has profile
		profile, _ := s.AttendeeRepo.FindByUserID(user.ID)
		if profile != nil {
			userResp.HasProfile = helper.BoolPtr(true)
		} else {
			userResp.HasProfile = helper.BoolPtr(false)
		}
	} else {
		profile, _ := s.OrganizerRepo.FindByUserID(user.ID)
		if profile != nil {
			userResp.HasProfile = helper.BoolPtr(true)
		} else {
			userResp.HasProfile = helper.BoolPtr(false)
		}
	}

	if role == string(model.Organizer) {

	}

	return userResp, nil
}

func (s *UserService) FindUser(id string) (*model.Users, error) {
	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) DeleteUnverifiedUser(id string) error {
	if err := s.UserRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func toUserResponse(user *model.Users) (*dto.UsersResponse, error) {
	var userResp *dto.UsersResponse

	if user == nil {
		return nil, errors.New("user data not found")
	}

	var deletedAt int64
	if user.DeletedAt != nil {
		deletedAt = helper.ConvertDatetoUnix(user.DeletedAt.Format(time.RFC3339))
	}

	userResp = &dto.UsersResponse{
		ID:         user.ID,
		Email:      user.Email,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		CreatedAt:  helper.ConvertDatetoUnix(user.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:  helper.ConvertDatetoUnix(user.UpdatedAt.Format(time.RFC3339)),
		DeletedAt:  deletedAt,
	}

	return userResp, nil
}
