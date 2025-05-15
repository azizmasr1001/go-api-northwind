package controllers

import (
	"github.com/aronipurwanto/go-api-northwind/config"
	"github.com/aronipurwanto/go-api-northwind/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthController struct {
	DB     *gorm.DB
	Config config.Config
}

func NewAuthController(db *gorm.DB, cfg config.Config) *AuthController {
	return &AuthController{DB: db, Config: cfg}
}

type LoginRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.Employee
	if err := a.DB.Where("first_name = ? AND last_name = ?", req.FirstName, req.LastName).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Create JWT Token
	claims := jwt.MapClaims{
		"user_id":   user.EmployeeID,
		"firstName": user.FirstName,
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.Config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  signedToken,
	})
}
