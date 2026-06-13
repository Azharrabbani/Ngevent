package handler

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/service"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoriesHandler struct {
	CategoriesService *service.CategoryService
	Validate          *validator.Validate
}

func NewCategoriesHandler(
	categoriesService *service.CategoryService,
	validate *validator.Validate,
) *CategoriesHandler {
	return &CategoriesHandler{
		CategoriesService: categoriesService,
		Validate:          validate,
	}
}

// CreateCategory godoc
// @Summary      Create a category
// @Description  Creates a new event category (Admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        body  body     dto.CreateCatReq  true  "Category name"
// @Success      201   {object} dto.Response{data=string}
// @Failure      400   {object} dto.Response
// @Failure      401   {object} dto.Response
// @Router       /category/ [post]
func (h *CategoriesHandler) CreateCategory(c *fiber.Ctx) error {
	var req *dto.CreateCatReq
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
			"error-validation",
			msg,
		))
	}

	if err := h.CategoriesService.Create(req); err != nil {
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
		"category created",
	))
}

// ListCategories godoc
// @Summary      List all categories
// @Description  Returns all event categories (public)
// @Tags         categories
// @Produce      json
// @Success      200  {object} dto.Response{data=[]dto.ListCatResp}
// @Failure      400  {object} dto.Response
// @Router       /category/ [get]
func (h *CategoriesHandler) ListCategories(c *fiber.Ctx) error {
	categories, err := h.CategoriesService.FindAll()
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
		categories,
	))
}

// ListWithPagination godoc
// @Summary      List categories with pagination
// @Description  Returns a paginated list of categories (Admin only)
// @Tags         categories
// @Produce      json
// @Security     CookieAuth
// @Param        name   query    string  false  "Filter by category name"
// @Param        page   query    int     false  "Page number"
// @Param        limit  query    int     false  "Items per page"
// @Param        sort   query    string  false  "Sort order"
// @Success      200    {object} dto.Response{data=model.PaginationRow}
// @Failure      400    {object} dto.Response
// @Failure      401    {object} dto.Response
// @Router       /category/list [get]
func (h *CategoriesHandler) ListWithPagination(c *fiber.Ctx) error {
	filter := new(dto.FilterCatReq)

	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	if filter.Name != nil {
		name := utils.CreateSlug(helper.StringValue(filter.Name))
		filter.Name = &name
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

	res, err := h.CategoriesService.GetWithPagination(*page, filter)
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
		res,
	))
}

// ListByCatName godoc
// @Summary      Filter categories by name
// @Description  Returns categories matching the given name (public)
// @Tags         categories
// @Produce      json
// @Param        name   query    string  false  "Category name to search"
// @Param        page   query    int     false  "Page number"
// @Param        limit  query    int     false  "Items per page"
// @Success      200    {object} dto.Response{data=[]dto.ListCatResp}
// @Failure      400    {object} dto.Response
// @Router       /category/filter [get]
func (h *CategoriesHandler) ListByCatName(c *fiber.Ctx) error {
	catName := new(dto.FindCatReq)

	if err := c.QueryParser(catName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
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

	categories, err := h.CategoriesService.FindBySlug(catName.Name, *page)
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
		categories,
	))
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Updates an existing event category by ID (Admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        id    path     string            true  "Category ID"
// @Param        body  body     dto.UpdateCatReq  true  "Updated category name"
// @Success      200   {object} dto.Response{data=string}
// @Failure      400   {object} dto.Response
// @Failure      401   {object} dto.Response
// @Router       /category/{id} [put]
func (h *CategoriesHandler) UpdateCategory(c *fiber.Ctx) error {
	catID := c.Params("id")

	if catID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			"missing category ID",
		))
	}

	var req *dto.UpdateCatReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error",
			err.Error(),
		))
	}

	req.CategoryID = catID

	if err := h.Validate.Struct(req); err != nil {
		msg := utils.GetValidationError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.Error(
			fiber.StatusBadRequest,
			"failed",
			"error-validation",
			msg,
		))
	}

	if err := h.CategoriesService.Update(req); err != nil {
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
		"category updated",
	))
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Deletes an event category by ID (Admin only)
// @Tags         categories
// @Produce      json
// @Security     CookieAuth
// @Param        id  path     string  true  "Category ID"
// @Success      200 {object} dto.Response{data=string}
// @Failure      400 {object} dto.Response
// @Failure      401 {object} dto.Response
// @Router       /category/{id} [delete]
func (h *CategoriesHandler) DeleteCategory(c *fiber.Ctx) error {
	catID := c.Params("id")

	if err := h.CategoriesService.Delete(catID); err != nil {
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
		"category deleted",
	))
}
