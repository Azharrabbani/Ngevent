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

// SearchLocation godoc
// @Summary      Search for a location
// @Description  Searches locations using Nominatim/OpenStreetMap (used during event creation)
// @Tags         location
// @Produce      json
// @Param        query  query    string  true  "Search query string"
// @Success      200    {object} dto.Response{data=[]dto.SearchResponse}
// @Failure      400    {object} dto.Response
// @Router       /location/ [get]
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
