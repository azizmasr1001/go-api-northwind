package repositories

import (
	"github.com/azizmasr1001/go-api-northwind/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetAll() ([]models.Employee, error)
	GetPaginated(page int, limit int) ([]models.Employee, int64, error)
	GetByID(id int) (*models.Employee, error)
	Create(employee *models.Employee) (models.Employee, error)
	Update(employee *models.Employee) (models.Employee, error)
	Delete(id int) error
}

type employeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepo{db}
}

func (r *employeeRepo) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Find(&employees).Error
	return employees, err
}

func (r *employeeRepo) GetByID(id int) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.First(&employee, id).Error
	return &employee, err
}

func (r *employeeRepo) GetPaginated(page int, limit int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	offset := (page - 1) * limit
	if err := r.db.Model(&models.Employee{}).Count(&total).Limit(limit).Offset(offset).Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *employeeRepo) Create(employee *models.Employee) (models.Employee, error) {
	if err := r.db.Create(employee).Error; err != nil {
		return models.Employee{}, err
	}
	return *employee, nil
}

func (r *employeeRepo) Update(employee *models.Employee) (models.Employee, error) {
	if err := r.db.Save(employee).Error; err != nil {
		return models.Employee{}, err
	}
	return *employee, nil
}

func (r *employeeRepo) Delete(id int) error {
	return r.db.Delete(&models.Employee{}, id).Error
}
