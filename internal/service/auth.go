package service

import (
	"errors"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NewTaskOTP interface {
	EnqueueOTPVerification(taskType string, payload *model.OTPPayload) error
	CancelOTPVerification(id string) error
}

type NewTaskEmail interface {
	Enqueue(taskType string, payload interface{}) error
}

type AuthService struct {
	userRepo           repository.UsersRepo
	sessionRepo        repository.SessionRepo
	otpRepo            repository.OtpRepo
	UserTaskPublisher  NewTaskUnverifiedUser
	OtpTaskPublisher   NewTaskOTP
	EmailTaskPublisher NewTaskEmail
}

func NewAuthService(
	userRepo repository.UsersRepo,
	sessionRepo repository.SessionRepo,
	otpRepo repository.OtpRepo,
	userTaskPublisher NewTaskUnverifiedUser,
	otpTaskPublisher NewTaskOTP,
	emailTaskPublisher NewTaskEmail,
) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		sessionRepo:        sessionRepo,
		otpRepo:            otpRepo,
		UserTaskPublisher:  userTaskPublisher,
		OtpTaskPublisher:   otpTaskPublisher,
		EmailTaskPublisher: emailTaskPublisher,
	}
}

func (s *AuthService) VerififyEmail(id, otpInput string) (int, error) {
	otpX := s.otpRepo.GetDB().Begin()
	authX := s.userRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			otpX.Rollback()
			authX.Rollback()
		}
	}()

	// Check OTP
	otp, err := s.otpRepo.FindByID(id)
	if err != nil {
		return fiber.StatusNotFound, errors.New("otp not found")
	}

	if otp.OTP != otpInput || otp.IsUsed || time.Now().UTC().After(otp.ExpiredAt) {
		return fiber.StatusBadRequest, errors.New("otp expired or not valid")
	}

	// Update OTP status
	otp.IsUsed = true
	// otp.ExpiredAt = time.Now().UTC()

	_, err = s.otpRepo.Update(otp)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Cancel unused otp task
	if err := s.OtpTaskPublisher.CancelOTPVerification(otp.ID); err != nil {
		otpX.Rollback()
		return fiber.StatusBadGateway, err
	}

	// Update user verification
	user, err := s.userRepo.FindByID(otp.UserID)
	if err != nil {
		otpX.Rollback()
		return fiber.StatusBadRequest, err
	}

	user.IsVerified = true
	user.UpdatedAt = time.Now().UTC()

	_, err = s.userRepo.Update(user)
	if err != nil {
		otpX.Rollback()
		return fiber.StatusBadRequest, err
	}

	// Cancel unverified user task
	if err := s.UserTaskPublisher.CancelUnverifiedUser(user.ID); err != nil {
		otpX.Rollback()
		authX.Rollback()
		return fiber.StatusBadGateway, err
	}

	// Commit all changes
	if err := otpX.Commit().Error; err != nil {
		return fiber.StatusBadGateway, err
	}

	if err := authX.Commit().Error; err != nil {
		return fiber.StatusBadGateway, err
	}

	return 0, nil
}

func (s *AuthService) Login(client *model.Client, email, password string) (*model.Users, int, string, error) {
	authX := s.userRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			authX.Rollback()
		}
	}()

	// Login by checking user account
	user, err := s.userRepo.Login(email, password)
	if err != nil {
		return nil, fiber.StatusNotFound, "", errors.New("incorrect password or user not found")
	}

	// is user verified
	if !user.IsVerified {
		return nil, fiber.StatusBadRequest, "", errors.New("user not verified")
	}

	// Check user session
	userSession, err := s.sessionRepo.FindByUserID(user.ID)
	if err != nil {
		// Session not found -> create new session
		accessToken, refreshToken, err := helper.GenerateToken(user)
		if err != nil {
			return nil, fiber.StatusBadRequest, "", err
		}

		newSession := &model.Sessions{
			UserID:       user.ID,
			RefreshToken: refreshToken,
			IPAddress:    client.IP,
			ExpiredAt:    time.Now().Add(time.Hour * 24 * 7).UTC(),
			UserAgent:    client.UserAgent,
		}

		// Save the session
		if err := s.sessionRepo.Create(newSession); err != nil {
			authX.Rollback()
			return nil, fiber.StatusBadRequest, "", err
		}

		return user, 0, accessToken, nil
	}

	// Session exist -> check if session expired
	if time.Now().After(userSession.ExpiredAt) {
		if err := s.sessionRepo.Delete(userSession.ID); err != nil {
			return nil, fiber.StatusBadRequest, "", err
		}

		return nil, fiber.StatusBadRequest, "", errors.New("session expired, please login again")
	}

	// Session valid -> update token
	accessToken, refreshToken, err := helper.GenerateToken(user)
	if err != nil {
		authX.Rollback()
		return nil, fiber.StatusBadRequest, "", err
	}

	// Set New Expired
	expiredAt := time.Now().Add(time.Hour * 24 * 7).UTC()
	updateAt := time.Now().UTC()
	userIP := client.IP
	userAgent := client.UserAgent

	// Update session
	if err := s.sessionRepo.Update(user.ID, refreshToken, userIP, userAgent, expiredAt, updateAt); err != nil {
		authX.Rollback()
		return nil, fiber.StatusBadRequest, "", err
	}

	// Commit all changes
	if err := authX.Commit().Error; err != nil {
		return nil, fiber.StatusBadGateway, "", err
	}

	return user, 0, accessToken, nil
}

