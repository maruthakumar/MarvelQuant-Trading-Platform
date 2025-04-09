package integration

import (
	"context"
	"testing"
	"time"
	"os"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/core"
	"github.com/trading-platform/backend/internal/database"
	"github.com/trading-platform/backend/internal/messagequeue"
)

// TestConfig is imported from auth_order_flow_test.go

// loadConfig is imported from auth_order_flow_test.go

// setupTestEnvironment is imported from auth_order_flow_test.go

func TestPortfolioStrategyFlow(t *testing.T) {
	// Set up the test environment
	dbService, mqService, authService, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a portfolio service
	executionEngine := &MockExecutionEngine{}
	portfolioService := core.NewPortfolioService(executionEngine)

	// Create a strategy service
	strategyService := core.NewStrategyService(portfolioService)

	// Create a risk manager
	riskManager := core.NewRiskManager(portfolioService)

	// Test the complete portfolio and strategy flow
	t.Run("CreateExecuteMonitorPortfolio", func(t *testing.T) {
		ctx := context.Background()

		// 1. Create a strategy
		strategy := core.Strategy{
			UserID:         "user123",
			Name:           "Test Strategy",
			Description:    "Integration test strategy",
			IsActive:       true,
			MaxProfitValue: 1000.0,
			MaxProfitType:  "ABSOLUTE",
			MaxLossValue:   500.0,
			MaxLossType:    "ABSOLUTE",
		}

		createdStrategy, err := strategyService.CreateStrategy(ctx, strategy)
		require.NoError(t, err)
		require.NotNil(t, createdStrategy)
		assert.NotEmpty(t, createdStrategy.ID)
		assert.Equal(t, strategy.Name, createdStrategy.Name)

		// 2. Create a portfolio
		portfolio := core.Portfolio{
			UserID:      "user123",
			StrategyID:  createdStrategy.ID,
			Name:        "Test Portfolio",
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
				{
					LegID:      2,
					BuySell:    "SELL",
					OptionType: "CE",
					Strike:     "18100",
					Lots:       1,
				},
			},
		}

		createdPortfolio, err := portfolioService.CreatePortfolio(ctx, portfolio)
		require.NoError(t, err)
		require.NotNil(t, createdPortfolio)
		assert.NotEmpty(t, createdPortfolio.ID)
		assert.Equal(t, portfolio.Name, createdPortfolio.Name)
		assert.Equal(t, 2, len(createdPortfolio.Legs))

		// 3. Execute the strategy
		executionEngine.ExecuteStrategyFunc = func(ctx context.Context, request core.ExecutionRequest) (*core.ExecutionResponse, error) {
			return &core.ExecutionResponse{
				Success:  true,
				OrderIDs: []string{"order123", "order456"},
			}, nil
		}

		responses, err := strategyService.ExecuteStrategy(ctx, createdStrategy.ID, "user123")
		require.NoError(t, err)
		require.NotNil(t, responses)
		assert.NotEmpty(t, responses)
		for _, response := range responses {
			assert.True(t, response.Success)
			assert.NotEmpty(t, response.OrderIDs)
		}

		// 4. Check portfolio risk
		riskBreached, err := riskManager.CheckPortfolioRisk(ctx, createdPortfolio.ID)
		require.NoError(t, err)
		assert.False(t, riskBreached)

		// 5. Apply risk action (square off)
		executionEngine.ExecuteStrategyFunc = func(ctx context.Context, request core.ExecutionRequest) (*core.ExecutionResponse, error) {
			return &core.ExecutionResponse{
				Success:  true,
				OrderIDs: []string{"order789"},
			}, nil
		}

		err = riskManager.ApplyRiskAction(ctx, createdPortfolio.ID, "SQUARE_OFF")
		require.NoError(t, err)

		// 6. Get portfolio positions
		positions, err := portfolioService.GetPortfolioPositions(ctx, createdPortfolio.ID)
		require.NoError(t, err)
		require.NotNil(t, positions)
		assert.NotEmpty(t, positions)
		assert.Equal(t, createdPortfolio.ID, positions[0].PortfolioID)
	})
}

func TestWebSocketMessageQueueFlow(t *testing.T) {
	// Set up the test environment
	dbService, mqService, authService, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test the WebSocket and message queue integration
	t.Run("PublishSubscribeMessages", func(t *testing.T) {
		ctx := context.Background()

		// 1. Initialize the message queue service
		err := mqService.Initialize()
		require.NoError(t, err)

		// 2. Set up a channel to receive messages
		messageReceived := make(chan bool)
		
		// 3. Subscribe to market data
		err = mqService.SubscribeToMarketData("NIFTY", func(data []byte) {
			// Verify the message content
			assert.Equal(t, []byte("test market data"), data)
			messageReceived <- true
		})
		require.NoError(t, err)

		// 4. Publish market data
		err = mqService.PublishMarketData("NIFTY", []byte("test market data"))
		require.NoError(t, err)

		// 5. Wait for the message to be received or timeout
		select {
		case <-messageReceived:
			// Message was received successfully
		case <-time.After(5 * time.Second):
			t.Fatal("Timed out waiting for message")
		}

		// 6. Set up a channel to receive order messages
		orderReceived := make(chan bool)
		
		// 7. Consume order messages
		err = mqService.ConsumeOrders(func(data []byte) {
			// Verify the message content
			assert.Equal(t, []byte("test order data"), data)
			orderReceived <- true
		})
		require.NoError(t, err)

		// 8. Publish an order message
		err = mqService.PublishOrder("new", []byte("test order data"))
		require.NoError(t, err)

		// 9. Wait for the order message to be received or timeout
		select {
		case <-orderReceived:
			// Order message was received successfully
		case <-time.After(5 * time.Second):
			t.Fatal("Timed out waiting for order message")
		}
	})
}

// MockExecutionEngine is a mock implementation of the execution engine for integration testing
type MockExecutionEngine struct {
	ExecuteStrategyFunc         func(ctx context.Context, request core.ExecutionRequest) (*core.ExecutionResponse, error)
	MonitorExecutionFunc        func(ctx context.Context, portfolioID string) error
	StartExecutionMonitoringFunc func(ctx context.Context) error
}

func (m *MockExecutionEngine) ExecuteStrategy(ctx context.Context, request core.ExecutionRequest) (*core.ExecutionResponse, error) {
	if m.ExecuteStrategyFunc != nil {
		return m.ExecuteStrategyFunc(ctx, request)
	}
	return &core.ExecutionResponse{
		Success:  true,
		OrderIDs: []string{"mock-order-id"},
	}, nil
}

func (m *MockExecutionEngine) MonitorExecution(ctx context.Context, portfolioID string) error {
	if m.MonitorExecutionFunc != nil {
		return m.MonitorExecutionFunc(ctx, portfolioID)
	}
	return nil
}

func (m *MockExecutionEngine) StartExecutionMonitoring(ctx context.Context) error {
	if m.StartExecutionMonitoringFunc != nil {
		return m.StartExecutionMonitoringFunc(ctx)
	}
	return nil
}
