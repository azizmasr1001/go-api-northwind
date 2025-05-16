package controllers

import (
	"github.com/azizmasr1001/go-api-northwind/models"
	"github.com/azizmasr1001/go-api-northwind/services"
	"github.com/azizmasr1001/go-api-northwind/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EmployeeController struct {
	service  services.EmployeeService
	validate *validator.Validate
}

func NewEmployeeController(service services.EmployeeService) *EmployeeController {
	return &EmployeeController{
		service:  service,
		validate: validator.New(),
	}
}

// GetAllEmployees godoc
// @Summary Get all employees
// @Description Retrieve a paginated list of employees
// @Tags Employees
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} utils.StandardListResponse
// @Failure 500 {object} utils.StandardErrorResponse
// @Router /employees [get]
// @Security BearerAuth
func (c *EmployeeController) GetAll(ctx *fiber.Ctx) error {
	page, limit := utils.GetPagination(ctx)

	data, total, err := c.service.GetAllPaginated(page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Failed to retrieve data", []utils.ErrorDetail{{Message: err.Error()}})
	}

	meta := utils.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}

	return utils.ListResponse(ctx, 200, "List retrieved successfully", data, meta)
}

// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Description Retrieve a single employee by its ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} models.Employee
// @Failure 400,404 {object} utils.StandardErrorResponse
// @Router /employees/{id} [get]
// @Security BearerAuth
func (c *EmployeeController) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	emp, err := c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Employee not found", []utils.ErrorDetail{{Message: err.Error()}})
	}
	return utils.SuccessResponse(ctx, 200, "Data retrieved successfully", emp)
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Add a new employee to the system
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee object"
// @Success 201 {object} models.Employee
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /employees [post]
// @Security BearerAuth
func (c *EmployeeController) Create(ctx *fiber.Ctx) error {
	emp, validationErrs, err := utils.BindAndValidate[models.Employee](ctx, c.validate)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if validationErrs != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", validationErrs)
	}
	created, err := c.service.Create(emp)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Failed to create employee", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 201, "Employee created successfully", created)
}

// UpdateEmployee godoc
// @Summary Update an employee
// @Description Update employee information by ID
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body models.Employee true "Updated employee object"
// @Success 200 {object} models.Employee
// @Failure 400,404,500 {object} utils.StandardErrorResponse
// @Router /employees/{id} [put]
// @Security BearerAuth
func (c *EmployeeController) Update(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	_, err = c.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(ctx, 404, "Employee not found", []utils.ErrorDetail{{Message: err.Error()}})
	}

	input, validationErrs, err := utils.BindAndValidate[models.Employee](ctx, c.validate)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}
	if validationErrs != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", validationErrs)
	}

	input.EmployeeID = id
	updated, err := c.service.Update(input)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Update failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Employee updated successfully", updated)
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Remove employee record by ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} fiber.Map
// @Failure 400,500 {object} utils.StandardErrorResponse
// @Router /employees/{id} [delete]
// @Security BearerAuth
func (c *EmployeeController) Delete(ctx *fiber.Ctx) error {
	id, err := utils.ParseID(ctx)
	if err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid ID", []utils.ErrorDetail{{Message: err.Error()}})
	}

	if err = c.service.Delete(id); err != nil {
		return utils.ErrorResponse(ctx, 500, "Delete failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Employee deleted successfully", nil)
}
