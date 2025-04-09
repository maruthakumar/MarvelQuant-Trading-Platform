package tests

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
	"go.mongodb.org/mongo-driver/mongo"

	"trading_platform/backend/internal/api"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/models"
)

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) (string, error) {
	args := m.Called(order)
	return args.String(0), args.Error(1)
}

func (m *MockOrderRepository) GetByID(id string) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrderRepository) Find(filter models.OrderFilter, page, limit int) ([]models.Order, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]models.Order), args.Int(1), args.Error(2)
}

// TestCreateOrder tests the create order endpoint
func TestCreateOrder(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Create test order
	order := models.Order{
		Symbol:      "AAPL",
		OrderType:   models.OrderTypeLimit,
		Side:        models.SideBuy,
		Quantity:    100,
		Price:       150.50,
		TimeInForce: models.TimeInForceGTC,
		Status:      models.OrderStatusNew,
	}

	// Set up expectations
	mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return("order123", nil)

	// Create request with authenticated context
	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.CreateOrder).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse response
	var response models.Order
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify order data
	assert.Equal(t, "order123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, models.OrderTypeLimit, response.OrderType)
	assert.Equal(t, models.SideBuy, response.Side)
	assert.Equal(t, float64(100), response.Quantity)
	assert.Equal(t, 150.50, response.Price)
	assert.Equal(t, models.TimeInForceGTC, response.TimeInForce)
	assert.Equal(t, models.OrderStatusNew, response.Status)

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOrder tests the get order endpoint
func TestGetOrder(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Create test order
	order := &models.Order{
		ID:          "order123",
		UserID:      "user123",
		Symbol:      "AAPL",
		OrderType:   models.OrderTypeLimit,
		Side:        models.SideBuy,
		Quantity:    100,
		Price:       150.50,
		TimeInForce: models.TimeInForceGTC,
		Status:      models.OrderStatusFilled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockOrderRepo.On("GetByID", "order123").Return(order, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/orders/order123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Order
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify order data
	assert.Equal(t, "order123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, models.OrderTypeLimit, response.OrderType)
	assert.Equal(t, models.SideBuy, response.Side)
	assert.Equal(t, float64(100), response.Quantity)
	assert.Equal(t, 150.50, response.Price)
	assert.Equal(t, models.TimeInForceGTC, response.TimeInForce)
	assert.Equal(t, models.OrderStatusFilled, response.Status)

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOrderUnauthorized tests the get order endpoint with unauthorized access
func TestGetOrderUnauthorized(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Create test order
	order := &models.Order{
		ID:          "order123",
		UserID:      "user456", // Different user ID
		Symbol:      "AAPL",
		OrderType:   models.OrderTypeLimit,
		Side:        models.SideBuy,
		Quantity:    100,
		Price:       150.50,
		TimeInForce: models.TimeInForceGTC,
		Status:      models.OrderStatusFilled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockOrderRepo.On("GetByID", "order123").Return(order, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/orders/order123", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response - should be forbidden
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}

// TestUpdateOrder tests the update order endpoint
func TestUpdateOrder(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Create existing order
	existingOrder := &models.Order{
		ID:          "order123",
		UserID:      "user123",
		Symbol:      "AAPL",
		OrderType:   models.OrderTypeLimit,
		Side:        models.SideBuy,
		Quantity:    100,
		Price:       150.50,
		TimeInForce: models.TimeInForceGTC,
		Status:      models.OrderStatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create order update
	orderUpdate := models.Order{
		Quantity: 200,
		Price:    155.75,
	}

	// Set up expectations
	mockOrderRepo.On("GetByID", "order123").Return(existingOrder, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(nil)

	// Create request with authenticated context
	body, _ := json.Marshal(orderUpdate)
	req, _ := http.NewRequest("PUT", "/orders/order123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.UpdateOrder).Methods("PUT")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Order
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify updated order data
	assert.Equal(t, "order123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, models.OrderTypeLimit, response.OrderType)
	assert.Equal(t, models.SideBuy, response.Side)
	assert.Equal(t, float64(200), response.Quantity) // Updated
	assert.Equal(t, 155.75, response.Price)          // Updated
	assert.Equal(t, models.TimeInForceGTC, response.TimeInForce)
	assert.Equal(t, models.OrderStatusNew, response.Status)

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}

// TestCancelOrder tests the cancel order endpoint
func TestCancelOrder(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Create existing order
	existingOrder := &models.Order{
		ID:          "order123",
		UserID:      "user123",
		Symbol:      "AAPL",
		OrderType:   models.OrderTypeLimit,
		Side:        models.SideBuy,
		Quantity:    100,
		Price:       150.50,
		TimeInForce: models.TimeInForceGTC,
		Status:      models.OrderStatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set up expectations
	mockOrderRepo.On("GetByID", "order123").Return(existingOrder, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("POST", "/orders/order123/cancel", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}/cancel", handler.CancelOrder).Methods("POST")
	router.ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.Order
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify order status is cancelled
	assert.Equal(t, "order123", response.ID)
	assert.Equal(t, "user123", response.UserID)
	assert.Equal(t, models.OrderStatusCancelled, response.Status)

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}

// TestGetOrders tests the get orders endpoint with filtering and pagination
func TestGetOrders(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Create test orders
	orders := []models.Order{
		{
			ID:          "order1",
			UserID:      "user123",
			Symbol:      "AAPL",
			OrderType:   models.OrderTypeLimit,
			Side:        models.SideBuy,
			Quantity:    100,
			Price:       150.50,
			TimeInForce: models.TimeInForceGTC,
			Status:      models.OrderStatusFilled,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "order2",
			UserID:      "user123",
			Symbol:      "MSFT",
			OrderType:   models.OrderTypeMarket,
			Side:        models.SideSell,
			Quantity:    50,
			TimeInForce: models.TimeInForceIOC,
			Status:      models.OrderStatusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Set up expectations
	mockOrderRepo.On("Find", mock.AnythingOfType("models.OrderFilter"), 1, 20).Return(orders, 2, nil)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/orders?symbol=AAPL&status=filled", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(handler.GetOrders).ServeHTTP(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify pagination data
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(20), response["limit"])
	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["totalPages"])

	// Verify orders data
	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}

// TestOrderNotFound tests the get order endpoint with non-existent order
func TestOrderNotFound(t *testing.T) {
	// Create mock repository
	mockOrderRepo := new(MockOrderRepository)

	// Create handler
	handler := api.NewOrderHandler(mockOrderRepo)

	// Set up expectations
	mockOrderRepo.On("GetByID", "nonexistent").Return(nil, mongo.ErrNoDocuments)

	// Create request with authenticated context
	req, _ := http.NewRequest("GET", "/orders/nonexistent", nil)
	ctx := auth.SetUserIDInContext(req.Context(), "user123")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// Set up router with URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check response - should be not found
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify expectations
	mockOrderRepo.AssertExpectations(t)
}
