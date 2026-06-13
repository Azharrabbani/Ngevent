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

type OrganizerProfileHandler struct {
	OrganizerService *service.OrganizerProfileService
	AuthService      *service.AuthService
	Validate         *validator.Validate
}

func NewOrganizerProfileHandler(
	organizerService *service.OrganizerProfileService,
	authService *service.AuthService,
	validate *validator.Validate,
) *OrganizerProfileHandler {
	return &OrganizerProfileHandler{
		OrganizerService: organizerService,
		AuthService:      authService,
		Validate:         validate,
	}
}

// CreateProfile godoc
// @Summary      Create organizer profile
// @Description  Creates an event organizer profile, including company NPWP/NIB documents (Organizer only)
// @Tags         organizer
// @Accept       multipart/form-data
// @Produce      json
// @Security     CookieAuth
// @Param        photo        formData  file    false  "Profile photo (max 5MB)"
// @Param        name         formData  string  true   "Organizer/company name"
// @Param        phonenumber  formData  string  true   "Phone number"
// @Param        iso          formData  string  false  "Phone country ISO code"
// @Param        address      formData  string  true   "Address"
// @Param        email        formData  string  false  "Public contact email"
// @Param        instagram    formData  string  false  "Instagram handle"
// @Param        description  formData  string  false  "Company description"
// @Param        npwp_number  formData  string  true   "NPWP number"
// @Param        npwp_file    formData  file    true   "NPWP document (max 100MB)"
// @Param        nib_number   formData  string  true   "NIB number"
// @Param        nib_file     formData  file    true   "NIB document (max 100MB)"
// @Success      201  {object} dto.Response{data=string}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/ [post]
func (h *OrganizerProfileHandler) CreateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	// Only event organizer can create EO profile
	if role != "event organizer" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
			fiber.StatusUnauthorized,
			"failed",
			"error",
			"unauthorized action",
		))
	}

	photo, _ := c.FormFile("photo")
	if photo != nil {
		// Check photo size
		if photo.Size > (5 * 1024 * 1024) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
				fiber.StatusBadRequest,
				"failed",
				"error",
				"File is too big",
			))
		}
	} else {
		photo = nil
	}

	npwpFile, _ := c.FormFile("npwp_file")

	nibFile, _ := c.FormFile("nib_file")

	if npwpFile != nil && nibFile != nil {
		// Check pdf size
		if npwpFile.Size > (100*1024*1024) || nibFile.Size > (100*1024*1024) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
				fiber.StatusBadRequest,
				"failed",
				"error",
				"file for npwp or nib is too big",
			))
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			"NPWP and NIB must be uploaded.",
		))
	}

	name := c.FormValue("name")
	phoneNumber := c.FormValue("phonenumber")
	iso := c.FormValue("iso")
	address := c.FormValue("address")
	email := c.FormValue("email")
	instagram := c.FormValue("instagram")
	desc := c.FormValue("description")
	npwpNumber := c.FormValue("npwp_number")
	nibNumber := c.FormValue("nib_number")

	req := &dto.CreateOrganizerProfileReq{
		UserID:       userID,
		PhotoProfile: photo,
		Name:         name,
		PhoneNumber:  phoneNumber,
		ISO:          iso,
		Address:      &address,
		SocialMedia: dto.OrganizerSocialMediaReq{
			Email:     &email,
			Instagram: &instagram,
		},
		CompanyDetail: dto.OrganizerCompDetailReq{
			Description: &desc,
			NPWP:        npwpNumber,
			NPWPFile:    npwpFile,
			NIB:         nibNumber,
			NIBFile:     nibFile,
		},
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

	if err := h.OrganizerService.CreateProfile(req); err != nil {
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

// GetProfileByID godoc
// @Summary      Get organizer profile by ID
// @Description  Returns an organizer profile by its profile ID (authenticated users)
// @Tags         organizer
// @Produce      json
// @Security     CookieAuth
// @Param        id  path     string  true  "Organizer profile ID"
// @Success      200  {object} dto.Response{data=dto.OrganizerProfilesResponse}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/{id} [get]
func (h *OrganizerProfileHandler) GetProfileByID(c *fiber.Ctx) error {
	id := c.Params("id")

	profile, err := h.OrganizerService.FindByID(id)
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
		profile,
	))
}

