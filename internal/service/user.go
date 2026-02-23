package service

import (
	"errors"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"
	"time"
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
}

func NewUserService(
	userRepo repository.UsersRepo,
	otpRepo repository.OtpRepo,
	userTaskPublisher NewTaskUnverifiedUser,
	otpTaskPublisher NewTaskOTP,
	emailTaskPublisher NewTaskEmail,
) *UserService {
	return &UserService{
		UserRepo:           userRepo,
		OtpRepo:            otpRepo,
		UserTaskPublisher:  userTaskPublisher,
		OtpTaskPublisher:   otpTaskPublisher,
		EmailTaskPublisher: emailTaskPublisher,
	}
}

func (s *UserService) CreateUser(email, password, role string) (*model.Users, error) {
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

		newUser, err := s.UserRepo.Create(user)
		if err != nil {
			return nil, errors.New("email already registred")
		}

		return newUser, nil
	}

	// User beside admin have to verified their email
	user := &model.Users{
		Email:    email,
		Password: password,
		Role:     role,
	}

	newUser, err := s.UserRepo.Create(user)
	if err != nil {
		return nil, errors.New("email already registred")
	}

	// Generate OTP
	otpCode, err := helper.GenerateOTP()
	if err != nil {
		return nil, err
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

	// Send to email
	emailPayload := &model.EmailPayload{
		To:    newUser.Email,
		OTP:   newOTP.OTP,
		OTPID: newOTP.ID,
	}

	s.EmailTaskPublisher.Enqueue(model.TypeEMailVerify, emailPayload)

	return newUser, nil
}

func (s *UserService) DeleteUnverifiedUser(id string) error {
	return s.UserRepo.Delete(id)
}
