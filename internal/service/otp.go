package service

import (
	"errors"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OTPService struct {
	UserRepo           repository.UsersRepo
	OtpRepo            repository.OtpRepo
	OtpTaskPublisher   NewTaskOTP
	UserTaskPublisher  NewTaskUnverifiedUser
	EmailTaskPublisher NewTaskEmail
}

func NewOTPService(
	userRepo repository.UsersRepo,
	otpRepo repository.OtpRepo,
	otpTaskPublisher NewTaskOTP,
	userTaskPublisher NewTaskUnverifiedUser,
	emailTaskPublisher NewTaskEmail,
) *OTPService {
	return &OTPService{
		UserRepo:           userRepo,
		OtpRepo:            otpRepo,
		OtpTaskPublisher:   otpTaskPublisher,
		UserTaskPublisher:  userTaskPublisher,
		EmailTaskPublisher: emailTaskPublisher,
	}
}

func (s *OTPService) ResendOTPCode(email string) (int, error) {
	otpX := s.OtpRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		otpX.Rollback()
	}()

	// Find user
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return fiber.StatusNotFound, errors.New("user with this email not found")
	}

	// Delete the task
	// Cancel unverified user task
	if err := s.UserTaskPublisher.CancelUnverifiedUser(user.ID); err != nil {
		return fiber.StatusBadGateway, err
	}

	// Generate new otp
	otpCode, err := helper.GenerateOTP()
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Save OTP
	now := time.Now().UTC()
	otp := helper.NewOTP(
		otpCode,
		user.ID,
		"verified_email",
		now.Add(3*time.Minute),
	)

	newOTP, err := s.OtpRepo.Create(otp)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Create unverified user task
	// This task function is to delete unverified user
	userPayload := &model.UnverifiedUserPayload{UserID: user.ID}
	if err := s.UserTaskPublisher.EnqueueUnverifiedUser(model.TypeVerifiedUser, userPayload); err != nil {
		return fiber.StatusBadGateway, err
	}

	// Create otp task
	// This task function is to delete unused otp
	otpPayload := &model.OTPPayload{OTPID: newOTP.ID}
	if err := s.OtpTaskPublisher.EnqueueOTPVerification(model.TypeVerifiedOTP, otpPayload); err != nil {
		otpX.Rollback()
		return fiber.StatusBadGateway, err
	}

	// Commit changes
	if err := otpX.Commit().Error; err != nil {
		return fiber.StatusBadGateway, err
	}

	emailPayload := &model.EmailPayload{
		To:    email,
		OTP:   newOTP.OTP,
	}
	// Send new email
	s.EmailTaskPublisher.Enqueue(model.TypeEMailVerify, emailPayload)

	return 0, nil
}
