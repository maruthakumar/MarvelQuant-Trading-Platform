package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"encoding/json"
	
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	"trading_platform/backend/internal/interfaces"
	"trading_platform/backend/internal/models"
)

// TestAPIGatewayIntegration tests the API Gateway integration with HTTP handlers
func TestAPIGatewayIntegration(t *testing.T) {
	// Create mock services
	mockExecutionPlatform := new(MockExecutionPlatform)
	
	// Create API Gateway with mock execution platform
	gateway := NewAPIGateway(mockExecutionPlatform)
	
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add middleware to inject user context
	router.Use(func(c *gin.Context) {
		// For testing, we'll inject a test user context
		ctx := context.WithValue(c.Request.Context(), "userID", "user123")
		ctx = context.WithValue(ctx, "userType", "SIM")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	
	// Add permissions for test user
	gateway.accessControlList["user123"] = []string{
		"simulation:account:read",
		"simulation:market:read",
		"backtest:session:create",
		"system:status:read",
	}
	
	// Setup routes
	router.GET("/api/simulation/accounts/:id", func(c *gin.Context) {
		accountID := c.Param("id")
		account, err := gateway.GetSimulationAccount(c.Request.Context(), accountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, account)
	})
	
	router.GET("/api/simulation/market/:symbol/price", func(c *gin.Context) {
		symbol := c.Param("symbol")
		price, err := gateway.GetCurrentMarketPrice(c.Request.Context(), symbol)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, price)
	})
	
	router.POST("/api/simulation/backtest", func(c *gin.Context) {
		var session models.BacktestSession
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		accountID := c.Query("accountId")
		result, err := gateway.CreateBacktestSession(c.Request.Context(), accountID, session)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})
	
	router.GET("/api/system/status", func(c *gin.Context) {
		status, err := gateway.GetSystemStatus(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, status)
	})
	
	// Setup mock responses
	mockAccount := &models.SimulationAccount{
		ID:             "sim123",
		Name:           "Test Account",
		InitialBalance: 100000.0,
		CurrentBalance: 100000.0,
		Currency:       "USD",
		SimulationType: "PAPER",
		IsActive:       true,
	}
	
	mockMarketData := &models.MarketDataSnapshot{
		Symbol:    "AAPL",
		Price:     150.25,
		Timestamp: time.Now(),
	}
	
	mockBacktestSession := &models.BacktestSession{
		ID:             "session123",
		Name:           "Test Backtest",
		StartDate:      time.Now().Add(-30 * 24 * time.Hour),
		EndDate:        time.Now(),
		Symbols:        []string{"AAPL", "MSFT"},
		InitialBalance: 100000.0,
		Status:         "PENDING",
	}
	
	mockSystemStatus := map[string]interface{}{
		"simulation_system": map[string]interface{}{
			"status":    "operational",
			"version":   "9.6.3",
			"uptime":    "3d 12h 45m",
			"load":      0.42,
			"memory":    "68%",
		},
		"execution_platform": map[string]interface{}{
			"status":    "operational",
			"version":   "9.6.3",
			"uptime":    "5d 8h 30m",
			"load":      0.35,
			"memory":    "72%",
		},
	}
	
	// Replace services with mocks
	mockSimulationService := new(MockSimulationService)
	mockMarketSimulationService := new(MockMarketSimulationService)
	mockBacktestService := new(MockBacktestService)
	
	gateway.simulationService = mockSimulationService
	gateway.marketSimulationService = mockMarketSimulationService
	gateway.backtestService = mockBacktestService
	
	// Setup mock expectations
	mockSimulationService.On("GetSimulationAccount", "sim123").Return(mockAccount, nil)
	mockMarketSimulationService.On("GetCurrentMarketPrice", "AAPL").Return(mockMarketData, nil)
	mockBacktestService.On("CreateBacktestSession", "sim123", mock.AnythingOfType("models.BacktestSession")).Return(mockBacktestSession, nil)
	
	// Setup market data synchronization
	symbols := []string{"AAPL", "MSFT", "GOOGL"}
	mockExecutionPlatform.On("SynchronizeMarketData", mock.Anything, symbols).Return(nil)
	
	t.Run("GetSimulationAccount", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/simulation/accounts/sim123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.SimulationAccount
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockAccount.ID, response.ID)
		assert.Equal(t, mockAccount.Name, response.Name)
		assert.Equal(t, mockAccount.InitialBalance, response.InitialBalance)
		assert.Equal(t, mockAccount.CurrentBalance, response.CurrentBalance)
		assert.Equal(t, mockAccount.Currency, response.Currency)
		assert.Equal(t, mockAccount.SimulationType, response.SimulationType)
		assert.Equal(t, mockAccount.IsActive, response.IsActive)
	})
	
	t.Run("GetCurrentMarketPrice", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/simulation/market/AAPL/price", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.MarketDataSnapshot
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockMarketData.Symbol, response.Symbol)
		assert.Equal(t, mockMarketData.Price, response.Price)
	})
	
	t.Run("GetSystemStatus", func(t *testing.T) {
		// Override the GetSystemStatus method to return our mock data
		originalGetSystemStatus := gateway.GetSystemStatus
		defer func() { gateway.GetSystemStatus = originalGetSystemStatus }()
		
		gateway.GetSystemStatus = func(ctx context.Context) (map[string]interface{}, error) {
			return mockSystemStatus, nil
		}
		
		req, _ := http.NewRequest("GET", "/api/system/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Check that the response contains the expected keys
		simSystem, ok := response["simulation_system"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "operational", simSystem["status"])
		
		execPlatform, ok := response["execution_platform"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "operational", execPlatform["status"])
	})
	
	t.Run("Permission Denied", func(t *testing.T) {
		// Create a new router with middleware that injects an unauthorized user
		unauthorizedRouter := gin.New()
		unauthorizedRouter.Use(func(c *gin.Context) {
			ctx := context.WithValue(c.Request.Context(), "userID", "unauthorized")
			ctx = context.WithValue(ctx, "userType", "SIM")
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		})
		
		unauthorizedRouter.GET("/api/simulation/accounts/:id", func(c *gin.Context) {
			accountID := c.Param("id")
			account, err := gateway.GetSimulationAccount(c.Request.Context(), accountID)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, account)
		})
		
		req, _ := http.NewRequest("GET", "/api/simulation/accounts/sim123", nil)
		w := httptest.NewRecorder()
		unauthorizedRouter.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "authorization")
	})
	
	t.Run("Rate Limit", func(t *testing.T) {
		// Set up a rate limit that will be exceeded
		gateway.rateLimits["account_management"] = RateLimit{
			MaxRequests:     1,
			TimeWindow:      time.Minute,
			CurrentRequests: make(map[string][]time.Time),
		}
		
		// First request should succeed
		req1, _ := http.NewRequest("GET", "/api/simulation/accounts/sim123", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)
		
		// Second request should fail due to rate limit
		req2, _ := http.NewRequest("GET", "/api/simulation/accounts/sim123", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)
		assert.Contains(t, w2.Body.String(), "rate limit")
	})
}

