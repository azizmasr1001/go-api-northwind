package models

import "time"

type Order struct {
	OrderID        int        `gorm:"primaryKey;column:OrderID" json:"order_id"`
	CustomerID     *string    `gorm:"column:CustomerID" json:"customer_id"`
	EmployeeID     *int       `gorm:"column:EmployeeID" json:"employee_id"`
	OrderDate      *time.Time `gorm:"column:OrderDate" json:"order_date"`
	RequiredDate   *time.Time `gorm:"column:RequiredDate" json:"required_date"`
	ShippedDate    *time.Time `gorm:"column:ShippedDate" json:"shipped_date"`
	ShipVia        *int       `gorm:"column:ShipVia" json:"ship_via"`
	Freight        float64    `gorm:"column:Freight" json:"freight"`
	ShipName       string     `gorm:"column:ShipName" json:"ship_name"`
	ShipAddress    string     `gorm:"column:ShipAddress" json:"ship_address"`
	ShipCity       string     `gorm:"column:ShipCity" json:"ship_city"`
	ShipRegion     string     `gorm:"column:ShipRegion" json:"ship_region"`
	ShipPostalCode string     `gorm:"column:ShipPostalCode" json:"ship_postal_code"`
	ShipCountry    string     `gorm:"column:ShipCountry" json:"ship_country"`
}

type OrderDetail struct {
	OrderID   int     `gorm:"column:OrderID;primaryKey" json:"order_id"`
	ProductID int     `gorm:"column:ProductID;primaryKey" json:"product_id"`
	UnitPrice float64 `gorm:"column:UnitPrice" json:"unit_price"`
	Quantity  int     `gorm:"column:Quantity" json:"quantity"`
	Discount  float64 `gorm:"column:Discount" json:"discount"`
}
