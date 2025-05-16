package main

import (
	"context"
	"github.com/aronipurwanto/go-api-northwind/config"
	"github.com/aronipurwanto/go-api-northwind/controllers"
	"github.com/aronipurwanto/go-api-northwind/repositories"
	"github.com/aronipurwanto/go-api-northwind/routes"
	"github.com/aronipurwanto/go-api-northwind/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"

	redisClient "github.com/aronipurwanto/go-api-northwind/internal/redis"
	"github.com/gofiber/swagger"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(sqlserver.Open(cfg.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	log.Println("Connected to DB")

	redis := redisClient.NewRedisClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPass)
	pong, err := redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis:", pong)

	empRepo := repositories.NewEmployeeRepository(db)
	empService := services.NewEmployeeService(empRepo)
	empController := controllers.NewEmployeeController(empService)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, cfg.JWTSecret, redis)
	authController := controllers.NewAuthController(authService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo)
	orderController := controllers.NewOrderController(orderService)

	app := fiber.New()

	// Swagger info
	// @title Northwind API
	// @version 1.0
	// @description REST API for Northwind orders
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetupRoutes(app, cfg, redis, authController, empController, categoryController, productController, orderController)

	log.Fatal(app.Listen(":" + cfg.Port))
}