// GetProfileBySlug godoc
// @Summary      Get organizer profile by slug
// @Description  Returns the public organizer profile matching the given slug
// @Tags         organizer
// @Produce      json
// @Param        slug  path     string  true  "Organizer profile slug"
// @Success      200  {object} dto.Response{data=dto.OrganizerProfilesResponse}
// @Failure      400  {object} dto.Response
// @Router       /organizer/public/{slug} [get]
func (h *OrganizerProfileHandler) GetProfileBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			"slug is required",
		))
	}

	profile, err := h.OrganizerService.FindBySlug(slug)
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
		profile,
	))
}

// GetProfileByUserID godoc
// @Summary      Get my organizer profile
// @Description  Returns the authenticated user's organizer profile
// @Tags         organizer
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object} dto.Response{data=dto.OrganizerProfilesResponse}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/me [get]
func (h *OrganizerProfileHandler) GetProfileByUserID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	profile, err := h.OrganizerService.FindByUserID(userID)
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
		profile,
	))
}

// GetAllProfile godoc
// @Summary      List all organizer profiles
// @Description  Returns a paginated list of organizer profiles, including pending/rejected ones (Admin only)
// @Tags         organizer
// @Produce      json
// @Security     CookieAuth
// @Param        filter  query    string  false  "Filter by name/email"
// @Param        status  query    string  false  "Filter by status"
// @Param        page    query    int     false  "Page number"
// @Param        limit   query    int     false  "Items per page"
// @Param        sort    query    string  false  "Sort order"
// @Success      200  {object} dto.Response{data=model.PaginationRow}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/profiles [get]
func (h *OrganizerProfileHandler) GetAllProfile(c *fiber.Ctx) error {
	filter := new(dto.FilterProfileReq)

	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"error",
			err.Error(),
		))
	}

	paginate := new(model.Pagination)

	if err := c.QueryParser(paginate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Page:  paginate.Page,
		Limit: paginate.Limit,
		Sort:  paginate.Sort,
	}

	organizers, err := h.OrganizerService.FindAll(*page, filter)
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
		organizers,
	))
}

// GetAllProfileForPublic godoc
// @Summary      List organizer profiles (public)
// @Description  Returns a paginated list of verified organizer profiles for public listing
// @Tags         organizer
// @Produce      json
// @Param        filter  query    string  false  "Filter by name"
// @Param        page    query    int     false  "Page number"
// @Param        limit   query    int     false  "Items per page"
// @Param        sort    query    string  false  "Sort order"
// @Success      200  {object} dto.Response{data=model.PaginationRow}
// @Failure      400  {object} dto.Response
// @Router       /organizer/ [get]
func (h *OrganizerProfileHandler) GetAllProfileForPublic(c *fiber.Ctx) error {
	filter := new(dto.FilterPublicProfileReq)

	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"error",
			err.Error(),
		))
	}

	paginate := new(model.Pagination)

	if err := c.QueryParser(paginate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	page := &model.Pagination{
		Page:  paginate.Page,
		Limit: paginate.Limit,
		Sort:  paginate.Sort,
	}

	organizers, err := h.OrganizerService.FindAllForPublic(*page, filter)
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
		organizers,
	))
}

