package helper

import (
	"ngevent/internal/dto"

	"github.com/gofiber/fiber/v2"
)

func IsAdmin(c *fiber.Ctx, role string) error {
	if role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
			fiber.StatusUnauthorized,
			"failed",
			"error",
			"unauthorized action",
		))
	}
	return nil
}
