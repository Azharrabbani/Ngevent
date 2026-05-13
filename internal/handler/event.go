package handler

import (
	"encoding/json"
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

type EventHandler struct {
	EventService *service.EventService
	Validate     *validator.Validate
}

func NewEventHandler(
	eventService *service.EventService,
	validate *validator.Validate,
) *EventHandler {
	return &EventHandler{
		EventService: eventService,
		Validate:     validate,
	}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	data := c.FormValue("data")

	var req dto.EventReq
	if err := json.Unmarshal([]byte(data), &req); err != nil {
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
	req.UserID = userID

	banner, _ := c.FormFile("banner")

	if req.Status != string(model.Draft) {
		if banner == nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
				fiber.StatusBadRequest,
				"failed",
				"validation-error",
				"banner is required",
			))
		}

		// Check pdf size
		if banner.Size > (5 * 1024 * 1024) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
				fiber.StatusBadRequest,
				"failed",
				"error",
				"file is too big",
			))
		}
	}

	if err := h.EventService.CreateEvent(banner, &req); err != nil {
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
		"new event created",
	))
}

func (h *EventHandler) GetEvents(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	filterReq := new(dto.EventFilterReq)
	if err := c.QueryParser(filterReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	if filterReq.Status != "" && role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
			fiber.StatusUnauthorized,
			"failed",
			"error",
			"unauthorized action",
		))
	}

	var start, end time.Time
	if filterReq.StartTime != 0 {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		unix := time.Unix(filterReq.StartTime, 0).In(loc)
		start = time.Date(unix.Year(), unix.Month(), unix.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	}

	var title string
	if filterReq.Title != "" {
		title = utils.CreateSlug(filterReq.Title)
	}

	filter := &dto.EventFilter{
		Title:    helper.StrPointerIfNotEmpty(title),
		Category: filterReq.Category,
		Status:   helper.StrPointerIfNotEmpty(filterReq.Status),
		Start:    helper.TimeToPointer(start),
		End:      helper.TimeToPointer(end),
		Location: helper.StrPointerIfNotEmpty(filterReq.Location),
	}

	pagination := new(model.Pagination)
	if err := c.QueryParser(pagination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Sort:  pagination.Sort,
		Limit: pagination.Limit,
		Page:  pagination.Page,
	}

	events, err := h.EventService.GetEvents(filter, *page)
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
		events,
	))
}

func (h *EventHandler) GetEventByID(c *fiber.Ctx) error {
	eventID := c.Params("id")

	userLat, err := strconv.ParseFloat(c.Query("lat", "0"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	userLon, err := strconv.ParseFloat(c.Query("lon", "0"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	event, err := h.EventService.GetEventByID(eventID, userLat, userLon)
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
		event,
	))
}

func (h *EventHandler) GetEventRoute(c *fiber.Ctx) error {
	eventID := c.Params("id")

	userLat, err := strconv.ParseFloat(c.Query("lat", "0"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	userLon, err := strconv.ParseFloat(c.Query("lon", "0"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	resp, err := h.EventService.GetEventRoute(eventID, userLat, userLon)
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
		resp,
	))
}

func (h *EventHandler) FindNearestEvents(c *fiber.Ctx) error {
	userLat, err := strconv.ParseFloat(c.Query("lat", "0"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	userLon, err := strconv.ParseFloat(c.Query("lon", "0"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	user := model.Location{
		Name: "user",
		Lat:  userLat,
		Lon:  userLon,
	}

	pagination := new(model.Pagination)
	if err := c.QueryParser(pagination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Sort:  pagination.Sort,
		Limit: pagination.Limit,
		Page:  pagination.Page,
	}

	resp, err := h.EventService.FindNearestEvent(user, *page)
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
		resp,
	))
}

func (h *EventHandler) GetEventsByProfileID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	filterReq := new(dto.EventFilterReq)
	if err := c.QueryParser(filterReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	var start, end time.Time
	if filterReq.StartTime != 0 {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		unix := time.Unix(filterReq.StartTime, 0).In(loc)
		start = time.Date(unix.Year(), unix.Month(), unix.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	}

	var title string
	if filterReq.Title != "" {
		title = utils.CreateSlug(filterReq.Title)
	}

	filter := &dto.EventFilter{
		Title:    helper.StrPointerIfNotEmpty(title),
		Category: filterReq.Category,
		Status:   helper.StrPointerIfNotEmpty(filterReq.Status),
		Start:    helper.TimeToPointer(start),
		End:      helper.TimeToPointer(end),
		Location: helper.StrPointerIfNotEmpty(filterReq.Location),
	}

	pagination := new(model.Pagination)
	if err := c.QueryParser(pagination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Sort:  pagination.Sort,
		Limit: pagination.Limit,
		Page:  pagination.Page,
	}

	events, err := h.EventService.GetEventsByProfileID(userID, filter, *page)
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
		events,
	))
}

func (h *EventHandler) ReviewEvent(c *fiber.Ctx) error {
	var req *dto.ReviewEventReq
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

	if err := h.EventService.ReviewEvent(req); err != nil {
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
		"success review event",
	))
}

func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	userID := c.Locals("user_id").(string)

	data := c.FormValue("data")

	var req *dto.EventReq
	if err := json.Unmarshal([]byte(data), &req); err != nil {
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

	req.ID = &eventID
	req.UserID = userID

	banner, _ := c.FormFile("banner")
	if banner != nil {
		// Check pdf size
		if banner.Size > (5 * 1024 * 1024) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
				fiber.StatusBadRequest,
				"failed",
				"error",
				"file is too big",
			))
		}
	} else {
		banner = nil
	}

	if err := h.EventService.UpdateEvent(banner, req); err != nil {
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
		"event updated",
	))

}

func (h *EventHandler) CancelEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	if err := h.EventService.CancelEvent(id, userID); err != nil {
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
		"event canceled",
	))
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	if err := h.EventService.DeleteEvent(id, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"Event deleted",
	))
}
