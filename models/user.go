package models

import "time"

type User struct {
	UserID       int       `gorm:"column:UserID;primaryKey" json:"user_id"`
	Username     string    `gorm:"column:Username" json:"username" validate:"required"`
	Email        string    `gorm:"column:Email" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"column:PasswordHash" json:"-"`
	Role         string    `gorm:"column:Role" json:"role"`
	IsActive     bool      `gorm:"column:IsActive" json:"is_active"`
	EmployeeID   *int      `gorm:"column:EmployeeID" json:"employee_id"`
	CreatedAt    time.Time `gorm:"column:CreatedAt" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:UpdatedAt" json:"updated_at"`
}

type RegisterRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	Role       string `json:"role"`
	EmployeeID *int   `json:"employee_id"` // opsional
}