// TestCrossSystemErrorHandling tests the error handling across system boundaries
func TestCrossSystemErrorHandling(t *testing.T) {
	// Create mock services
	mockExecutionPlatform := new(MockExecutionPlatform)
	
	// Create API Gateway with mock execution platform
	gateway := NewAPIGateway(mockExecutionPlatform)
	
	// Add permissions for test user
	gateway.accessControlList["user123"] = []string{
		"simulation:market:read",
		"system:sync:execute",
	}
	
	// Create context with user ID and user type
	ctx := context.WithValue(context.Background(), "userID", "user123")
	ctx = context.WithValue(ctx, "userType", "SIM")
	
	t.Run("ExecutionPlatformError", func(t *testing.T) {
		// Setup mock to return an error
		mockExecutionPlatform.On("GetRealTimeMarketData", ctx, "AAPL").Return(nil, errors.New("execution platform error"))
		
		// Setup market data synchronization to fail
		symbols := []string{"AAPL"}
		mockExecutionPlatform.On("SynchronizeMarketData", ctx, symbols).Return(errors.New("synchronization failed"))
		
		// Call method that depends on execution platform
		_, err := gateway.GetCurrentMarketPrice(ctx, "AAPL")
		
		// Assert that the error is properly handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal system error")
		
		// The original error should be logged but not exposed to the caller
		assert.NotContains(t, err.Error(), "execution platform error")
	})
	
	t.Run("SimulationSystemError", func(t *testing.T) {
		// Replace services with mocks
		mockMarketSimulationService := new(MockMarketSimulationService)
		gateway.marketSimulationService = mockMarketSimulationService
		
		// Setup mock to return an error
		mockMarketSimulationService.On("GetCurrentMarketPrice", "AAPL").Return(nil, errors.New("simulation system error"))
		
		// Setup market data synchronization to succeed
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		mockExecutionPlatform.On("SynchronizeMarketData", ctx, symbols).Return(nil)
		
		// Call method
		_, err := gateway.GetCurrentMarketPrice(ctx, "AAPL")
		
		// Assert that the error is properly handled
		assert.Error(t, err)
		
		// Validation errors should be passed through with details
		assert.Contains(t, err.Error(), "simulation system error")
	})
	
	t.Run("DataSynchronizationError", func(t *testing.T) {
		// Setup market data synchronization to fail
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		mockExecutionPlatform.On("SynchronizeMarketData", ctx, symbols).Return(errors.New("synchronization failed"))
		
		// Call method that requires data synchronization
		err := gateway.SynchronizeMarketData(ctx, symbols)
		
		// Assert that the error is properly handled
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal system error")
	})
}

