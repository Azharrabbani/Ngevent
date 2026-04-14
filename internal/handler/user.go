package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"

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

	// Store new user
	user, err := h.UserService.CreateUser(input.Email, input.Password, input.ConfirmPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(
		fiber.StatusCreated,
		"success",
		"register-success",
		user,
	))
}

func (h *UserHandler) SelectRole(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req dto.RoleInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	if err := h.Validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"validation-error",
			msg,
		))
	}

	if err := h.UserService.UpdateRole(userId, req.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	user, err := h.UserService.FindUser(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Error(
			fiber.StatusInternalServerError,
			"failed",
			"error",
			"failed to fetch updated user",
		))
	}

	accessToken, err := helper.GenerateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Error(
			fiber.StatusInternalServerError,
			"failed",
			"error",
			"failed to generate token",
		))
	}

	// update cookie
	c.Cookie(&fiber.Cookie{
		Name:     "ngevent_cookie",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   60 * 60 * 3,
		SameSite: "None",
	})

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"Role selected",
	))
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	filterUser := new(dto.ListUsersReq)
	if err := c.QueryParser(filterUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	paginate := new(model.Pagination)
	if err := c.QueryParser(paginate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Page:  paginate.Page,
		Limit: paginate.Limit,
		Sort:  paginate.Sort,
	}

	users, err := h.UserService.FindAllUsers(filterUser, *page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		users,
	))
}

func (h *UserHandler) FindUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.UserService.FindUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusFound).JSON(dto.Success(
		fiber.StatusFound,
		"success",
		"user-found",
		user,
	))
}

func (h *UserHandler) FindCurrentUser(c *fiber.Ctx) error {
	id := c.Locals("user_id").(string)

	user, err := h.UserService.FindUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"user-found",
		user,
	))
}
