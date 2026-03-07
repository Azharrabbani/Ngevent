package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
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
	OtpRepo            repository.OtpRepo
	UserTaskPublisher  NewTaskUnverifiedUser
	OtpTaskPublisher   NewTaskOTP
	EmailTaskPublisher NewTaskEmail
	rdb                *redis.Client
}

func NewUserService(
	userRepo repository.UsersRepo,
	otpRepo repository.OtpRepo,
	userTaskPublisher NewTaskUnverifiedUser,
	otpTaskPublisher NewTaskOTP,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *UserService {
	return &UserService{
		UserRepo:           userRepo,
		OtpRepo:            otpRepo,
		UserTaskPublisher:  userTaskPublisher,
		OtpTaskPublisher:   otpTaskPublisher,
		EmailTaskPublisher: emailTaskPublisher,
		rdb:                rdb,
	}
}

func (s *UserService) InvalidateCache() {
	ctx := context.Background()

	patterns := []string{
		"users:all:*",
	}

	for _, pattern := range patterns {
		iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			s.rdb.Del(ctx, iter.Val())
		}
	}

	// Use SCAN for pattern keys to avoid blocking
	log.Println("[CACHE] users cache invalidated")
}

func (s *UserService) CreateUser(email, password, role string) error {
	userX := s.UserRepo.GetDB().Begin()
	otpX := s.OtpRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			userX.Rollback()
			otpX.Rollback()
		}
	}()

	if role == "admin" {
		user := &model.Users{
			Email:      email,
			Password:   password,
			Role:       role,
			IsVerified: true,
		}

		if err := userX.Create(user); err != nil {
			return errors.New("email already registred")
		}

		// Commit all changes
		if err := userX.Commit().Error; err != nil {
			return err
		}

		if err := otpX.Commit().Error; err != nil {
			return err
		}

		// Invalidate cache after update
		s.InvalidateCache()

		return nil
	}

	// User beside admin have to verified their email
	user := &model.Users{
		Email:    email,
		Password: password,
		Role:     role,
	}

	newUser, err := s.UserRepo.Create(user)
	if err != nil {
		errors.New("email already registred")
	}

	// Generate OTP
	otpCode, err := helper.GenerateOTP()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	otp := helper.NewOTP(
		otpCode,
		user.ID,
		"verified_email",
		now.Add(3*time.Minute),
	)

	// Save OTP
	newOTP, err := s.OtpRepo.Create(otp)
	if err != nil {
		userX.Rollback()
		return err
	}

	// Create unverified user task
	// This task function is to delete unverified user
	userPayload := &model.UnverifiedUserPayload{UserID: newUser.ID}
	if err := s.UserTaskPublisher.EnqueueUnverifiedUser(model.TypeVerifiedUser, userPayload); err != nil {
		userX.Rollback()
		return err
	}

	// Create otp task
	// This task function is to delete unused otp
	otpPayload := &model.OTPPayload{OTPID: newOTP.ID}
	if err := s.OtpTaskPublisher.EnqueueOTPVerification(model.TypeVerifiedOTP, otpPayload); err != nil {
		userX.Rollback()
		otpX.Rollback()
		return err
	}

	// Commit all changes
	if err := userX.Commit().Error; err != nil {
		return err
	}

	if err := otpX.Commit().Error; err != nil {
		return err
	}

	// Invalidate cache after update
	s.InvalidateCache()

	// Send to email
	emailPayload := &model.EmailPayload{
		To:    newUser.Email,
		OTP:   newOTP.OTP,
		OTPID: newOTP.ID,
	}

	s.EmailTaskPublisher.Enqueue(model.TypeEMailVerify, emailPayload)

	return nil
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

	return userResp, nil
}

func (s *UserService) DeleteUnverifiedUser(id string) error {
	if err := s.UserRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache after update
	s.InvalidateCache()

	return nil
}

func toUserResponse(user *model.Users) (*dto.UsersResponse, error) {
	var userResp *dto.UsersResponse

	if user == nil {
		return nil, errors.New("user data not found")
	}

	userResp = &dto.UsersResponse{
		ID:         user.ID,
		Email:      user.Email,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		CreatedAt:  helper.ConvertDatetoUnix(user.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:  helper.ConvertDatetoUnix(user.UpdatedAt.Format(time.RFC3339)),
		DeletedAt:  helper.ConvertDatetoUnix(user.DeletedAt.Format(time.RFC3339)),
	}

	return userResp, nil
}
