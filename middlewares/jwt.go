package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Protected(jwtSecret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"code":    401,
				"message": "Unauthorized",
			})
		},
	})
}
