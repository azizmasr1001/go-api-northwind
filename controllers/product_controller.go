package controllers

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/services"
	"github.com/aronipurwanto/go-api-northwind/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	service  services.ProductService
	validate *validator.Validate
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{
		service:  service,
		validate: validator.New(),
	}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve a paginated list of products
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.StandardListResponse
// @Failure 500 {object} utils.StandardErrorResponse
// @Router /products [get]
// @Security BearerAuth
func (c *ProductController) GetAll(ctx *fiber.Ctx) error {
	//page, limit := utils.GetPagination(ctx)
	page := ctx.Locals("page").(int)
	limit := ctx.Locals("limit").(int)

	products, total, err := c.service.GetAllPaginated(page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Failed to retrieve products", []utils.ErrorDetail{{Message: err.Error()}})
	}

	meta := utils.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}

	return utils.ListResponse(ctx, 200, "Products retrieved", products, meta)
}

// SearchProducts godoc
// @Summary Search products by name
// @Tags Products
// @Produce json
// @Param name query string true "Search keyword"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} utils.StandardListResponse
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /products/search [get]
// @Security BearerAuth
func (c *ProductController) Search(ctx *fiber.Ctx) error {
	name := ctx.Query("name", "")
	if name == "" {
		return utils.ErrorResponse(ctx, 400, "Missing query parameter: name", nil)
	}

	//page := ctx.QueryInt("page", 1)
	//limit := ctx.QueryInt("limit", 10)
	page, limit := utils.GetPagination(ctx)

	products, total, err := c.service.SearchByName(name, page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Search failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	meta := utils.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}

	return utils.ListResponse(ctx, 200, "Search results", products, meta)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieve a single product by its ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400,404 {object} utils.StandardErrorResponse
// @Router /products/{id} [get]
// @Security BearerAuth
func (c *ProductController) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	data, err := c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Product not found", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.SuccessResponse(ctx, 200, "Product retrieved", data)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the database
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} models.Product
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /products [post]
// @Security BearerAuth
func (c *ProductController) Create(ctx *fiber.Ctx) error {
	prod, validationErrs, err := utils.BindAndValidate[models.Product](ctx, c.validate)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if validationErrs != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", validationErrs)
	}

	created, err := c.service.Create(prod)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Create failed", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.SuccessResponse(ctx, 201, "Product created", created)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update product data by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Updated product object"
// @Success 200 {object} models.Product
// @Failure 400,404,500 {object} utils.StandardErrorResponse
// @Router /products/{id} [put]
// @Security BearerAuth
func (c *ProductController) Update(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	_, err = c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Product not found", []utils.ErrorDetail{{Message: err.Error()}})
	}

	prod, validationErrs, err := utils.BindAndValidate[models.Product](ctx, c.validate)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if validationErrs != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", validationErrs)
	}

	updated, err := c.service.Update(prod)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Update failed", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.SuccessResponse(ctx, 200, "Product updated", updated)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Remove product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} fiber.Map
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /products/{id} [delete]
// @Security BearerAuth
func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals("id").(int)

	if err := c.service.Delete(id); err != nil {
		return utils.ErrorResponse(ctx, 500, "Delete failed", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.SuccessResponse(ctx, 200, "Product deleted", nil)
}
