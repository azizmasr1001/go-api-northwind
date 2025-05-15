package routes

import (
	"github.com/aronipurwanto/go-api-northwind/config"
	"github.com/aronipurwanto/go-api-northwind/controllers"
	"github.com/aronipurwanto/go-api-northwind/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg config.Config,
	authCtrl *controllers.AuthController,
	empCtrl *controllers.EmployeeController,
	catCtrl *controllers.CategoryController,
	prodCtrl *controllers.ProductController,
	orderCtrl *controllers.OrderController) {

	api := app.Group("/api")
	api.Post("/login", authCtrl.Login)

	employee := api.Group("/employees", middlewares.Protected(cfg.JWTSecret))
	employee.Get("/", empCtrl.GetAll)
	employee.Get("/:id", empCtrl.GetByID)
	employee.Post("/", empCtrl.Create)
	employee.Put("/:id", empCtrl.Update)
	employee.Delete("/:id", empCtrl.Delete)

	category := api.Group("/categories", middlewares.Protected(cfg.JWTSecret))
	category.Get("/", catCtrl.GetAll)
	category.Post("/", catCtrl.Create)
	category.Get("/:id", catCtrl.GetByID)
	category.Put("/:id", catCtrl.Update)
	category.Delete("/:id", catCtrl.Delete)

	product := api.Group("/products", middlewares.Protected(cfg.JWTSecret))
	//product.Get("/", prodCtrl.GetAll)
	product.Get("/", middlewares.ValidateQueryPagination(1, 10), prodCtrl.GetAll)
	product.Post("/", prodCtrl.Create)
	product.Get("/:id", middlewares.ValidateIDParam("id"), prodCtrl.GetByID)
	product.Put("/:id", middlewares.ValidateIDParam("id"), prodCtrl.Update)
	product.Delete("/:id", middlewares.ValidateIDParam("id"), prodCtrl.Delete)

	order := api.Group("/orders", middlewares.Protected(cfg.JWTSecret))
	order.Get("/", middlewares.ValidateQueryPagination(1, 10), orderCtrl.GetAll)
	order.Get("/:id", middlewares.ValidateIDParam("id"), orderCtrl.GetByID)
	order.Post("/", orderCtrl.Create)
	order.Put("/:id", middlewares.ValidateIDParam("id"), orderCtrl.Update)
	order.Delete("/:id", middlewares.ValidateIDParam("id"), orderCtrl.Delete)
	order.Post("/", orderCtrl.Create)
}
