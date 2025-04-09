package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// MockOrderRepository is a mock implementation of the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByID(id string) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAll(filter models.OrderFilter, offset, limit int) ([]models.Order, int, error) {
	args := m.Called(filter, offset, limit)
	return args.Get(0).([]models.Order), args.Int(1), args.Error(2)
}

func (m *MockOrderRepository) Update(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockOrderRepository)
	
	// Create a sample order
	order := &models.Order{
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		OptionType:     models.OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
	}
	
	// Set up the mock repository expectations
	mockRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(order, nil)
	
	// Create the service with the mock repository
	service := NewOrderService(mockRepo)
	
	// Call the service method
	createdOrder, err := service.CreateOrder(order)
	
	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.Equal(t, models.OrderStatusPending, createdOrder.Status)
	assert.Equal(t, 0, createdOrder.FilledQuantity)
	
	// Verify that the mock repository was called
	mockRepo.AssertExpectations(t)
}

func TestGetOrderByID(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockOrderRepository)
	
	// Create a sample order
	order := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
	}
	
	// Set up the mock repository expectations
	mockRepo.On("GetByID", "order123").Return(order, nil)
	mockRepo.On("GetByID", "nonexistent").Return(nil, errors.New("order not found"))
	
	// Create the service with the mock repository
	service := NewOrderService(mockRepo)
	
	// Test successful retrieval
	retrievedOrder, err := service.GetOrderByID("order123")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	assert.Equal(t, order.ID, retrievedOrder.ID)
	
	// Test error case
	retrievedOrder, err = service.GetOrderByID("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, retrievedOrder)
	
	// Test empty ID
	retrievedOrder, err = service.GetOrderByID("")
	assert.Error(t, err)
	assert.Nil(t, retrievedOrder)
	
	// Verify that the mock repository was called
	mockRepo.AssertExpectations(t)
}

