package utils

import "github.com/gofiber/fiber/v2"

type Meta struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
	Total int `json:"total,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type StandardErrorResponse struct {
	Status  string        `json:"status"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

type StandardListResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

func SuccessResponse(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"status":  "success",
		"code":    statusCode,
		"message": message,
		"data":    data,
	})
}

func ListResponse(ctx *fiber.Ctx, statusCode int, message string, data interface{}, meta Meta) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"status":  "success",
		"code":    statusCode,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

func ErrorResponse(ctx *fiber.Ctx, statusCode int, message string, errors []ErrorDetail) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"code":    statusCode,
		"message": message,
		"errors":  errors,
	})
}
