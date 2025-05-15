package utils

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ParseID parses ":id" from URL and returns it as int
func ParseID(ctx *fiber.Ctx) (int, error) {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// BindAndValidate parses body JSON and validates it with validator
func BindAndValidate[T any](ctx *fiber.Ctx, validate *validator.Validate) (*T, []ErrorDetail, error) {
	var payload T

	if err := ctx.BodyParser(&payload); err != nil {
		return nil, nil, err
	}

	if err := validate.Struct(payload); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return nil, FormatValidationErrors(errs), nil
		}
		return nil, nil, err
	}

	return &payload, nil, nil
}

// FormatValidationErrors converts validator.ValidationErrors to []ErrorDetail
func FormatValidationErrors(errs validator.ValidationErrors) []ErrorDetail {
	var details []ErrorDetail
	for _, e := range errs {
		details = append(details, ErrorDetail{
			Field:   e.Field(),
			Message: e.Error(),
		})
	}
	return details
}
