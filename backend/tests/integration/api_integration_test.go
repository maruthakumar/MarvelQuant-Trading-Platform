package integration

import (
	"context"
	"testing"
	"time"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/gin-gonic/gin"
	"github.com/trading-platform/backend/internal/api"
	"github.com/trading-platform/backend/internal/auth"
	"github.com/trading-platform/backend/internal/middleware"
	"github.com/trading-platform/backend/internal/core"
)

func TestAPIIntegration(t *testing.T) {
	// Set up the test environment
	dbService, mqService, authService, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create mock services
	mockOrderService := &MockOrderService{}
	mockPortfolioService := &MockPortfolioService{}
	mockStrategyService := &MockStrategyService{}

	// Create API handlers
	handlers := api.NewHandlers(mockOrderService, mockPortfolioService, mockStrategyService)

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Set up Gin for testing
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register API routes with authentication
	api.RegisterAuthRoutes(router, authService)
	
	// Create authenticated group
	authGroup := router.Group("/api")
	authGroup.Use(authMiddleware.RequireAuth())
	
	// Register protected routes
	api.RegisterProtectedRoutes(authGroup, handlers)

	// Test the complete API flow with authentication
	t.Run("AuthenticatedAPIFlow", func(t *testing.T) {
		// 1. Register a user
		registerRequest := auth.RegisterRequest{
			Username: "apitest",
			Email:    "api@test.com",
			Password: "password123",
		}

		registerBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var registerResponse auth.AuthResponse
		err = json.Unmarshal(w.Body.Bytes(), &registerResponse)
		require.NoError(t, err)
		assert.NotEmpty(t, registerResponse.Token)

		// 2. Use the token to access a protected endpoint
		token := registerResponse.Token

		// Set up the mock order service to return orders
		mockOrderService.GetOrdersFunc = func(ctx context.Context, brokerName string) ([]core.Order, error) {
			return []core.Order{
				{
					ID:              "order123",
					UserID:          registerResponse.User.ID,
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
			}, nil
		}

		// Make a request to get orders
		req = httptest.NewRequest("GET", "/api/orders?broker=xts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var orders []core.Order
		err = json.Unmarshal(w.Body.Bytes(), &orders)
		require.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, "order123", orders[0].ID)

		// 3. Test portfolio creation
		mockPortfolioService.CreatePortfolioFunc = func(ctx context.Context, portfolio core.Portfolio) (*core.Portfolio, error) {
			portfolio.ID = "portfolio123"
			portfolio.CreatedAt = time.Now()
			portfolio.UpdatedAt = time.Now()
			for i := range portfolio.Legs {
				portfolio.Legs[i].ID = "leg" + string(rune(i+1))
				portfolio.Legs[i].PortfolioID = portfolio.ID
			}
			return &portfolio, nil
		}

		portfolioRequest := core.Portfolio{
			UserID:      registerResponse.User.ID,
			StrategyID:  "strategy123",
			Name:        "API Test Portfolio",
			Symbol:      "NIFTY",
			Exchange:    "NSE",
			Expiry:      "25APR2025",
			DefaultLots: 1,
			IsActive:    true,
			Legs: []core.PortfolioLeg{
				{
					LegID:      1,
					BuySell:    "BUY",
					OptionType: "CE",
					Strike:     "18000",
					Lots:       1,
				},
			},
		}

		portfolioBody, err := json.Marshal(portfolioRequest)
		require.NoError(t, err)

		req = httptest.NewRequest("POST", "/api/portfolios", bytes.NewBuffer(portfolioBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var createdPortfolio core.Portfolio
		err = json.Unmarshal(w.Body.Bytes(), &createdPortfolio)
		require.NoError(t, err)
		assert.Equal(t, "portfolio123", createdPortfolio.ID)
		assert.Equal(t, "API Test Portfolio", createdPortfolio.Name)

		// 4. Test strategy execution
		mockStrategyService.ExecuteStrategyFunc = func(ctx context.Context, strategyID string, userID string) (map[string]*core.ExecutionResponse, error) {
			return map[string]*core.ExecutionResponse{
				"portfolio123": {
					Success:  true,
					OrderIDs: []string{"order123", "order456"},
				},
			}, nil
		}

		req = httptest.NewRequest("POST", "/api/strategies/strategy123/execute", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var executionResponse map[string]*core.ExecutionResponse
		err = json.Unmarshal(w.Body.Bytes(), &executionResponse)
		require.NoError(t, err)
		assert.Contains(t, executionResponse, "portfolio123")
		assert.True(t, executionResponse["portfolio123"].Success)
		assert.Contains(t, executionResponse["portfolio123"].OrderIDs, "order123")
	})

	// Test unauthorized access
	t.Run("UnauthorizedAccess", func(t *testing.T) {
		// Try to access a protected endpoint without a token
		req := httptest.NewRequest("GET", "/api/orders?broker=xts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Try to access a protected endpoint with an invalid token
		req = httptest.NewRequest("GET", "/api/orders?broker=xts", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// Mock services for API integration testing

type MockOrderService struct {
	PlaceOrderFunc   func(ctx context.Context, request core.OrderRequest) (*core.OrderResponse, error)
	ModifyOrderFunc  func(ctx context.Context, orderID string, request core.OrderRequest) (*core.OrderResponse, error)
	CancelOrderFunc  func(ctx context.Context, orderID string, brokerName string) (*core.OrderResponse, error)
	GetOrderFunc     func(ctx context.Context, orderID string, brokerName string) (*core.Order, error)
	GetOrdersFunc    func(ctx context.Context, brokerName string) ([]core.Order, error)
}

func (m *MockOrderService) PlaceOrder(ctx context.Context, request core.OrderRequest) (*core.OrderResponse, error) {
	if m.PlaceOrderFunc != nil {
		return m.PlaceOrderFunc(ctx, request)
	}
	return &core.OrderResponse{Success: true, OrderID: "mock-order-id"}, nil
}

func (m *MockOrderService) ModifyOrder(ctx context.Context, orderID string, request core.OrderRequest) (*core.OrderResponse, error) {
	if m.ModifyOrderFunc != nil {
		return m.ModifyOrderFunc(ctx, orderID, request)
	}
	return &core.OrderResponse{Success: true, OrderID: orderID}, nil
}

func (m *MockOrderService) CancelOrder(ctx context.Context, orderID string, brokerName string) (*core.OrderResponse, error) {
	if m.CancelOrderFunc != nil {
		return m.CancelOrderFunc(ctx, orderID, brokerName)
	}
	return &core.OrderResponse{Success: true, OrderID: orderID}, nil
}

func (m *MockOrderService) GetOrder(ctx context.Context, orderID string, brokerName string) (*core.Order, error) {
	if m.GetOrderFunc != nil {
		return m.GetOrderFunc(ctx, orderID, brokerName)
	}
	return &core.Order{ID: orderID, Status: "COMPLETED"}, nil
}

func (m *MockOrderService) GetOrders(ctx context.Context, brokerName string) ([]core.Order, error) {
	if m.GetOrdersFunc != nil {
		return m.GetOrdersFunc(ctx, brokerName)
	}
	return []core.Order{{ID: "mock-order-id", Status: "COMPLETED"}}, nil
}

type MockPortfolioService struct {
	CreatePortfolioFunc       func(ctx context.Context, portfolio core.Portfolio) (*core.Portfolio, error)
	UpdatePortfolioFunc       func(ctx context.Context, portfolio core.Portfolio) (*core.Portfolio, error)
	GetPortfolioFunc          func(ctx context.Context, portfolioID string) (*core.Portfolio, error)
	DeletePortfolioFunc       func(ctx context.Context, portfolioID string) error
	ExecutePortfolioFunc      func(ctx context.Context, portfolioID string, userID string) (*core.ExecutionResponse, error)
	SquareOffPortfolioFunc    func(ctx context.Context, portfolioID string, userID string) (*core.ExecutionResponse, error)
	GetPortfolioPositionsFunc func(ctx context.Context, portfolioID string) ([]core.Position, error)
}

func (m *MockPortfolioService) CreatePortfolio(ctx context.Context, portfolio core.Portfolio) (*core.Portfolio, error) {
	if m.CreatePortfolioFunc != nil {
		return m.CreatePortfolioFunc(ctx, portfolio)
	}
	portfolio.ID = "mock-portfolio-id"
	return &portfolio, nil
}

func (m *MockPortfolioService) UpdatePortfolio(ctx context.Context, portfolio core.Portfolio) (*core.Portfolio, error) {
	if m.UpdatePortfolioFunc != nil {
		return m.UpdatePortfolioFunc(ctx, portfolio)
	}
	return &portfolio, nil
}

func (m *MockPortfolioService) GetPortfolio(ctx context.Context, portfolioID string) (*core.Portfolio, error) {
	if m.GetPortfolioFunc != nil {
		return m.GetPortfolioFunc(ctx, portfolioID)
	}
	return &core.Portfolio{ID: portfolioID, Name: "Mock Portfolio"}, nil
}

func (m *MockPortfolioService) DeletePortfolio(ctx context.Context, portfolioID string) error {
	if m.DeletePortfolioFunc != nil {
		return m.DeletePortfolioFunc(ctx, portfolioID)
	}
	return nil
}

func (m *MockPortfolioService) ExecutePortfolio(ctx context.Context, portfolioID string, userID string) (*core.ExecutionResponse, error) {
	if m.ExecutePortfolioFunc != nil {
		return m.ExecutePortfolioFunc(ctx, portfolioID, userID)
	}
	return &core.ExecutionResponse{Success: true, OrderIDs: []string{"mock-order-id"}}, nil
}

func (m *MockPortfolioService) SquareOffPortfolio(ctx context.Context, portfolioID string, userID string) (*core.ExecutionResponse, error) {
	if m.SquareOffPortfolioFunc != nil {
		return m.SquareOffPortfolioFunc(ctx, portfolioID, userID)
	}
	return &core.ExecutionResponse{Success: true, OrderIDs: []string{"mock-order-id"}}, nil
}

func (m *MockPortfolioService) GetPortfolioPositions(ctx context.Context, portfolioID string) ([]core.Position, error) {
	if m.GetPortfolioPositionsFunc != nil {
		return m.GetPortfolioPositionsFunc(ctx, portfolioID)
	}
	return []core.Position{{PortfolioID: portfolioID, Symbol: "NIFTY"}}, nil
}

type MockStrategyService struct {
	CreateStrategyFunc            func(ctx context.Context, strategy core.Strategy) (*core.Strategy, error)
	UpdateStrategyFunc            func(ctx context.Context, strategy core.Strategy) (*core.Strategy, error)
	GetStrategyFunc               func(ctx context.Context, strategyID string) (*core.Strategy, error)
	DeleteStrategyFunc            func(ctx context.Context, strategyID string) error
	GetPortfoliosForStrategyFunc  func(ctx context.Context, strategyID string) ([]core.Portfolio, error)
	ExecuteStrategyFunc           func(ctx context.Context, strategyID string, userID string) (map[string]*core.ExecutionResponse, error)
}

func (m *MockStrategyService) CreateStrategy(ctx context.Context, strategy core.Strategy) (*core.Strategy, error) {
	if m.CreateStrategyFunc != nil {
		return m.CreateStrategyFunc(ctx, strategy)
	}
	strategy.ID = "mock-strategy-id"
	return &strategy, nil
}

func (m *MockStrategyService) UpdateStrategy(ctx context.Context, strategy core.Strategy) (*core.Strategy, error) {
	if m.UpdateStrategyFunc != nil {
		return m.UpdateStrategyFunc(ctx, strategy)
	}
	return &strategy, nil
}

func (m *MockStrategyService) GetStrategy(ctx context.Context, strategyID string) (*core.Strategy, error) {
	if m.GetStrategyFunc != nil {
		return m.GetStrategyFunc(ctx, strategyID)
	}
	return &core.Strategy{ID: strategyID, Name: "Mock Strategy"}, nil
}

func (m *MockStrategyService) DeleteStrategy(ctx context.Context, strategyID string) error {
	if m.DeleteStrategyFunc != nil {
		return m.DeleteStrategyFunc(ctx, strategyID)
	}
	return nil
}

func (m *MockStrategyService) GetPortfoliosForStrategy(ctx context.Context, strategyID string) ([]core.Portfolio, error) {
	if m.GetPortfoliosForStrategyFunc != nil {
		return m.GetPortfoliosForStrategyFunc(ctx, strategyID)
	}
	return []core.Portfolio{{ID: "mock-portfolio-id", StrategyID: strategyID}}, nil
}

func (m *MockStrategyService) ExecuteStrategy(ctx context.Context, strategyID string, userID string) (map[string]*core.ExecutionResponse, error) {
	if m.ExecuteStrategyFunc != nil {
		return m.ExecuteStrategyFunc(ctx, strategyID, userID)
	}
	return map[string]*core.ExecutionResponse{
		"mock-portfolio-id": {Success: true, OrderIDs: []string{"mock-order-id"}},
	}, nil
}
