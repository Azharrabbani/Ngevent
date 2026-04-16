package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrganizerProfileHandler struct {
	OrganizerService *service.OrganizerProfileService
	Validate         *validator.Validate
}

func NewOrganizerProfileHandler(
	organizerService *service.OrganizerProfileService,
	validate *validator.Validate,
) *OrganizerProfileHandler {
	return &OrganizerProfileHandler{
		OrganizerService: organizerService,
		Validate:         validate,
	}
}

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
		// Check pdf size
		if photo.Size > (5 * 1024 * 1024) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
				fiber.StatusBadRequest,
				"failed",
				"error",
				"file is too big",
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
			NPWPFile:    *npwpFile,
			NIB:         nibNumber,
			NIBFile:     *nibFile,
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

	return c.Status(fiber.StatusFound).JSON(dto.Success(
		fiber.StatusFound,
		"success",
		"success",
		profile,
	))
}

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

	return c.Status(fiber.StatusFound).JSON(dto.Success(
		fiber.StatusFound,
		"success",
		"success",
		profile,
	))
}

func (h *OrganizerProfileHandler) GetAllProfile(c *fiber.Ctx) error {
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

	organizers, err := h.OrganizerService.FindAll(*page)
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

func (h *OrganizerProfileHandler) FilterProfile(c *fiber.Ctx) error {
	filter := new(dto.FilterReq)
	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	if err := h.Validate.Struct(filter); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"error",
			"validation-error",
			msg,
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
		Limit: paginate.Limit,
		Page:  paginate.Page,
		Sort:  paginate.Sort,
	}

	profiles, err := h.OrganizerService.FindByCountry(filter.Country, *page)
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
		profiles,
	))
}

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

func (h *OrganizerProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	npwpFile, err := c.FormFile("npwp_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	nibFile, err := c.FormFile("nib_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	// Check if npwp or nib null
	if nibFile == nil || npwpFile == nil {
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
			NPWPFile:    *npwpFile,
			NIBFile:     *nibFile,
		},
	}

	status, err := h.OrganizerService.UpdateProfile(userID, req)
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
		"update success",
	))
}

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
