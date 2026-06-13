package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		validate:    validate,
	}
}

// ListPhoneCodes godoc
// @Summary      List international phone codes
// @Description  Returns a paginated list of international phone dial codes
// @Tags         auth
// @Produce      json
// @Param        page   query    int  false  "Page number (default 1)"
// @Param        limit  query    int  false  "Items per page (default 20)"
// @Success      200    {object} dto.Response{data=[]dto.PhoneCodeResp}
// @Router       /phone-codes [get]
func (h *AuthHandler) ListPhoneCodes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	phoneCodes := utils.ListAllPhoneCodes(page, limit)

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		phoneCodes,
	))
}

// VerififyEmail godoc
// @Summary      Verify email address
// @Description  Verifies a user's email using the OTP sent during registration
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body     dto.VerifyEmailInput  true  "Email and OTP"
// @Success      200   {object} dto.Response{data=string}
// @Failure      400   {object} dto.Response
// @Router       /verify-email [put]
func (h *AuthHandler) VerififyEmail(c *fiber.Ctx) error {
	var req dto.VerifyEmailInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"validation-error",
			msg,
		))
	}

	// Verify Email
	status, err := h.AuthService.VerififyEmail(req.Email, req.OTP)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"verify-error",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"success verified email",
	))

}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns JWT tokens via cookies and response body
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body     dto.LoginInput  true  "Login credentials"
// @Success      200   {object} dto.Response{data=dto.LoginResponse}
// @Failure      400   {object} dto.Response
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Validate the req
	var req dto.LoginInput

	// Get user client
	ip := c.IP()
	userAgent := utils.Handler(c)

	client := &model.Client{
		IP:        ip,
		UserAgent: userAgent,
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
		))
	}

	// Login
	user, accessToken, refreshToken, refreshExp, err := h.AuthService.Login(client, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"login-failed",
			err.Error(),
		))
	}

	// Set Access Token cookie (short lived)
	c.Cookie(&fiber.Cookie{
		Name:     "ngevent_cookie",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   60 * 60 * 3,
		SameSite: "None",
	})

	// Set Refresh Token Cookie (long lived)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		Expires:  refreshExp,
		SameSite: "None",
	})

	loginUser := &dto.LoginResponse{
		ID:              user.ID,
		Email:           user.Email,
		Role:            user.Role,
		NgeventToken:    accessToken,
		NgeventRefToken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"login-success",
		loginUser,
	))
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Issues a new access token using the refresh_token cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object} dto.Response{data=dto.RefreshTokenResp}
// @Failure      401  {object} dto.Response
// @Router       /refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
			fiber.StatusUnauthorized,
			"failed",
			"error",
			"missing refresh token",
		))
	}

	accessToken, err := h.AuthService.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
			fiber.StatusUnauthorized,
			"failed",
			"error",
			err.Error(),
		))
	}

	// Set new access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "ngevent_cookie",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   60 * 60 * 3,
		SameSite: "None",
	})

	token := &dto.RefreshTokenResp{
		NgeventToken: accessToken,
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"token-refreshed",
		token,
	))
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Sends a password reset link to the provided email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body     dto.ForgetPasswordInput  true  "Email address"
// @Success      202   {object} dto.Response{data=string}
// @Failure      400   {object} dto.Response
// @Router       /forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgetPasswordInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
		))
	}

	// Forgot password
	status, err := h.AuthService.ForgotPassword(req.Email)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"error",
			"invalid-request",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusAccepted).JSON(dto.Success(
		fiber.StatusAccepted,
		"success",
		"success",
		"Reset link have been send to your email",
	))
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets a user's password using the OTP token from the reset link
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        id    path     string                 true  "OTP token ID"
// @Param        body  body     dto.ResetPasswordInput true  "New password"
// @Success      200   {object} dto.Response{data=string}
// @Failure      400   {object} dto.Response
// @Router       /reset-password/{id} [put]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	otpID := c.Params("id")

	var req dto.ResetPasswordInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	if err := h.validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
		))
	}

	// Reset password
	status, err := h.AuthService.ResetPassword(otpID, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"invalid-request",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"password have been reset",
	))
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidates the user's session and clears auth cookies
// @Tags         auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object} dto.Response
// @Failure      400  {object} dto.Response
// @Router       /logout/ [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			"missing refresh token",
		))
	}

	err := h.AuthService.Logout(refreshToken)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	helper.ClearAuthCookies(c)

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		nil,
	))
}
