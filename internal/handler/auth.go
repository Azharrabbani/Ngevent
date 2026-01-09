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

type UsersHandler struct {
	userService *service.UsersService
	validate    *validator.Validate
}

func NewUserHandler(userService *service.UsersService, validate *validator.Validate) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *UsersHandler) Register(c *fiber.Ctx) error {
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

func (h *UsersHandler) Login(c *fiber.Ctx) error {
	// Validate the req
	var req dto.LoginInput

	// Get user client
	ip := c.IP()
	userAgent := Handler(c)

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

	user, accessToken, err := h.userService.Login(client, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
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

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"login-success",
		userLogin,
	))
}