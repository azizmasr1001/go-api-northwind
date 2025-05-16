package middlewares

import (
	redisclient "github.com/aronipurwanto/go-api-northwind/internal/redis"
	"github.com/aronipurwanto/go-api-northwind/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"strings"
)

func Protected(jwtSecret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ContextKey: "user", // this is required to later retrieve token
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"code":    401,
				"message": "Unauthorized",
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			c.Locals("user_id", claims["user_id"])
			c.Locals("first_name", claims["firstName"])
			return c.Next()
		},
	})
}

func ProtectedWithRedis(jwtSecret string, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return utils.ErrorResponse(c, 401, "Missing or invalid token", nil)
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		blacklisted, err := redisclient.IsTokenBlacklisted(c.Context(), rdb, tokenStr)
		if err != nil {
			return utils.ErrorResponse(c, 500, "Redis error", []utils.ErrorDetail{{Message: err.Error()}})
		}
		if blacklisted {
			return utils.ErrorResponse(c, 401, "Token has been revoked", nil)
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.ErrorResponse(c, 401, "Invalid token", nil)
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("first_name", claims["firstName"])
		return c.Next()
	}
}
