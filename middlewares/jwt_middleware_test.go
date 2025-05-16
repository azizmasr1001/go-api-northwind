package middlewares

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http/httptest"
)

const jwtSecret = "test-secret"

func generateToken(userID int, firstName string, expMinutes int) string {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"firstName": firstName,
		"exp":       time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(jwtSecret))
	return signed
}

func setupRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Pastikan Redis aktif saat test
		DB:   1,                // Gunakan DB test
	})
	_ = rdb.FlushDB(context.Background()).Err()
	return rdb
}

func TestProtectedWithRedis_ValidToken(t *testing.T) {
	app := fiber.New()
	rdb := setupRedisClient()
	token := generateToken(1, "TestUser", 15)

	app.Use(ProtectedWithRedis(jwtSecret, rdb))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestProtectedWithRedis_InvalidToken(t *testing.T) {
	app := fiber.New()
	rdb := setupRedisClient()

	app.Use(ProtectedWithRedis(jwtSecret, rdb))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestProtectedWithRedis_BlacklistedToken(t *testing.T) {
	app := fiber.New()
	rdb := setupRedisClient()

	token := generateToken(2, "BlockedUser", 10)
	_ = rdb.Set(context.Background(), token, "blacklisted", 10*time.Minute).Err()

	app.Use(ProtectedWithRedis(jwtSecret, rdb))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
