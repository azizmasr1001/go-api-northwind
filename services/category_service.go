package services

import (
	"github.com/azizmasr1001/go-api-northwind/models"
	"github.com/azizmasr1001/go-api-northwind/repositories"
)

type CategoryService interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(cat *models.Category) (models.Category, error)
	Update(cat *models.Category) (models.Category, error)
	Delete(id int) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Create(cat *models.Category) (models.Category, error) {
	return s.repo.Create(cat)
}

func (s *categoryService) Update(cat *models.Category) (models.Category, error) {
	return s.repo.Update(cat)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
