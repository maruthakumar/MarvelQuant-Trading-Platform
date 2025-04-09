package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/api"
	"github.com/trading-platform/backend/internal/core"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"
)

// MockOrderService is a mock implementation of the order service
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) PlaceOrder(ctx context.Context, request core.OrderRequest) (*core.OrderResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*core.OrderResponse), args.Error(1)
}

func (m *MockOrderService) ModifyOrder(ctx context.Context, orderID string, request core.OrderRequest) (*core.OrderResponse, error) {
	args := m.Called(ctx, orderID, request)
	return args.Get(0).(*core.OrderResponse), args.Error(1)
}

func (m *MockOrderService) CancelOrder(ctx context.Context, orderID string, brokerName string) (*core.OrderResponse, error) {
	args := m.Called(ctx, orderID, brokerName)
	return args.Get(0).(*core.OrderResponse), args.Error(1)
}

func (m *MockOrderService) GetOrder(ctx context.Context, orderID string, brokerName string) (*core.Order, error) {
	args := m.Called(ctx, orderID, brokerName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Order), args.Error(1)
}

func (m *MockOrderService) GetOrders(ctx context.Context, brokerName string) ([]core.Order, error) {
	args := m.Called(ctx, brokerName)
	return args.Get(0).([]core.Order), args.Error(1)
}

func TestAPIRoutes(t *testing.T) {
	// Create mock services
	mockOrderService := new(MockOrderService)
	
	// Create API handlers with mock services
	handlers := api.NewHandlers(mockOrderService, nil, nil)
	
	// Set up Gin for testing
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Register API routes
	api.RegisterRoutes(router, handlers)
	
	// Test placing an order
	t.Run("PlaceOrder", func(t *testing.T) {
		// Create a test order request
		orderRequest := core.OrderRequest{
			UserID:          "user123",
			BrokerName:      "xts",
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       "MARKET",
			TransactionType: "BUY",
			ProductType:     "NRML",
			Quantity:        1,
		}
		
		// Convert the request to JSON
		requestBody, err := json.Marshal(orderRequest)
		require.NoError(t, err)
		
		// Set up the mock to return a successful response
		expectedResponse := &core.OrderResponse{
			Success: true,
			OrderID: "order123",
		}
		mockOrderService.On("PlaceOrder", mock.Anything, mock.Anything).Return(expectedResponse, nil)
		
		// Create a test request
		req := httptest.NewRequest("POST", "/api/orders", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		
		// Create a test response recorder
		w := httptest.NewRecorder()
		
		// Serve the request
		router.ServeHTTP(w, req)
		
		// Check that the response status is OK
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Parse the response body
		var response core.OrderResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Check that the response is as expected
		assert.Equal(t, expectedResponse.Success, response.Success)
		assert.Equal(t, expectedResponse.OrderID, response.OrderID)
		
		// Verify that the mock method was called
		mockOrderService.AssertExpectations(t)
	})
	
	// Test getting an order
	t.Run("GetOrder", func(t *testing.T) {
		// Set up the mock to return an order
		expectedOrder := &core.Order{
			ID:              "order123",
			UserID:          "user123",
			BrokerName:      "xts",
			Symbol:          "NIFTY",
			Exchange:        "NSE",
			OrderType:       "MARKET",
			TransactionType: "BUY",
			ProductType:     "NRML",
			Quantity:        1,
			Status:          "COMPLETED",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockOrderService.On("GetOrder", mock.Anything, "order123", "xts").Return(expectedOrder, nil)
		
		// Create a test request
		req := httptest.NewRequest("GET", "/api/orders/order123?broker=xts", nil)
		
		// Create a test response recorder
		w := httptest.NewRecorder()
		
		// Serve the request
		router.ServeHTTP(w, req)
		
		// Check that the response status is OK
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Parse the response body
		var order core.Order
		err := json.Unmarshal(w.Body.Bytes(), &order)
		require.NoError(t, err)
		
		// Check that the order is as expected
		assert.Equal(t, expectedOrder.ID, order.ID)
		assert.Equal(t, expectedOrder.Symbol, order.Symbol)
		assert.Equal(t, expectedOrder.OrderType, order.OrderType)
		assert.Equal(t, expectedOrder.Status, order.Status)
		
		// Verify that the mock method was called
		mockOrderService.AssertExpectations(t)
	})
	
	// Test getting all orders
	t.Run("GetOrders", func(t *testing.T) {
		// Set up the mock to return orders
		expectedOrders := []core.Order{
			{
				ID:              "order123",
				UserID:          "user123",
				BrokerName:      "xts",
				Symbol:          "NIFTY",
				Exchange:        "NSE",
				OrderType:       "MARKET",
				TransactionType: "BUY",
				ProductType:     "NRML",
				Quantity:        1,
				Status:          "COMPLETED",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              "order456",
				UserID:          "user123",
				BrokerName:      "xts",
				Symbol:          "BANKNIFTY",
				Exchange:        "NSE",
				OrderType:       "LIMIT",
				TransactionType: "SELL",
				ProductType:     "NRML",
				Quantity:        2,
				Price:           40000.0,
				Status:          "OPEN",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}
		mockOrderService.On("GetOrders", mock.Anything, "xts").Return(expectedOrders, nil)
		
		// Create a test request
		req := httptest.NewRequest("GET", "/api/orders?broker=xts", nil)
		
		// Create a test response recorder
		w := httptest.NewRecorder()
		
		// Serve the request
		router.ServeHTTP(w, req)
		
		// Check that the response status is OK
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Parse the response body
		var orders []core.Order
		err := json.Unmarshal(w.Body.Bytes(), &orders)
		require.NoError(t, err)
		
		// Check that the orders are as expected
		assert.Equal(t, len(expectedOrders), len(orders))
		assert.Equal(t, expectedOrders[0].ID, orders[0].ID)
		assert.Equal(t, expectedOrders[1].ID, orders[1].ID)
		
		// Verify that the mock method was called
		mockOrderService.AssertExpectations(t)
	})
	
	// Test cancelling an order
	t.Run("CancelOrder", func(t *testing.T) {
		// Set up the mock to return a successful response
		expectedResponse := &core.OrderResponse{
			Success: true,
			OrderID: "order123",
		}
		mockOrderService.On("CancelOrder", mock.Anything, "order123", "xts").Return(expectedResponse, nil)
		
		// Create a test request
		req := httptest.NewRequest("DELETE", "/api/orders/order123?broker=xts", nil)
		
		// Create a test response recorder
		w := httptest.NewRecorder()
		
		// Serve the request
		router.ServeHTTP(w, req)
		
		// Check that the response status is OK
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Parse the response body
		var response core.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Check that the response is as expected
		assert.Equal(t, expectedResponse.Success, response.Success)
		assert.Equal(t, expectedResponse.OrderID, response.OrderID)
		
		// Verify that the mock method was called
		mockOrderService.AssertExpectations(t)
	})
}
