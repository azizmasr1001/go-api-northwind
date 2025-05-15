package services

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/repositories"
)

type EmployeeService interface {
	GetAll() ([]models.Employee, error)
	GetAllPaginated(page, limit int) ([]models.Employee, int64, error)
	GetByID(id int) (*models.Employee, error)
	Create(emp *models.Employee) (models.Employee, error)
	Update(emp *models.Employee) (models.Employee, error)
	Delete(id int) error
}

type employeeService struct {
	repo repositories.EmployeeRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{repo}
}

func (s *employeeService) GetAll() ([]models.Employee, error) {
	return s.repo.GetAll()
}

func (s *employeeService) GetAllPaginated(page, limit int) ([]models.Employee, int64, error) {
	return s.repo.GetPaginated(page, limit)
}

func (s *employeeService) GetByID(id int) (*models.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *employeeService) Create(emp *models.Employee) (models.Employee, error) {
	return s.repo.Create(emp)
}

func (s *employeeService) Update(emp *models.Employee) (models.Employee, error) {
	return s.repo.Update(emp)
}

func (s *employeeService) Delete(id int) error {
	return s.repo.Delete(id)
}
