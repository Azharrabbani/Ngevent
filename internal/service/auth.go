package service

import (
	"errors"
	"fmt"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

type UsersService struct {
	userRepo    repository.UsersRepo
	sessionRepo repository.SessionRepo
	otpRepo     repository.OtpRepo
}

func NewUsersService(
	userRepo repository.UsersRepo,
	sessionRepo repository.SessionRepo,
	otpRepo repository.OtpRepo) *UsersService {
	return &UsersService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		otpRepo:     otpRepo,
	}
}

func (s *UsersService) CreateUser(email, password, role string) (*model.Users, error) {

	if role == "admin" {
		user := &model.Users{
			Email:      email,
			Password:   password,
			Role:       role,
			IsVerified: true,
		}

		newUser, err := s.userRepo.Create(user)
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

	newUser, err := s.userRepo.Create(user)
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
		"reset_password",
		now.Add(5*time.Minute),
	)

	// Save OTP
	newOTP, err := s.otpRepo.Create(otp)
	if err != nil {
		return nil, err
	}

	urlHost := os.Getenv("APP_HOST")
	urlPort := os.Getenv("APP_PORT")

	// Send to email
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Verifify Email")

	verifyLink := fmt.Sprintf(
		"%s:%s/api/v1/verify-email/%s",
		urlHost,
		urlPort,
		newOTP.ID,
	)

	utils.VerifyEmailMail(m, verifyLink, newOTP.OTP)

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	go func() {
		d := gomail.NewDialer(host, port, username, smtpPassword)
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}()

	return newUser, nil
}

func (s *UsersService) VerififyEmail(id, otpInput string) (int, error) {
	// Check OTP
	otp, err := s.otpRepo.FindByID(id)
	if err != nil {
		return fiber.StatusNotFound, errors.New("otp not found")
	}

	if otp.OTP != otpInput || otp.IsUsed {
		return fiber.StatusBadRequest, errors.New("otp expired or not valid")
	}

	// Update OTP status
	otp.IsUsed = true
	otp.ExpiredAt = time.Now().UTC()

	_, err = s.otpRepo.Update(otp)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	// Update user verification
	user, err := s.userRepo.FindByID(otp.UserID)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	user.IsVerified = true
	user.UpdatedAt = time.Now().UTC()

	_, err = s.userRepo.Update(user)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	return 0, nil
}

func (s *UsersService) Login(client *model.Client, email, password string) (*model.Users, int, string, error) {
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
		return nil, fiber.StatusBadRequest, "", err
	}

	// Set New Expired
	expiredAt := time.Now().Add(time.Hour * 24 * 7).UTC()
	updateAt := time.Now().UTC()
	userIP := client.IP
	userAgent := client.UserAgent

	// Update session
	s.sessionRepo.Update(user.ID, refreshToken, userIP, userAgent, expiredAt, updateAt)

	return user, 0, accessToken, nil
}

func (s *UsersService) ForgotPassword(email string) (int, error) {
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
		now.Add(5*time.Minute),
	)

	newOTP, err := s.otpRepo.Create(otp)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	urlHost := os.Getenv("APP_HOST")
	urlPort := os.Getenv("APP_PORT")

	// Send to email
	m := gomail.NewMessage()
	m.SetHeader("From", "ngevent@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Reset Password")

	resetLink := fmt.Sprintf(
		"%s:%s/api/v1/reset-password/%s",
		urlHost,
		urlPort,
		newOTP.ID,
	)

	utils.ForgotPasswordMail(m, resetLink)

	// SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	go func() {
		d := gomail.NewDialer(host, port, username, smtpPassword)
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}()

	return 0, nil
}

func (s *UsersService) ResetPassword(id, newPassword, confirmPassword string) (int, error) {
	// Check the OTP
	userOTP, err := s.otpRepo.FindByID(id)
	if err != nil {
		return fiber.StatusNotFound, errors.New("otp not found")
	}

	user, err := s.userRepo.FindByID(userOTP.UserID)
	if err != nil {
		return fiber.StatusNotFound, err
	}

	// Check if otp expired
	if time.Now().UTC().After(userOTP.ExpiredAt) || userOTP.IsUsed {
		return fiber.StatusBadRequest, errors.New("otp expired")
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

	// Update user password
	user.Password = newHashPassword
	user.UpdatedAt = time.Now().UTC()

	_, err = s.userRepo.Update(user)
	if err != nil {
		return fiber.StatusBadRequest, errors.New("failed to reset password")
	}

	return 0, nil
}

func (s *UsersService) Logout(id string) error {
	return s.sessionRepo.DeleteByUserID(id)
}
