package main

import (
	"github.com/aronipurwanto/go-api-northwind/config"
	"github.com/aronipurwanto/go-api-northwind/controllers"
	"github.com/aronipurwanto/go-api-northwind/repositories"
	"github.com/aronipurwanto/go-api-northwind/routes"
	"github.com/aronipurwanto/go-api-northwind/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"

	"github.com/gofiber/swagger"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(sqlserver.Open(cfg.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	empRepo := repositories.NewEmployeeRepository(db)
	empService := services.NewEmployeeService(empRepo)
	empController := controllers.NewEmployeeController(empService)

	app := fiber.New()
	authController := controllers.NewAuthController(db, cfg)

	// inisialisasi repo, service, controller
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo)
	orderController := controllers.NewOrderController(orderService)

	routes.SetupRoutes(app, cfg, authController, empController, categoryController, productController, orderController)

	// @title Northwind API
	// @version 1.0
	// @description REST API for Northwind orders
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":" + cfg.Port))
}
