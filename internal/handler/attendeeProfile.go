package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/service"
	"ngevent/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AttendeeProfileHandler struct {
	AttendeeProfileService *service.AttendeeProfileService
	Validate               *validator.Validate
}

func NewAttendeeProfileService(
	attendeeProfileService *service.AttendeeProfileService,
	validate *validator.Validate,
) *AttendeeProfileHandler {
	return &AttendeeProfileHandler{
		AttendeeProfileService: attendeeProfileService,
		Validate:               validate,
	}
}

func (h *AttendeeProfileHandler) CreateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	photo, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	name := c.FormValue("name")
	username := c.FormValue("username")
	phoneNumber := c.FormValue("phonenumber")
	iso := c.FormValue("iso")
	address := c.FormValue("address")

	req := &dto.CreateAttendeeProfileReq{
		UserID:       userID,
		PhotoProfile: photo,
		Name:         name,
		Username:     &username,
		PhoneNumber:  phoneNumber,
		ISO:          iso,
		Address:      &address,
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

	if err := h.AttendeeProfileService.Create(req); err != nil {
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
		"success",
		"success create the profile",
	))
}

func (h *AttendeeProfileHandler) HasProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	hasProfile, err := h.AttendeeProfileService.HasProfile(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Error(
			fiber.StatusInternalServerError,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		hasProfile,
	))
}
func (h *AttendeeProfileHandler) GetProfileByID(c *fiber.Ctx) error {
	id := c.Params("id")

	profile, err := h.AttendeeProfileService.FindByID(id)
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
		"success",
		profile,
	))
}

func (h *AttendeeProfileHandler) GetProfileByUserID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	profile, err := h.AttendeeProfileService.FindByUserID(userID)
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
		"success",
		profile,
	))
}

func (h *AttendeeProfileHandler) UpdateProfilePhoto(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	photo, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	status, err := h.AttendeeProfileService.UpdatePhotoProfile(photo, userID)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"photo profile have been updated",
	))
}

func (h *AttendeeProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req *dto.UpdateAttendeeProfileReq

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
			"error",
			msg,
		))
	}

	status, err := h.AttendeeProfileService.UpdateProfile(userID, req)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"update success",
	))
}
