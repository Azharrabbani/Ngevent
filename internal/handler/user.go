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

type UserHandler struct {
	UserService *service.UserService
	Validate    *validator.Validate
}

func NewUserHandler(userService *service.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Validate:    validate,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.Validate.Struct(input); err != nil {
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
	users, err := h.UserService.CreateUser(input.Email, string(hashPassword), input.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	newUser := &model.RegisterResponse{
		ID:         users.ID,
		Email:      users.Email,
		Password:   users.Password,
		Role:       users.Role,
		IsVerified: users.IsVerified,
		CreatedAt:  helper.ConvertDatetoUnix(users.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:  helper.ConvertDatetoUnix(users.UpdatedAt.Format(time.RFC3339)),
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"register-success",
		newUser,
	))
}
