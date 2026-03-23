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

type UpdatedEventHandler struct {
	UpdatedEventService *service.UpdatedEventService
	Validate            *validator.Validate
}

func NewUpdatedEventHandler(
	updatedEventService *service.UpdatedEventService,
	validate *validator.Validate,
) *UpdatedEventHandler {
	return &UpdatedEventHandler{
		UpdatedEventService: updatedEventService,
		Validate:            validate,
	}
}

func (h *UpdatedEventHandler) ListAllUpdated(c *fiber.Ctx) error {
	filterReq := new(dto.UpdatedEventFilterReq)
	if err := c.QueryParser(filterReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	if err := h.Validate.Struct(filterReq); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"validation-error",
			msg,
		))
	}

	var start, end time.Time
	if filterReq.Date != 0 {
		unix := time.Unix(filterReq.Date, 0)
		start = time.Date(unix.Year(), unix.Month(), unix.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	}

	var title string
	if filterReq.Title != "" {
		title = utils.CreateSlug(filterReq.Title)
	}

	filter := &dto.UpdatedEventFilter{
		Title:  helper.StrPointerIfNotEmpty(title),
		Status: helper.StrPointerIfNotEmpty(filterReq.Status),
		Start:  helper.TimeToPointer(start),
		End:    helper.TimeToPointer(end),
	}

	pagination := new(model.Pagination)
	if err := c.QueryParser(filterReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Sort:  pagination.Sort,
	}

	updatedEvents, err := h.UpdatedEventService.ListAllUpdatedEvents(filter, *page)
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
		updatedEvents,
	))
}

func (h *UpdatedEventHandler) ListAllUpdatedByEventID(c *fiber.Ctx) error {
	eventID := c.Params("event_id")
	filterReq := new(dto.UpdatedEventFilterReq)
	if err := c.QueryParser(filterReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	if err := h.Validate.Struct(filterReq); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"validation-error",
			msg,
		))
	}

	var start, end time.Time
	if filterReq.Date != 0 {
		unix := time.Unix(filterReq.Date, 0)
		start = time.Date(unix.Year(), unix.Month(), unix.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	}

	var title string
	if filterReq.Title != "" {
		title = utils.CreateSlug(filterReq.Title)
	}

	filter := &dto.UpdatedEventFilter{
		Title:  helper.StrPointerIfNotEmpty(title),
		Status: helper.StrPointerIfNotEmpty(filterReq.Status),
		Start:  helper.TimeToPointer(start),
		End:    helper.TimeToPointer(end),
	}

	pagination := new(model.Pagination)
	if err := c.QueryParser(filterReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Sort:  pagination.Sort,
	}

	updatedEvents, err := h.UpdatedEventService.ListAllUpdatedEventsByEventID(filter, *page, eventID)
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
		updatedEvents,
	))
}

func (h *UpdatedEventHandler) GetUpdatedByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	updatedEvent, err := h.UpdatedEventService.GetUpdateEventByID(id, userID, role)
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
		updatedEvent,
	))
}

func (h *UpdatedEventHandler) ReviewUpdate(c *fiber.Ctx) error {
	var req dto.ReviewUpdatedEventReq
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

	if err := h.UpdatedEventService.ReviewUpdated(&req); err != nil {
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
		"review success",
	))
}

func (h *UpdatedEventHandler) CancelUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.UpdatedEventService.CancelUpdate(id); err != nil {
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
		"updated event canceled",
	))
}
