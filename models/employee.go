package models

import "time"

type Employee struct {
	EmployeeID      int        `gorm:"primaryKey;column:EmployeeID" json:"employee_id"`
	LastName        string     `json:"last_name" validate:"required"`
	FirstName       string     `json:"first_name" validate:"required"`
	Title           string     `json:"title"`
	TitleOfCourtesy string     `json:"title_of_courtesy"`
	BirthDate       *time.Time `json:"birth_date"`
	HireDate        *time.Time `json:"hire_date"`
	Address         string     `json:"address"`
	City            string     `json:"city"`
	Region          string     `json:"region"`
	PostalCode      string     `json:"postal_code"`
	Country         string     `json:"country"`
	HomePhone       string     `json:"home_phone"`
	Extension       string     `json:"extension"`
	Notes           string     `json:"notes"`
	ReportsTo       *int       `json:"reports_to"`
	PhotoPath       string     `json:"photo_path"`
}
