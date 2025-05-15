package controllers

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/services"
	"github.com/aronipurwanto/go-api-northwind/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	service  services.CategoryService
	validate *validator.Validate
}

func NewCategoryController(service services.CategoryService) *CategoryController {
	return &CategoryController{
		service:  service,
		validate: validator.New(),
	}
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve a paginated list of product categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.StandardListResponse
// @Failure 500 {object} utils.StandardErrorResponse
// @Router /categories [get]
// @Security BearerAuth
func (c *CategoryController) GetAll(ctx *fiber.Ctx) error {
	data, err := c.service.GetAll()
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Failed to fetch categories", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.ListResponse(ctx, 200, "Categories fetched successfully", data, utils.Meta{})
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Retrieve a single category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400,404 {object} utils.StandardErrorResponse
// @Router /categories/{id} [get]
// @Security BearerAuth
func (c *CategoryController) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	data, err := c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Category not found", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.SuccessResponse(ctx, 200, "Category fetched", data)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Add a new category to the database
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object"
// @Success 201 {object} models.Category
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /categories [post]
// @Security BearerAuth
func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	cat, validationErrs, err := utils.BindAndValidate[models.Category](ctx, c.validate)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if validationErrs != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", validationErrs)
	}
	created, err := c.service.Create(cat)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Create failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 201, "Category created", created)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update category data by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Updated category object"
// @Success 200 {object} models.Category
// @Failure 400,404,500 {object} utils.StandardErrorResponse
// @Router /categories/{id} [put]
// @Security BearerAuth
func (c *CategoryController) Update(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	_, err = c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Category not found", []utils.ErrorDetail{{Message: err.Error()}})
	}

	// Parse + Validate
	input, validationErrs, err := utils.BindAndValidate[models.Category](ctx, c.validate)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}
	if validationErrs != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", validationErrs)
	}

	// Force ID from URL to match body
	input.CategoryID = id
	updated, err := c.service.Update(input)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Update failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Category updated", updated)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} fiber.Map
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /categories/{id} [delete]
// @Security BearerAuth
func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if err = c.service.Delete(id); err != nil {
		return utils.ErrorResponse(ctx, 500, "Delete failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Category deleted", nil)
}
