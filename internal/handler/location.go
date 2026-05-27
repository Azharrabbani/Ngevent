package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/service"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	LocationService *service.LocationService
}

func NewLocationHandler(
	locatonService *service.LocationService,
) *LocationHandler {
	return &LocationHandler{
		LocationService: locatonService,
	}
}

func (h *LocationHandler) SearchLocation(c *fiber.Ctx) error {
	req := new(dto.SearchReq)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	locations, err := h.LocationService.SearchLocation(req.Query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		locations,
	))
}
