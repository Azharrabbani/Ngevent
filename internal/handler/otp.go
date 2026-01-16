package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/service"
	"ngevent/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OTPHandler struct {
	OtpService *service.OTPService
	Validate   *validator.Validate
}

func NewOTPHandler(otpService *service.OTPService, validate *validator.Validate) *OTPHandler {
	return &OTPHandler{
		OtpService: otpService,
		Validate:   validate,
	}
}

func (h *OTPHandler) ResendOTP(c *fiber.Ctx) error {
	var req *dto.ResentOTPInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
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

	status, err := h.OtpService.ResendOTPCode(req.Email)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"resend-failed",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusAccepted).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"success",
		"OTP have been resend to your email",
	))
}
