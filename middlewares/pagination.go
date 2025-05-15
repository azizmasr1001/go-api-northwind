package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ValidateQueryPagination(defaultPage, defaultLimit int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageStr := c.Query("page", strconv.Itoa(defaultPage))
		limitStr := c.Query("limit", strconv.Itoa(defaultLimit))

		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)

		if err1 != nil || err2 != nil || page <= 0 || limit <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"code":    400,
				"message": "Invalid pagination parameters",
				"errors": []fiber.Map{
					{"field": "page", "message": "must be a positive number"},
					{"field": "limit", "message": "must be a positive number"},
				},
			})
		}

		// Simpan ke context
		c.Locals("page", page)
		c.Locals("limit", limit)

		return c.Next()
	}
}
