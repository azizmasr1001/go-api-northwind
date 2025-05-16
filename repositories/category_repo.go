package repositories

import (
	"github.com/azizmasr1001/go-api-northwind/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(cat *models.Category) (models.Category, error)
	Update(cat *models.Category) (models.Category, error)
	Delete(id int) error
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db}
}

func (r *categoryRepo) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepo) GetByID(id int) (*models.Category, error) {
	var cat models.Category
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *categoryRepo) Create(cat *models.Category) (models.Category, error) {
	if err := r.db.Create(cat).Error; err != nil {
		return models.Category{}, err
	}
	return *cat, nil
}

func (r *categoryRepo) Update(cat *models.Category) (models.Category, error) {
	if err := r.db.Save(cat).Error; err != nil {
		return models.Category{}, err
	}
	return *cat, nil
}

func (r *categoryRepo) Delete(id int) error {
	return r.db.Delete(&models.Category{}, id).Error
}
