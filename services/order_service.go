package services

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/repositories"
)

type OrderService interface {
	Create(order *models.Order, details []models.OrderDetail) (*models.Order, error)
	GetAll(page, limit int) ([]models.Order, int64, error)
	GetByID(id int) (*models.Order, []models.OrderDetail, error)
	Update(id int, order *models.Order, details []models.OrderDetail) (*models.Order, error)
	Delete(id int) error
}

type orderService struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{repo}
}

func (s *orderService) Create(order *models.Order, details []models.OrderDetail) (*models.Order, error) {
	err := s.repo.CreateOrderWithDetails(order, details)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetAll(page, limit int) ([]models.Order, int64, error) {
	return s.repo.GetAll(page, limit)
}

func (s *orderService) GetByID(id int) (*models.Order, []models.OrderDetail, error) {
	return s.repo.GetByID(id)
}

func (s *orderService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *orderService) Update(id int, order *models.Order, details []models.OrderDetail) (*models.Order, error) {
	err := s.repo.Update(id, order, details)
	if err != nil {
		return nil, err
	}
	order.OrderID = id
	return order, nil
}
