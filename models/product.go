package models

type Product struct {
	ProductID       int     `gorm:"primaryKey;column:ProductID" json:"product_id"`
	ProductName     string  `gorm:"column:ProductName" json:"product_name" validate:"required"`
	SupplierID      *int    `gorm:"column:SupplierID" json:"supplier_id"`
	CategoryID      *int    `gorm:"column:CategoryID" json:"category_id"`
	QuantityPerUnit string  `gorm:"column:QuantityPerUnit" json:"quantity_per_unit"`
	UnitPrice       float64 `gorm:"column:UnitPrice" json:"unit_price"`
	UnitsInStock    int     `gorm:"column:UnitsInStock" json:"units_in_stock"`
	UnitsOnOrder    int     `gorm:"column:UnitsOnOrder" json:"units_on_order"`
	ReorderLevel    int     `gorm:"column:ReorderLevel" json:"reorder_level"`
	Discontinued    bool    `gorm:"column:Discontinued" json:"discontinued"`
}
