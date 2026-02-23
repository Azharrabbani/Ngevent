package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrganizerUpdateHandler struct {
	OrganizerUpdateService *service.OrganizerUpdateService
	Validate               *validator.Validate
}

func NewOrganizerUpdateHandler(
	organizerUpdateService *service.OrganizerUpdateService,
	validate *validator.Validate,
) *OrganizerUpdateHandler {
	return &OrganizerUpdateHandler{
		OrganizerUpdateService: organizerUpdateService,
		Validate:               validate,
	}
}

func (h *OrganizerUpdateHandler) ValidateUpdate(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
			fiber.StatusUnauthorized,
			"failed",
			"error",
			"unauthorized action",
		))
	}

	id := c.Params("id")

	var req *dto.ValidateUpdateReq
	req.UpdateID = id
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
			"error",
			"validation-error",
			msg,
		))
	}

	if err := h.OrganizerUpdateService.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"approval-error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		nil,
	))
}

func (h *OrganizerUpdateHandler) FindUpdateByID(c *fiber.Ctx) error {
	id := c.Params("id")

	update, err := h.OrganizerUpdateService.FindByID(id)
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
		update,
	))
}

func (h *OrganizerUpdateHandler) FindUpdateByProfileID(c *fiber.Ctx) error {
	id := c.Params("id")

	paginate := new(model.Pagination)
	if err := c.QueryParser(paginate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Limit: paginate.Limit,
		Page:  paginate.Page,
		Sort:  paginate.Sort,
	}

	updates, err := h.OrganizerUpdateService.FindByProfileID(id, *page)
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
		updates,
	))
}
