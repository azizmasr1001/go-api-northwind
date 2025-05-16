package controllers_test

import (
	"encoding/json"
	"fmt"
	"github.com/azizmasr1001/go-api-northwind/config"
	"github.com/azizmasr1001/go-api-northwind/controllers"
	"github.com/azizmasr1001/go-api-northwind/internal/redis"
	"github.com/azizmasr1001/go-api-northwind/middlewares"
	"github.com/azizmasr1001/go-api-northwind/models"
	"github.com/azizmasr1001/go-api-northwind/repositories"
	"github.com/azizmasr1001/go-api-northwind/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setupAuthTestApp(t *testing.T) (*fiber.App, string) {
	cfg := config.LoadConfig()
	db, _ := gorm.Open(sqlserver.Open(cfg.DBSource), &gorm.Config{})
	rdb := redis.NewRedisClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPass)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)
	password := "secret123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := db.Create(&models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         "admin",
		IsActive:     true,
	}).Error
	assert.NoError(t, err)

	repo := repositories.NewAuthRepository(db)
	service := services.NewAuthService(repo, cfg.JWTSecret, rdb)
	controller := controllers.NewAuthController(service)

	app := fiber.New()
	app.Post("/login", controller.Login)
	app.Post("/refresh", controller.Refresh)
	app.Post("/logout", middlewares.ProtectedWithRedis(cfg.JWTSecret, rdb), controller.Logout)
	app.Get("/me", middlewares.ProtectedWithRedis(cfg.JWTSecret, rdb), controller.Me)
	app.Post("/register", controller.Register)

	loginPayload := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	return app, loginPayload
}

func readBody(t *testing.T, resp *http.Response) string {
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	_ = resp.Body.Close()
	return string(b)
}

func TestRegister_Success(t *testing.T) {
	app, _ := setupAuthTestApp(t)

	payload := fmt.Sprintf(`{
		"username": "user_%d",
		"email": "user_%d@example.com",
		"password": "secret123"
	}`, time.Now().UnixNano(), time.Now().UnixNano())

	req := httptest.NewRequest("POST", "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestRegister_Conflict(t *testing.T) {
	app, _ := setupAuthTestApp(t)
	username := fmt.Sprintf("dupe_%d", time.Now().UnixNano())
	email := fmt.Sprintf("dupe_%d@example.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"username":"%s","email":"%s","password":"test123"}`, username, email)

	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, 201, resp.StatusCode)

	req2 := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2, -1)
	assert.Equal(t, 409, resp2.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	app, loginPayload := setupAuthTestApp(t)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLogin_Invalid(t *testing.T) {
	app, _ := setupAuthTestApp(t)
	payload := `{"username":"invaliduser", "password":"wrongpass"}`

	req := httptest.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestRefresh_Success(t *testing.T) {
	app, loginPayload := setupAuthTestApp(t)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	_ = resp.Body.Close()

	token := result["data"].(map[string]interface{})["token"].(string)
	body := fmt.Sprintf(`{"refresh_token":"%s"}`, token)

	reqRefresh := httptest.NewRequest("POST", "/refresh", strings.NewReader(body))
	reqRefresh.Header.Set("Content-Type", "application/json")
	respRefresh, _ := app.Test(reqRefresh, -1)

	assert.Equal(t, 200, respRefresh.StatusCode)
}

func TestLogout_Success(t *testing.T) {
	app, loginPayload := setupAuthTestApp(t)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	_ = resp.Body.Close()

	token := result["data"].(map[string]interface{})["token"].(string)

	reqLogout := httptest.NewRequest("POST", "/logout", nil)
	reqLogout.Header.Set("Authorization", "Bearer "+token)
	respLogout, _ := app.Test(reqLogout, -1)

	assert.Equal(t, 200, respLogout.StatusCode)
}

func TestMe_Success(t *testing.T) {
	app, loginPayload := setupAuthTestApp(t)

	req := httptest.NewRequest("POST", "/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	_ = resp.Body.Close()

	token := result["data"].(map[string]interface{})["token"].(string)

	reqMe := httptest.NewRequest("GET", "/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+token)
	respMe, _ := app.Test(reqMe, -1)

	body := readBody(t, respMe)
	if respMe.StatusCode != 200 {
		t.Logf("Me failed: %s", body)
	}
	assert.Equal(t, 200, respMe.StatusCode)
}
