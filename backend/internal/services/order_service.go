package services

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// OrderService defines the interface for order-related operations
type OrderService interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	GetOrders(filter models.OrderFilter, page, limit int) ([]models.Order, int, error)
	UpdateOrder(order *models.Order) (*models.Order, error)
	CancelOrder(id string) error
}

// OrderServiceImpl implements the OrderService interface
type OrderServiceImpl struct {
	orderRepo repositories.OrderRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo repositories.OrderRepository) OrderService {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
	}
}

// CreateOrder creates a new order
func (s *OrderServiceImpl) CreateOrder(order *models.Order) (*models.Order, error) {
	// Validate the order
	if err := order.Validate(); err != nil {
		return nil, err
	}

	// Set initial values
	order.Status = models.OrderStatusPending
	order.FilledQuantity = 0
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Create the order
	createdOrder, err := s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderServiceImpl) GetOrderByID(id string) (*models.Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrders retrieves orders with filtering and pagination
func (s *OrderServiceImpl) GetOrders(filter models.OrderFilter, page, limit int) ([]models.Order, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100 // Maximum limit to prevent excessive queries
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get orders with pagination
	orders, total, err := s.orderRepo.GetAll(filter, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateOrder updates an existing order
func (s *OrderServiceImpl) UpdateOrder(order *models.Order) (*models.Order, error) {
	// Validate the order
	if err := order.Validate(); err != nil {
		return nil, err
	}

	// Check if order exists
	existingOrder, err := s.orderRepo.GetByID(order.ID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if order can be updated
	if existingOrder.Status == models.OrderStatusCancelled || 
	   existingOrder.Status == models.OrderStatusRejected {
		return nil, errors.New("cancelled or rejected orders cannot be updated")
	}

	// Preserve certain fields from the existing order
	order.CreatedAt = existingOrder.CreatedAt
	order.UpdatedAt = time.Now()

	// Update the order
	updatedOrder, err := s.orderRepo.Update(order)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

// CancelOrder cancels an existing order
func (s *OrderServiceImpl) CancelOrder(id string) error {
	if id == "" {
		return errors.New("order ID is required")
	}

	// Check if order exists
	existingOrder, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	// Check if order can be cancelled
	if existingOrder.Status != models.OrderStatusPending && 
	   existingOrder.Status != models.OrderStatusPartial {
		return errors.New("only pending or partially filled orders can be cancelled")
	}

	// Update order status
	existingOrder.Status = models.OrderStatusCancelled
	existingOrder.UpdatedAt = time.Now()

	// Save the updated order
	_, err = s.orderRepo.Update(existingOrder)
	if err != nil {
		return err
	}

	return nil
}