func (s *AuthService) ForgotPassword(email string) (int, error) {
	otpX := s.otpRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			otpX.Rollback()
		}
	}()

	// Search the user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return fiber.StatusNotFound, errors.New("user with this email not found")
	}

	// Generate OTP
	otpCode, err := helper.GenerateOTP()
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Save the otp
	now := time.Now().UTC()
	otp := helper.NewOTP(
		otpCode,
		user.ID,
		"reset_password",
		now.Add(3*time.Minute),
	)

	newOTP, err := s.otpRepo.Create(otp)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Create otp task
	// This task function is to delete unused user
	payload := &model.OTPPayload{OTPID: newOTP.ID}
	if err := s.OtpTaskPublisher.EnqueueOTPVerification(model.TypeVerifiedOTP, payload); err != nil {
		otpX.Rollback()
		return fiber.StatusBadGateway, err
	}

	// Commit all changes
	if err := otpX.Commit().Error; err != nil {
		return fiber.StatusBadGateway, err
	}

	emailPayload := &model.EmailPayload{
		To:    user.Email,
		OTPID: newOTP.ID,
	}

	// Send to email
	s.EmailTaskPublisher.Enqueue(model.TypeEmailForgetPassword, emailPayload)

	return 0, nil
}

func (s *AuthService) ResetPassword(id, newPassword, confirmPassword string) (int, error) {
	otpX := s.otpRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			otpX.Rollback()
		}
	}()

	// Check the OTP
	userOTP, err := s.otpRepo.FindByID(id)
	if err != nil {
		return fiber.StatusNotFound, errors.New("otp expired or not found")
	}

	user, err := s.userRepo.FindByID(userOTP.UserID)
	if err != nil {
		return fiber.StatusNotFound, err
	}

	// Check if otp expired
	if time.Now().UTC().After(userOTP.ExpiredAt) || userOTP.IsUsed {
		return fiber.StatusBadRequest, errors.New("otp expired or invalid")
	}

	if newPassword != confirmPassword {
		return fiber.StatusBadRequest, errors.New("passwords do not match")
	}

	newHashPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Update the otp status
	userOTP.IsUsed = true
	userOTP.UpdatedAt = time.Now().UTC()

	_, err = s.otpRepo.Update(userOTP)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Cancel unused otp task
	if err := s.OtpTaskPublisher.CancelOTPVerification(userOTP.ID); err != nil {
		otpX.Rollback()
		return fiber.StatusBadGateway, err
	}

	// Update user password
	user.Password = newHashPassword
	user.UpdatedAt = time.Now().UTC()

	_, err = s.userRepo.Update(user)
	if err != nil {
		otpX.Rollback()
		return fiber.StatusBadRequest, errors.New("failed to reset password")
	}

	// Commit all changes
	if err := otpX.Commit().Error; err != nil {
		return fiber.StatusBadGateway, err
	}

	return 0, nil
}

func (s *AuthService) Logout(id string) error {
	return s.sessionRepo.DeleteByUserID(id)
}

func (s *AuthService) DeleteUnusedOTP(id string) error {
	return s.otpRepo.Delete(id)
}
