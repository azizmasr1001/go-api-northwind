package repositories

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetPaginated(page, limit int) ([]models.Product, int64, error)
	SearchByName(name string, page, limit int) ([]models.Product, int64, error)
	GetByID(id int) (*models.Product, error)
	Create(prod *models.Product) (*models.Product, error)
	Update(prod *models.Product) (*models.Product, error)
	Delete(id int) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db}
}

func (r *productRepo) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepo) GetPaginated(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Product{}).Count(&total).Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *productRepo) SearchByName(name string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("ProductName LIKE ?", "%"+name+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepo) GetByID(id int) (*models.Product, error) {
	var prod models.Product
	err := r.db.First(&prod, id).Error
	return &prod, err
}

func (r *productRepo) Create(prod *models.Product) (*models.Product, error) {

	if err := r.db.Create(prod).Error; err != nil {
		return nil, err
	}
	return prod, nil
}

func (r *productRepo) Update(prod *models.Product) (*models.Product, error) {
	if err := r.db.Save(prod).Error; err != nil {
		return nil, err
	}
	return prod, nil
}

func (r *productRepo) Delete(id int) error {
	return r.db.Delete(&models.Product{}, id).Error
}
