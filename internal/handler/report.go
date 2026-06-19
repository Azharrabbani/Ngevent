package handler

import (
	"fmt"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	ReportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		ReportService: reportService,
	}
}

func (h *ReportHandler) DownloadPDF(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	filter := new(dto.EventReportFilter)
	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest, "failed", "error", err.Error(),
		))
	}

	if filter.Period != "monthly" && filter.Period != "yearly" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest, "failed", "validation-error",
			"period must be 'monthly' or 'yearly'",
		))
	}

	if filter.Year == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest, "failed", "validation-error",
			"year is required",
		))
	}

	pdf, filename, err := h.ReportService.GenerateEventReportPDF(userID, filter)
	if err != nil {
		log.Printf("[DownloadPDF] error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Error(
			fiber.StatusInternalServerError, "error", "internal server error", err.Error(),
		))
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(pdf)
}
