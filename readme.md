# 📘 Northwind REST API - Golang (Fiber, GORM, Viper, SQL Server)

A lightweight REST API for the Northwind sample database using:

- Golang
- Fiber - Web framework
- GORM - ORM
- Viper - Config loader
- Validator - Input validation
- SQL Server as database
- Clean architecture using MVC + Repository Pattern + Dependency Injection

## Dependency
- go get github.com/gofiber/fiber/v2 
- go get gorm.io/gorm 
- go get gorm.io/driver/sqlserver 
- go get github.com/spf13/viper 
- go get github.com/go-playground/validator/v10
- go get github.com/swaggo/fiber-swagger
- go get github.com/swaggo/swag/cmd/swag
- go get github.com/golang-jwt/jwt/v5
- go get github.com/gofiber/jwt/v3
---

## 📁 Project Structure
````
northwind-api/
├── config/                # Load .env config
├── controllers/           # HTTP handlers
├── models/                # GORM model structs
├── repositories/          # Data access layer
├── services/              # Business logic
├── routes/                # Route registration
├── main.go                # App entry point
├── .env                   # App configuration
└── go.mod
````
---

## ✅ Features

- CRUD for `Employees`
- Clean separation of concerns with MVC
- SQL Server support with GORM
- Input validation with go-playground/validator

---

## 🛠️ Requirements

- Go 1.18+
- SQL Server (Local or Remote)
- Northwind database installed

---