// UpdatePhotoProfile godoc
// @Summary      Update organizer profile photo
// @Description  Updates the authenticated organizer's profile photo (Organizer/Admin only)
// @Tags         organizer
// @Accept       multipart/form-data
// @Produce      json
// @Security     CookieAuth
// @Param        photo  formData  file  true  "New profile photo"
// @Success      200  {object} dto.Response{data=string}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/photo [put]
func (h *OrganizerProfileHandler) UpdatePhotoProfile(c *fiber.Ctx) error {
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

	status, err := h.OrganizerService.UpdatePhotoProfile(photo, userID)
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

// UpdateProfile godoc
// @Summary      Update organizer profile
// @Description  Updates the authenticated organizer's profile (Organizer/Admin only). Changes to critical fields (NPWP/NIB) are submitted for admin review instead of being applied immediately
// @Tags         organizer
// @Accept       multipart/form-data
// @Produce      json
// @Security     CookieAuth
// @Param        name         formData  string  true   "Organizer/company name"
// @Param        phonenumber  formData  string  true   "Phone number"
// @Param        iso          formData  string  false  "Phone country ISO code"
// @Param        address      formData  string  false  "Address"
// @Param        email        formData  string  false  "Public contact email"
// @Param        instagram    formData  string  false  "Instagram handle"
// @Param        description  formData  string  false  "Company description"
// @Param        npwp_number  formData  string  false  "New NPWP number (triggers review)"
// @Param        npwp_file    formData  file    false  "New NPWP document (triggers review)"
// @Param        nib_number   formData  string  false  "New NIB number (triggers review)"
// @Param        nib_file     formData  file    false  "New NIB document (triggers review)"
// @Success      200  {object} dto.Response{data=string}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/ [put]
func (h *OrganizerProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// OPTIONAL FILE
	npwpFile, _ := c.FormFile("npwp_file")
	nibFile, _ := c.FormFile("nib_file")

	name := c.FormValue("name")
	phoneNumber := c.FormValue("phonenumber")
	iso := c.FormValue("iso")
	address := c.FormValue("address")
	email := c.FormValue("email")
	instagram := c.FormValue("instagram")
	desc := c.FormValue("description")
	npwpNumber := c.FormValue("npwp_number")
	nibNumber := c.FormValue("nib_number")

	req := &dto.UpdateOrganizerProfileReq{
		Name:        name,
		PhoneNumber: phoneNumber,
		ISO:         iso,
		Address:     &address,
		SocialMedia: dto.OrganizerSocialMediaReq{
			Email:     &email,
			Instagram: &instagram,
		},
		CompanyDetail: dto.OrganizerCompDetailReq{
			Description: &desc,
			NPWP:        npwpNumber,
			NIB:         nibNumber,
			NPWPFile:    npwpFile,
			NIBFile:     nibFile,
		},
	}

	status, isCritical, err := h.OrganizerService.UpdateProfile(userID, req)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"error",
			err.Error(),
		))
	}

	if isCritical {
		return c.Status(fiber.StatusOK).JSON(dto.Success(
			fiber.StatusOK,
			"success",
			"success",
			"Update will be reviewed",
		))
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(
		fiber.StatusOK,
		"success",
		"success",
		"Profile updated",
	))
}

// ApprovedProfile godoc
// @Summary      Approve organizer profile
// @Description  Marks an organizer profile as verified/approved (Admin only)
// @Tags         organizer
// @Produce      json
// @Security     CookieAuth
// @Param        id  path     string  true  "Organizer profile ID"
// @Success      200  {object} dto.Response{data=string}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/approve/{id} [put]
func (h *OrganizerProfileHandler) ApprovedProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	profileID := c.Params("id")

	approvedReq := &dto.ApprovedReq{
		ReviewedBy: userID,
		ReviewedAt: time.Now().UTC(),
	}

	if err := h.OrganizerService.VerifiedProfile(profileID, approvedReq); err != nil {
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
		"approved success",
	))
}

// RejectProfile godoc
// @Summary      Reject organizer profile
// @Description  Rejects an organizer profile registration with a reason (Admin only)
// @Tags         organizer
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id    path     string            true  "Organizer profile ID"
// @Param        body  body     dto.RejectedReq   true  "Rejection reason"
// @Success      200  {object} dto.Response{data=string}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/reject/{id} [put]
func (h *OrganizerProfileHandler) RejectProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	profileID := c.Params("id")

	var req *dto.RejectedReq

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
			"error",
			msg,
		))
	}

	rejectedReq := &dto.RejectedReq{
		Reason:     req.Reason,
		ReviewedBy: userID,
		ReviewedAt: time.Now().UTC(),
	}

	if err := h.OrganizerService.RejectProfile(profileID, rejectedReq); err != nil {
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
		"rejected success",
	))
}

// CloseAccount godoc
// @Summary      Close organizer account
// @Description  Closes the authenticated organizer's account, logs out the session and clears auth cookies (Organizer only)
// @Tags         organizer
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object} dto.Response{data=string}
// @Failure      400  {object} dto.Response
// @Failure      401  {object} dto.Response
// @Router       /organizer/close-account [delete]
func (h *OrganizerProfileHandler) CloseAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			"missing refresh token",
		))
	}

	status, err := h.OrganizerService.CloseAccount(userID)
	if err != nil {
		return c.Status(status).JSON(dto.Error(
			status,
			"failed",
			"error",
			err.Error(),
		))
	}

	err = h.AuthService.Logout(refreshToken)

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
		"your account has been closed",
	))
}
