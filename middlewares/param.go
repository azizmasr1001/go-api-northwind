package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ValidateIDParam(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params(paramName)
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"code":    400,
				"message": "Invalid ID parameter",
				"errors": []fiber.Map{
					{"field": paramName, "message": "must be a valid positive number"},
				},
			})
		}

		// Simpan ID hasil parsing agar bisa dipakai di controller
		c.Locals("id", id)

		return c.Next()
	}
}
