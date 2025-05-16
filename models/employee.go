package models

import "time"

type Employee struct {
	EmployeeID      int        `gorm:"column:EmployeeID;primaryKey" json:"employee_id"`
	LastName        string     `gorm:"column:LastName" json:"last_name" validate:"required"`
	FirstName       string     `gorm:"column:FirstName" json:"first_name" validate:"required"`
	Title           string     `gorm:"column:Title" json:"title"`
	TitleOfCourtesy string     `gorm:"column:TitleOfCourtesy" json:"title_of_courtesy"`
	BirthDate       *time.Time `gorm:"column:BirthDate" json:"birth_date"`
	HireDate        *time.Time `gorm:"column:HireDate" json:"hire_date"`
	Address         string     `gorm:"column:Address" json:"address"`
	City            string     `gorm:"column:City" json:"city"`
	Region          string     `gorm:"column:Region" json:"region"`
	PostalCode      string     `gorm:"column:PostalCode" json:"postal_code"`
	Country         string     `gorm:"column:Country" json:"country"`
	HomePhone       string     `gorm:"column:HomePhone" json:"home_phone"`
	Extension       string     `gorm:"column:Extension" json:"extension"`
	Notes           string     `gorm:"column:Notes" json:"notes"`
	ReportsTo       *int       `gorm:"column:ReportsTo" json:"reports_to"`
	PhotoPath       string     `gorm:"column:PhotoPath" json:"photo_path"`
}
