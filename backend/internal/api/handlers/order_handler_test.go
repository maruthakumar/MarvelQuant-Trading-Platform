package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trading-platform/backend/internal/models"
)

// MockOrderService is a mock implementation of the OrderService interface
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrderByID(id string) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrders(filter models.OrderFilter, page, limit int) ([]models.Order, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]models.Order), args.Int(1), args.Error(2)
}

func (m *MockOrderService) UpdateOrder(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) CancelOrder(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	// Create a mock order service
	mockService := new(MockOrderService)
	
	// Create a sample order
	order := &models.Order{
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
		OptionType:     models.OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
	}
	
	// Set up the mock service expectations
	mockService.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(order, nil)
	
	// Create the handler with the mock service
	handler := NewOrderHandler(mockService)
	
	// Create a request body
	orderJSON, _ := json.Marshal(order)
	req, err := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(orderJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.CreateOrder(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
	
	// Parse the response
	var response models.Order
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, order.UserID, response.UserID)
	assert.Equal(t, order.Symbol, response.Symbol)
	assert.Equal(t, order.OrderType, response.OrderType)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetOrder(t *testing.T) {
	// Create a mock order service
	mockService := new(MockOrderService)
	
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
		OptionType:     models.OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
	}
	
	// Set up the mock service expectations
	mockService.On("GetOrderByID", "order123").Return(order, nil)
	
	// Create the handler with the mock service
	handler := NewOrderHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/orders/order123", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}", handler.GetOrder)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.Order
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, order.ID, response.ID)
	assert.Equal(t, order.UserID, response.UserID)
	assert.Equal(t, order.Symbol, response.Symbol)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetOrders(t *testing.T) {
	// Create a mock order service
	mockService := new(MockOrderService)
	
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
	
	// Set up the mock service expectations
	mockService.On("GetOrders", mock.AnythingOfType("models.OrderFilter"), 1, 50).Return(orders, 2, nil)
	
	// Create the handler with the mock service
	handler := NewOrderHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/orders?userId=user123", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetOrders(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(50), response["limit"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestUpdateOrder(t *testing.T) {
	// Create a mock order service
	mockService := new(MockOrderService)
	
	// Create a sample order
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
		CreatedAt:      time.Now(),
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
		CreatedAt:      existingOrder.CreatedAt,
	}
	
	// Set up the mock service expectations
	mockService.On("GetOrderByID", "order123").Return(existingOrder, nil)
	mockService.On("UpdateOrder", mock.AnythingOfType("*models.Order")).Return(updatedOrder, nil)
	
	// Create the handler with the mock service
	handler := NewOrderHandler(mockService)
	
	// Create a request body
	orderJSON, _ := json.Marshal(updatedOrder)
	req, err := http.NewRequest("PUT", "/api/orders/order123", bytes.NewBuffer(orderJSON))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}", handler.UpdateOrder)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response models.Order
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, updatedOrder.ID, response.ID)
	assert.Equal(t, updatedOrder.Quantity, response.Quantity)
	assert.Equal(t, updatedOrder.Price, response.Price)
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestCancelOrder(t *testing.T) {
	// Create a mock order service
	mockService := new(MockOrderService)
	
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
	
	// Set up the mock service expectations
	mockService.On("GetOrderByID", "order123").Return(order, nil)
	mockService.On("CancelOrder", "order123").Return(nil)
	
	// Create the handler with the mock service
	handler := NewOrderHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("POST", "/api/orders/order123/cancel", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/orders/{id}/cancel", handler.CancelOrder)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, "Order cancelled successfully", response["message"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}

func TestGetOrdersByUser(t *testing.T) {
	// Create a mock order service
	mockService := new(MockOrderService)
	
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
	
	// Set up the mock service expectations
	mockService.On("GetOrders", mock.MatchedBy(func(filter models.OrderFilter) bool {
		return filter.UserID == "user123"
	}), 1, 50).Return(orders, 2, nil)
	
	// Create the handler with the mock service
	handler := NewOrderHandler(mockService)
	
	// Create a request
	req, err := http.NewRequest("GET", "/api/users/user123/orders", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Set up the router to get the URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{userId}/orders", handler.GetOrdersByUser)
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	router.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check the response
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(50), response["limit"])
	
	// Verify that the mock service was called
	mockService.AssertExpectations(t)
}