// TestDataSynchronization tests the data synchronization between execution platform and simulation system
func TestDataSynchronization(t *testing.T) {
	// Create mock services
	mockExecutionPlatform := new(MockExecutionPlatform)
	
	// Create API Gateway with mock execution platform
	gateway := NewAPIGateway(mockExecutionPlatform)
	
	// Replace services with mocks
	mockMarketSimulationService := new(MockMarketSimulationService)
	gateway.marketSimulationService = mockMarketSimulationService
	
	// Add permissions for test user
	gateway.accessControlList["user123"] = []string{
		"simulation:market:read",
		"system:sync:execute",
	}
	
	// Create context with user ID and user type
	ctx := context.WithValue(context.Background(), "userID", "user123")
	ctx = context.WithValue(ctx, "userType", "SIM")
	
	t.Run("InitialSynchronization", func(t *testing.T) {
		// Setup mocks
		symbols := []string{"AAPL", "MSFT"}
		
		marketData1 := &models.MarketDataSnapshot{
			Symbol:    "AAPL",
			Price:     150.25,
			Timestamp: time.Now(),
		}
		
		marketData2 := &models.MarketDataSnapshot{
			Symbol:    "MSFT",
			Price:     290.75,
			Timestamp: time.Now(),
		}
		
		mockExecutionPlatform.On("GetRealTimeMarketData", ctx, "AAPL").Return(marketData1, nil)
		mockExecutionPlatform.On("GetRealTimeMarketData", ctx, "MSFT").Return(marketData2, nil)
		
		mockMarketSimulationService.On("UpdateMarketData", "AAPL", marketData1).Return(nil)
		mockMarketSimulationService.On("UpdateMarketData", "MSFT", marketData2).Return(nil)
		
		// Call method
		err := gateway.SynchronizeMarketData(ctx, symbols)
		
		// Assert
		assert.NoError(t, err)
		mockExecutionPlatform.AssertExpectations(t)
		mockMarketSimulationService.AssertExpectations(t)
	})
	
	t.Run("CachedSynchronization", func(t *testing.T) {
		// Setup initial synchronization time
		gateway.lastSyncTime["market_data"] = time.Now()
		
		// Setup mock for market data
		mockMarketData := &models.MarketDataSnapshot{
			Symbol:    "AAPL",
			Price:     150.25,
			Timestamp: time.Now(),
		}
		
		mockMarketSimulationService.On("GetCurrentMarketPrice", "AAPL").Return(mockMarketData, nil)
		
		// Call method that would normally trigger synchronization
		result, err := gateway.GetCurrentMarketPrice(ctx, "AAPL")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, mockMarketData, result)
		
		// Synchronization should not be called again since it was recently done
		mockExecutionPlatform.AssertNotCalled(t, "SynchronizeMarketData")
	})
	
	t.Run("ExpiredCacheSynchronization", func(t *testing.T) {
		// Setup expired synchronization time (more than 5 minutes ago)
		gateway.lastSyncTime["market_data"] = time.Now().Add(-10 * time.Minute)
		
		// Setup mocks
		symbols := []string{"AAPL", "MSFT", "GOOGL"}
		mockExecutionPlatform.On("SynchronizeMarketData", ctx, symbols).Return(nil)
		
		mockMarketData := &models.MarketDataSnapshot{
			Symbol:    "AAPL",
			Price:     150.25,
			Timestamp: time.Now(),
		}
		
		mockMarketSimulationService.On("GetCurrentMarketPrice", "AAPL").Return(mockMarketData, nil)
		
		// Call method that would trigger synchronization
		result, err := gateway.GetCurrentMarketPrice(ctx, "AAPL")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, mockMarketData, result)
		
		// Synchronization should be called again since cache expired
		mockExecutionPlatform.AssertCalled(t, "SynchronizeMarketData", ctx, symbols)
	})
}
