package models

type Category struct {
	CategoryID   int    `gorm:"primaryKey;column:CategoryID" json:"category_id"`
	CategoryName string `gorm:"column:CategoryName" json:"category_name" validate:"required"`
	Description  string `gorm:"column:Description" json:"description"`
}
