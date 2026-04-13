package service

import (
	"errors"
	"ngevent/internal/dto"
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

func (s *AuthService) VerififyEmail(email, otpInput string) (int, error) {
	otpX := s.otpRepo.GetDB().Begin()
	authX := s.userRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			otpX.Rollback()
			authX.Rollback()
		}
	}()

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return fiber.StatusNotFound, errors.New("user not found")
	}

	// Check OTP
	otp, err := s.otpRepo.FindByUserID(user.ID)
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

func (s *AuthService) Login(client *model.Client, req dto.LoginInput) (*model.Users, string, string, time.Time, error) {
	sesX := s.sessionRepo.GetDB().Begin()

	// Rollback if failed
	defer func() {
		if r := recover(); r != nil {
			sesX.Rollback()
		}
	}()

	// Login by checking user account
	user, err := s.userRepo.Login(req.Email, req.Password)
	if err != nil {
		return nil, "", "", time.Time{}, errors.New("incorrect password or user not found")
	}

	// is user verified
	if !user.IsVerified {
		return nil, "", "", time.Time{}, errors.New("user not verified")
	}

	// Generate access token
	accessToken, err := helper.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", time.Time{}, err
	}

	// Determine refresh token expire time
	var refreshExpire time.Time
	if req.RememberMe {
		refreshExpire = time.Now().UTC().Add(7 * 24 * time.Hour)
	} else {
		refreshExpire = time.Now().UTC().Add(24 * time.Hour)
	}

	// Generate refresh token
	refreshToken, jti, err := helper.GenerateRefreshToken(user.ID, refreshExpire)
	if err != nil {
		return nil, "", "", time.Time{}, err
	}

	// Save new session
	session := &model.Sessions{
		UserID:       user.ID,
		JTI:          jti,
		RefreshToken: refreshToken,
		ExpiredAt:    refreshExpire,
		IPAddress:    client.IP,
		UserAgent:    client.UserAgent,
	}

	if err := sesX.Create(session).Error; err != nil {
		sesX.Rollback()
		return nil, "", "", time.Time{}, err
	}

	// Commit all change
	if err := sesX.Commit().Error; err != nil {
		sesX.Rollback()
		return nil, "", "", time.Time{}, err
	}

	return user, accessToken, refreshToken, refreshExpire, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// Parse refresh token
	userID, jti, err := helper.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Check session
	session, err := s.sessionRepo.FindByJTI(jti)
	if err != nil {
		return "", errors.New("session not found")
	}

	// Check expired
	now := time.Now().UTC()
	if now.After(session.ExpiredAt.UTC()) {
		return "", errors.New("refresh token expired")
	}

	// Generate new access token
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	return helper.GenerateAccessToken(user)
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
		return fiber.StatusBadRequest, errors.New("password not match")
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

func (s *AuthService) Logout(refreshToken string) error {
	// Parse refresh token
	_, jti, err := helper.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return s.sessionRepo.Revoke(jti)
}

func (s *AuthService) DeleteUnusedOTP(id string) error {
	return s.otpRepo.Delete(id)
}
