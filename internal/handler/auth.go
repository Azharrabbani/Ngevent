package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		validate:    validate,
	}
}

func (h *AuthHandler) ListPhoneCodes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	phoneCodes := utils.ListAllPhoneCodes(page, limit)

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		phoneCodes,
	))
}

func (h *AuthHandler) VerififyEmail(c *fiber.Ctx) error {
	otpID := c.Params("id")

	var req dto.VerifyEmailInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"validation-error",
			msg,
		))
	}

	// Verify Email
	status, err := h.AuthService.VerififyEmail(otpID, req.OTP)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"verify-error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"success verified email",
	))

}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Validate the req
	var req dto.LoginInput

	// Get user client
	ip := c.IP()
	userAgent := utils.Handler(c)

	client := &model.Client{
		IP:        ip,
		UserAgent: userAgent,
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
		))
	}

	// Login
	user, status, accessToken, err := h.AuthService.Login(client, req.Email, req.Password)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"login-failed",
			err.Error(),
		))
	}

	userLogin := &model.LoginResponse{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		AccessToken: accessToken,
		LoginAt:     helper.ConvertDatetoUnix(user.UpdatedAt.Format(time.RFC3339)),
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "ngevent-cookie",
		Value:    accessToken,
		HTTPOnly: true,
		MaxAge:   60 * 60 * 3,
	})

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"login-success",
		userLogin,
	))
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgetPasswordInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
		))
	}

	// Forgot password
	status, err := h.AuthService.ForgotPassword(req.Email)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"error",
			"invalid-request",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusAccepted).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"success",
		"Reset link have been send to your email",
	))
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	otpID := c.Params("id")

	var req dto.ResetPasswordInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
		))
	}

	// Reset password
	status, err := h.AuthService.ResetPassword(otpID, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"password have been reset",
	))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// Logout
	if err := h.AuthService.Logout(userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"invalid-request",
			err.Error(),
		))
	}

	// Clear cookie
	c.ClearCookie()

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"success",
		"Logout success",
	))
}
