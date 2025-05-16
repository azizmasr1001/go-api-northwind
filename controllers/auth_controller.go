package controllers

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/services"
	"github.com/aronipurwanto/go-api-northwind/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type AuthController struct {
	service  services.AuthService
	validate *validator.Validate
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{
		service:  service,
		validate: validator.New(),
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400,401 {object} utils.StandardErrorResponse
// @Router /login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var input models.LoginRequest
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", []utils.ErrorDetail{{Message: err.Error()}})
	}
	if err := c.validate.Struct(input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", utils.FormatValidationErrors(err.(validator.ValidationErrors)))
	}

	user, err := c.service.Login(ctx.Context(), input.Username, input.Password)
	if err != nil {
		return utils.ErrorResponse(ctx, 401, "Unauthorized", []utils.ErrorDetail{{Message: err.Error()}})
	}

	accessToken, err := c.service.GenerateAccessToken(user)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Token generation failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	refreshToken, err := c.service.GenerateRefreshToken(user)
	if err != nil {
		return utils.ErrorResponse(ctx, 500, "Refresh token generation failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Login successful", fiber.Map{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

// Register godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.RegisterRequest true "User registration"
// @Success 201 {object} models.User
// @Failure 400,409,500 {object} utils.StandardErrorResponse
// @Router /register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var input models.RegisterRequest
	if err := ctx.BodyParser(&input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", nil)
	}
	if err := c.validate.Struct(input); err != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", utils.FormatValidationErrors(err.(validator.ValidationErrors)))
	}

	user, err := c.service.Register(ctx.Context(), &input)
	if err != nil {
		if strings.Contains(err.Error(), "username already taken") || strings.Contains(err.Error(), "email already taken") {
			return utils.ErrorResponse(ctx, 409, "Conflict", []utils.ErrorDetail{{Message: err.Error()}})
		}
		return utils.ErrorResponse(ctx, 500, "Registration failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 201, "User registered", user)
}

// Me godoc
// @Summary Get current user info
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} utils.StandardErrorResponse
// @Router /me [get]
// @Security BearerAuth
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")
	firstName := ctx.Locals("first_name")

	if userID == nil || firstName == nil {
		return utils.ErrorResponse(ctx, 401, "Unauthorized", nil)
	}

	return utils.SuccessResponse(ctx, 200, "User info retrieved", fiber.Map{
		"user_id":    userID,
		"first_name": firstName,
	})
}

// Logout godoc
// @Summary User logout
// @Tags Auth
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 400,401 {object} utils.StandardErrorResponse
// @Router /logout [post]
// @Security BearerAuth
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	header := ctx.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return utils.ErrorResponse(ctx, 400, "Invalid token format", nil)
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.service.GetSecret()), nil
	})
	if err != nil || !token.Valid {
		return utils.ErrorResponse(ctx, 401, "Invalid token", nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	expUnix := int64(claims["exp"].(float64))
	expDuration := time.Until(time.Unix(expUnix, 0))

	if err := c.service.Logout(ctx.Context(), tokenStr, expDuration); err != nil {
		return utils.ErrorResponse(ctx, 500, "Logout failed", []utils.ErrorDetail{{Message: err.Error()}})
	}

	return utils.SuccessResponse(ctx, 200, "Logout successful", nil)
}

// Refresh godoc
// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 400,401,500 {object} utils.StandardErrorResponse
// @Router /refresh [post]
func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, 400, "Invalid input", nil)
	}
	if err := c.validate.Struct(req); err != nil {
		return utils.ErrorResponse(ctx, 400, "Validation failed", utils.FormatValidationErrors(err.(validator.ValidationErrors)))
	}

	token, err := c.service.RefreshToken(ctx.Context(), req.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(ctx, 401, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, 200, "Token refreshed", fiber.Map{"access_token": token})
}

// SendOTP godoc
// @Summary Send OTP to email for password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body map[string]string true "Email"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} utils.StandardErrorResponse
// @Router /send-otp [post]
func (c *AuthController) SendOTP(ctx *fiber.Ctx) error {
	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil || body["email"] == "" {
		return utils.ErrorResponse(ctx, 400, "Email is required", nil)
	}

	if err := c.service.SendResetOTP(ctx.Context(), body["email"]); err != nil {
		return utils.ErrorResponse(ctx, 400, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, 200, "OTP sent to email", nil)
}

// VerifyOTP godoc
// @Summary Verify OTP for password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Email and OTP"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} utils.StandardErrorResponse
// @Router /verify-otp [post]
func (c *AuthController) VerifyOTP(ctx *fiber.Ctx) error {
	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil || body["email"] == "" || body["otp"] == "" {
		return utils.ErrorResponse(ctx, 400, "Email and OTP are required", nil)
	}

	ok, err := c.service.VerifyResetOTP(ctx.Context(), body["email"], body["otp"])
	if err != nil || !ok {
		return utils.ErrorResponse(ctx, 400, "Invalid OTP", nil)
	}

	return utils.SuccessResponse(ctx, 200, "OTP verified", nil)
}

// ResetPassword godoc
// @Summary Reset password using email and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Email and new password"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} utils.StandardErrorResponse
// @Router /reset-password [post]
func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	var body map[string]string
	if err := ctx.BodyParser(&body); err != nil || body["email"] == "" || body["new_password"] == "" {
		return utils.ErrorResponse(ctx, 400, "Email and new password are required", nil)
	}

	if err := c.service.ResetPassword(ctx.Context(), body["email"], body["new_password"]); err != nil {
		return utils.ErrorResponse(ctx, 400, err.Error(), nil)
	}

	return utils.SuccessResponse(ctx, 200, "Password reset successful", nil)
}
