package services

import (
	"github.com/azizmasr1001/go-api-northwind/models"
	"github.com/azizmasr1001/go-api-northwind/repositories"
)

type ProductService interface {
	GetAll() ([]models.Product, error)
	GetAllPaginated(page, limit int) ([]models.Product, int64, error)
	SearchByName(name string, page, limit int) ([]models.Product, int64, error)
	GetByID(id int) (*models.Product, error)
	Create(prod *models.Product) (*models.Product, error)
	Update(prod *models.Product) (*models.Product, error)
	Delete(id int) error
}

type productServiceImpl struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productServiceImpl{repo}
}

func (s *productServiceImpl) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *productServiceImpl) GetAllPaginated(page, limit int) ([]models.Product, int64, error) {
	return s.repo.GetPaginated(page, limit)
}

func (s *productServiceImpl) SearchByName(name string, page, limit int) ([]models.Product, int64, error) {
	return s.repo.SearchByName(name, page, limit)
}

func (s *productServiceImpl) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productServiceImpl) Create(prod *models.Product) (*models.Product, error) {
	return s.repo.Create(prod)
}

func (s *productServiceImpl) Update(prod *models.Product) (*models.Product, error) {
	return s.repo.Update(prod)
}

func (s *productServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}