func TestGetOrders(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockOrderRepository)
	
	// Create sample orders
	orders := []models.Order{
		{
			ID:             "order123",
			UserID:         "user123",
			Symbol:         "NIFTY",
			Exchange:       "NSE",
			OrderType:      models.OrderTypeLimit,
			Direction:      models.OrderDirectionBuy,
			Quantity:       10,
			Price:          500.50,
			Status:         models.OrderStatusPending,
			ProductType:    models.ProductTypeMIS,
			InstrumentType: models.InstrumentTypeOption,
		},
		{
			ID:             "order456",
			UserID:         "user123",
			Symbol:         "BANKNIFTY",
			Exchange:       "NSE",
			OrderType:      models.OrderTypeMarket,
			Direction:      models.OrderDirectionSell,
			Quantity:       5,
			Price:          0,
			Status:         models.OrderStatusExecuted,
			ProductType:    models.ProductTypeMIS,
			InstrumentType: models.InstrumentTypeFuture,
		},
	}
	
	// Set up the mock repository expectations
	mockRepo.On("GetAll", mock.AnythingOfType("models.OrderFilter"), 0, 50).Return(orders, 2, nil)
	mockRepo.On("GetAll", mock.AnythingOfType("models.OrderFilter"), 50, 50).Return([]models.Order{}, 2, nil)
	
	// Create the service with the mock repository
	service := NewOrderService(mockRepo)
	
	// Test successful retrieval with default pagination
	filter := models.OrderFilter{UserID: "user123"}
	retrievedOrders, total, err := service.GetOrders(filter, 1, 50)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedOrders))
	assert.Equal(t, 2, total)
	
	// Test pagination (page 2, no results)
	retrievedOrders, total, err = service.GetOrders(filter, 2, 50)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(retrievedOrders))
	assert.Equal(t, 2, total)
	
	// Test invalid pagination parameters (should use defaults)
	mockRepo.On("GetAll", mock.AnythingOfType("models.OrderFilter"), 0, 50).Return(orders, 2, nil)
	retrievedOrders, total, err = service.GetOrders(filter, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedOrders))
	assert.Equal(t, 2, total)
	
	// Test limit capping
	mockRepo.On("GetAll", mock.AnythingOfType("models.OrderFilter"), 0, 100).Return(orders, 2, nil)
	retrievedOrders, total, err = service.GetOrders(filter, 1, 200)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedOrders))
	assert.Equal(t, 2, total)
	
	// Verify that the mock repository was called
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrder(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockOrderRepository)
	
	// Create sample orders
	existingOrder := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
		CreatedAt:      time.Now().Add(-time.Hour),
	}
	
	updatedOrder := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       15, // Updated quantity
		Price:          550.75, // Updated price
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
	}
	
	cancelledOrder := &models.Order{
		ID:             "order456",
		UserID:         "user123",
		Symbol:         "BANKNIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeMarket,
		Direction:      models.OrderDirectionSell,
		Quantity:       5,
		Price:          0,
		Status:         models.OrderStatusCancelled,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeFuture,
	}
	
	// Set up the mock repository expectations
	mockRepo.On("GetByID", "order123").Return(existingOrder, nil)
	mockRepo.On("GetByID", "order456").Return(cancelledOrder, nil)
	mockRepo.On("GetByID", "nonexistent").Return(nil, errors.New("order not found"))
	mockRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(updatedOrder, nil)
	
	// Create the service with the mock repository
	service := NewOrderService(mockRepo)
	
	// Test successful update
	result, err := service.UpdateOrder(updatedOrder)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedOrder.Quantity, result.Quantity)
	assert.Equal(t, updatedOrder.Price, result.Price)
	
	// Test update of cancelled order (should fail)
	_, err = service.UpdateOrder(cancelledOrder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cancelled or rejected orders cannot be updated")
	
	// Test update of non-existent order
	nonexistentOrder := &models.Order{ID: "nonexistent"}
	_, err = service.UpdateOrder(nonexistentOrder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
	
	// Verify that the mock repository was called
	mockRepo.AssertExpectations(t)
}

func TestCancelOrder(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockOrderRepository)
	
	// Create sample orders
	pendingOrder := &models.Order{
		ID:             "order123",
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeLimit,
		Direction:      models.OrderDirectionBuy,
		Quantity:       10,
		Price:          500.50,
		Status:         models.OrderStatusPending,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeOption,
	}
	
	executedOrder := &models.Order{
		ID:             "order456",
		UserID:         "user123",
		Symbol:         "BANKNIFTY",
		Exchange:       "NSE",
		OrderType:      models.OrderTypeMarket,
		Direction:      models.OrderDirectionSell,
		Quantity:       5,
		Price:          0,
		Status:         models.OrderStatusExecuted,
		ProductType:    models.ProductTypeMIS,
		InstrumentType: models.InstrumentTypeFuture,
	}
	
	// Set up the mock repository expectations
	mockRepo.On("GetByID", "order123").Return(pendingOrder, nil)
	mockRepo.On("GetByID", "order456").Return(executedOrder, nil)
	mockRepo.On("GetByID", "nonexistent").Return(nil, errors.New("order not found"))
	mockRepo.On("Update", mock.MatchedBy(func(order *models.Order) bool {
		return order.ID == "order123" && order.Status == models.OrderStatusCancelled
	})).Return(pendingOrder, nil)
	
	// Create the service with the mock repository
	service := NewOrderService(mockRepo)
	
	// Test successful cancellation
	err := service.CancelOrder("order123")
	assert.NoError(t, err)
	
	// Test cancellation of executed order (should fail)
	err = service.CancelOrder("order456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only pending or partially filled orders can be cancelled")
	
	// Test cancellation of non-existent order
	err = service.CancelOrder("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
	
	// Test cancellation with empty ID
	err = service.CancelOrder("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order ID is required")
	
	// Verify that the mock repository was called
	mockRepo.AssertExpectations(t)
}
