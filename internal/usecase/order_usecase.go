package usecase

import (
	"delivery-state-manager/internal/models"
	"delivery-state-manager/pkg/errs"
)

// OrderRepository defines the interface for order operations
type OrderRepository interface {
	CreateOrder(order *models.Order)
	GetOrder(id string) (*models.Order, error)
	GetAllOrders() []*models.Order
	UpdateOrderStatus(id string, status models.OrderStatus) error
}

// OrderUseCase handles order-related use cases
type OrderUseCase struct {
	repo OrderRepository
}

// NewOrderUseCase creates a new OrderUseCase instance
func NewOrderUseCase(repo OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		repo: repo,
	}
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(order *models.Order) error {
	// Validate required fields
	if order.ID == "" || order.Customer == "" {
		return errs.ErrMissingRequiredField
	}

	uc.repo.CreateOrder(order)
	return nil
}

// GetOrder retrieves an order by ID
func (uc *OrderUseCase) GetOrder(id string) (*models.Order, error) {
	return uc.repo.GetOrder(id)
}

// GetAllOrders returns all orders
func (uc *OrderUseCase) GetAllOrders() []*models.Order {
	return uc.repo.GetAllOrders()
}

// UpdateOrderStatus updates the status of an order
func (uc *OrderUseCase) UpdateOrderStatus(id string, status models.OrderStatus) error {
	return uc.repo.UpdateOrderStatus(id, status)
}
