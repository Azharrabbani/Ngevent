package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService *service.UsersService
	validate    *validator.Validate
}

func NewAuthHandler(userService *service.UsersService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(input); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"validation-error",
			msg,
		))
	}

	// Hashing password
	hashPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	// Store new user
	users, err := h.userService.CreateUser(input.Email, string(hashPassword), input.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"register-success",
		users,
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
			"failed",
			"validation-error",
			msg,
		))
	}

	user, status, accessToken, err := h.userService.Login(client, req.Email, req.Password)
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

	status, err := h.userService.ForgotPassword(req.Email)
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

	status, err := h.userService.ResetPassword(otpID, req.NewPassword, req.ConfirmPassword)
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

	if err := h.userService.Logout(userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"invalid-request",
			err.Error(),
		))
	}

	c.ClearCookie()
	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"success",
		"Logout success",
	))
}
