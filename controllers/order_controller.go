package controllers

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/services"
	"github.com/aronipurwanto/go-api-northwind/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	service  services.OrderService
	validate *validator.Validate
}

func NewOrderController(service services.OrderService) *OrderController {
	return &OrderController{
		service:  service,
		validate: validator.New(),
	}
}

type OrderInput struct {
	models.Order
	OrderDetails []models.OrderDetail `json:"order_details" validate:"required,min=1,dive"`
}

// CreateOrder godoc
// @Summary Create new order
// @Description Create order and order details in one transaction
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body OrderInput true "Order with details"
// @Success 201 {object} models.Order
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /orders [post]
// @Security BearerAuth
func (c *OrderController) Create(ctx *fiber.Ctx) error {
	var input OrderInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if err := c.validate.Struct(input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", utils.FormatValidationErrors(err.(validator.ValidationErrors)))
	}

	created, err := c.service.Create(&input.Order, input.OrderDetails)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Failed to create order", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 201, "Order created successfully", created)
}

// GetAllOrders godoc
// @Summary Get list of orders
// @Tags Orders
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.StandardListResponse
// @Failure 500 {object} utils.StandardErrorResponse
// @Router /orders [get]
// @Security BearerAuth
func (c *OrderController) GetAll(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	data, total, err := c.service.GetAll(page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Failed to fetch orders", []utils.ErrorDetail{{Message: err.Error()}})
	}

	meta := utils.Meta{Page: page, Limit: limit, Total: int(total)}
	return utils.ListResponse(ctx, 200, "Orders retrieved", data, meta)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} fiber.Map
// @Failure 400,404 {object} utils.StandardErrorResponse
// @Router /orders/{id} [get]
// @Security BearerAuth
func (c *OrderController) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	order, details, err := c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Order not found", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Order found", fiber.Map{
		"order":   order,
		"details": details,
	})
}

// DeleteOrder godoc
// @Summary Delete order by ID
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} fiber.Map
// @Failure 400,404,500 {object} utils.StandardErrorResponse
// @Router /orders/{id} [delete]
// @Security BearerAuth
func (c *OrderController) Delete(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if err = c.service.Delete(id); err != nil {
		return utils.ErrorResponse(ctx, 500, "Delete failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Order deleted", nil)
}

// UpdateOrder godoc
// @Summary Update an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param order body OrderInput true "Updated order and details"
// @Success 200 {object} models.Order
// @Failure 400,404,500 {object} utils.StandardErrorResponse
// @Router /orders/{id} [put]
// @Security BearerAuth
func (c *OrderController) Update(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	var input OrderInput
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}
	if err := c.validate.Struct(input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", utils.FormatValidationErrors(err.(validator.ValidationErrors)))
	}

	updated, err := c.service.Update(id, &input.Order, input.OrderDetails)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Update failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Order updated", updated)
}
