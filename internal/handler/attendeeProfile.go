package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AttendeeProfileHandler struct {
	AttendeeProfileService *service.AttendeeProfileService
}

func NewAttendeeProfileService(attendeeProfileService *service.AttendeeProfileService) *AttendeeProfileHandler {
	return &AttendeeProfileHandler{AttendeeProfileService: attendeeProfileService}
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

	req := &dto.CreateProfileReq{
		UserID:      userID,
		Name:        name,
		Username:    &username,
		PhoneNumber: phoneNumber,
		ISO:         iso,
		Address:     &address,
	}

	if err := h.AttendeeProfileService.Create(photo, req); err != nil {
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
	profileID := c.Params("id")
	userID := c.Locals("user_id").(string)

	if profileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			"missing profile id",
		))
	}

	// Validate user


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

// func (h *AttendeeProfileHandler) UpdateProfile(c *fiber.Ctx) error {

// }
