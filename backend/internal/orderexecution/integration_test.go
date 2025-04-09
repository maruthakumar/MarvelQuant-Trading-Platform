package orderexecution

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// IntegrationTest represents an integration test for the order execution engine
type IntegrationTest struct {
	Name        string
	Description string
	Setup       func(t *testing.T) (context.Context, *HighPerformanceOrderExecutionEngine, func())
	Test        func(t *testing.T, ctx context.Context, engine *HighPerformanceOrderExecutionEngine)
}

// RunIntegrationTests runs all integration tests
func RunIntegrationTests(t *testing.T) {
	tests := []IntegrationTest{
		{
			Name:        "BasicOrderExecution",
			Description: "Test basic order execution flow",
			Setup:       setupBasicOrderExecutionTest,
			Test:        testBasicOrderExecution,
		},
		{
			Name:        "OrderLifecycle",
			Description: "Test order lifecycle management",
			Setup:       setupOrderLifecycleTest,
			Test:        testOrderLifecycle,
		},
		{
			Name:        "ErrorHandling",
			Description: "Test error handling and recovery",
			Setup:       setupErrorHandlingTest,
			Test:        testErrorHandling,
		},
		{
			Name:        "RiskManagement",
			Description: "Test risk management system",
			Setup:       setupRiskManagementTest,
			Test:        testRiskManagement,
		},
		{
			Name:        "BracketOrders",
			Description: "Test bracket order execution",
			Setup:       setupBracketOrderTest,
			Test:        testBracketOrder,
		},
		{
			Name:        "HighFrequencyTrading",
			Description: "Test high-frequency trading performance",
			Setup:       setupHighFrequencyTradingTest,
			Test:        testHighFrequencyTrading,
		},
		{
			Name:        "MonitoringAndLogging",
			Description: "Test monitoring and logging",
			Setup:       setupMonitoringAndLoggingTest,
			Test:        testMonitoringAndLogging,
		},
		{
			Name:        "ConcurrentOrders",
			Description: "Test concurrent order execution",
			Setup:       setupConcurrentOrdersTest,
			Test:        testConcurrentOrders,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctx, engine, cleanup := test.Setup(t)
			defer cleanup()
			test.Test(t, ctx, engine)
		})
	}
}

// setupBasicOrderExecutionTest sets up a basic order execution test
func setupBasicOrderExecutionTest(t *testing.T) (context.Context, *HighPerformanceOrderExecutionEngine, func()) {
	// Create logger
	logger := NewConsoleLogger()
	logger.SetPrefix("BasicOrderExecutionTest")

	// Create mock broker
	broker := NewMockBroker("XTS", logger)

	// Create order execution engine
	engine := createTestExecutionEngine(broker, logger)

	// Create context
	ctx := context.Background()

	// Return cleanup function
	cleanup := func() {
		engine.Close()
	}

	return ctx, engine, cleanup
}

// testBasicOrderExecution tests basic order execution
func testBasicOrderExecution(t *testing.T, ctx context.Context, engine *HighPerformanceOrderExecutionEngine) {
	// Create test order
	order := Order{
		ID:          "TEST-ORDER-001",
		PortfolioID: "PORTFOLIO-001",
		StrategyID:  "STRATEGY-001",
		Symbol:      "RELIANCE",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    100,
		Price:       2500.0,
		Status:      OrderStatusNew,
	}

	// Execute order
	response, err := engine.ExecuteOrder(ctx, order)
	if err != nil {
		t.Fatalf("Failed to execute order: %v", err)
	}

	// Verify response
	if response.Order.ID != order.ID {
		t.Errorf("Expected order ID %s, got %s", order.ID, response.Order.ID)
	}
	if response.Status != "SUCCESS" && response.Status != "BATCHED" {
		t.Errorf("Expected status SUCCESS or BATCHED, got %s", response.Status)
	}

	// If order was batched, wait for processing
	if response.Status == "BATCHED" {
		time.Sleep(500 * time.Millisecond)
	}

	// Get order status
	status, err := engine.GetOrderStatus(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order status: %v", err)
	}

	// Verify status
	if status != OrderLifecycleStateCompleted && status != OrderLifecycleStateAcknowledged {
		t.Errorf("Expected status COMPLETED or ACKNOWLEDGED, got %s", status)
	}

	// Get order details
	details, err := engine.GetOrderDetails(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order details: %v", err)
	}

	// Verify details
	if details.Order.ID != order.ID {
		t.Errorf("Expected order ID %s, got %s", order.ID, details.Order.ID)
	}
	if details.Order.Symbol != order.Symbol {
		t.Errorf("Expected symbol %s, got %s", order.Symbol, details.Order.Symbol)
	}
	if details.Order.Quantity != order.Quantity {
		t.Errorf("Expected quantity %d, got %d", order.Quantity, details.Order.Quantity)
	}
	if details.Order.Price != order.Price {
		t.Errorf("Expected price %f, got %f", order.Price, details.Order.Price)
	}
}

// setupOrderLifecycleTest sets up an order lifecycle test
func setupOrderLifecycleTest(t *testing.T) (context.Context, *HighPerformanceOrderExecutionEngine, func()) {
	// Create logger
	logger := NewConsoleLogger()
	logger.SetPrefix("OrderLifecycleTest")

	// Create mock broker
	broker := NewMockBroker("XTS", logger)

	// Create order execution engine
	engine := createTestExecutionEngine(broker, logger)

	// Create context
	ctx := context.Background()

	// Return cleanup function
	cleanup := func() {
		engine.Close()
	}

	return ctx, engine, cleanup
}

// testOrderLifecycle tests order lifecycle management
func testOrderLifecycle(t *testing.T, ctx context.Context, engine *HighPerformanceOrderExecutionEngine) {
	// Create test order
	order := Order{
		ID:          "TEST-ORDER-002",
		PortfolioID: "PORTFOLIO-001",
		StrategyID:  "STRATEGY-001",
		Symbol:      "INFY",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    50,
		Price:       1500.0,
		Status:      OrderStatusNew,
	}

	// Execute order
	response, err := engine.ExecuteOrder(ctx, order)
	if err != nil {
		t.Fatalf("Failed to execute order: %v", err)
	}

	// Verify response
	if response.Order.ID != order.ID {
		t.Errorf("Expected order ID %s, got %s", order.ID, response.Order.ID)
	}

	// If order was batched, wait for processing
	if response.Status == "BATCHED" {
		time.Sleep(500 * time.Millisecond)
	}

	// Get order status
	status, err := engine.GetOrderStatus(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order status: %v", err)
	}

	// Verify initial status
	if status != OrderLifecycleStateAcknowledged && status != OrderLifecycleStateCompleted {
		t.Errorf("Expected status ACKNOWLEDGED or COMPLETED, got %s", status)
	}

	// If order is still in progress, cancel it
	if status == OrderLifecycleStateAcknowledged {
		// Cancel order
		err = engine.CancelOrder(ctx, order.ID)
		if err != nil {
			t.Fatalf("Failed to cancel order: %v", err)
		}

		// Get updated status
		status, err = engine.GetOrderStatus(ctx, order.ID)
		if err != nil {
			t.Fatalf("Failed to get order status: %v", err)
		}

		// Verify cancelled status
		if status != OrderLifecycleStateCancelled {
			t.Errorf("Expected status CANCELLED, got %s", status)
		}
	}

	// Get order details
	details, err := engine.GetOrderDetails(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order details: %v", err)
	}

	// Verify events
	if len(details.Events) == 0 {
		t.Errorf("Expected events, got none")
	}

	// Verify event types
	eventTypes := make(map[string]bool)
	for _, event := range details.Events {
		eventTypes[event.Type] = true
	}

	// Check for required events
	if !eventTypes["ORDER_CREATED"] {
		t.Errorf("Missing ORDER_CREATED event")
	}
	if !eventTypes["ORDER_SUBMITTED"] {
		t.Errorf("Missing ORDER_SUBMITTED event")
	}
	if status == OrderLifecycleStateCancelled && !eventTypes["ORDER_CANCELLED"] {
		t.Errorf("Missing ORDER_CANCELLED event")
	}
	if status == OrderLifecycleStateCompleted && !eventTypes["ORDER_COMPLETED"] {
		t.Errorf("Missing ORDER_COMPLETED event")
	}
}

// setupErrorHandlingTest sets up an error handling test
func setupErrorHandlingTest(t *testing.T) (context.Context, *HighPerformanceOrderExecutionEngine, func()) {
	// Create logger
	logger := NewConsoleLogger()
	logger.SetPrefix("ErrorHandlingTest")

	// Create mock broker with error simulation
	broker := NewMockBroker("XTS", logger)
	broker.SetErrorSimulation(true)
	broker.SetErrorRate(50) // 50% error rate

	// Create order execution engine
	engine := createTestExecutionEngine(broker, logger)

	// Create context
	ctx := context.Background()

	// Return cleanup function
	cleanup := func() {
		engine.Close()
	}

	return ctx, engine, cleanup
}

// testErrorHandling tests error handling and recovery
func testErrorHandling(t *testing.T, ctx context.Context, engine *HighPerformanceOrderExecutionEngine) {
	// Create test orders
	orders := []Order{
		{
			ID:          "ERROR-TEST-001",
			PortfolioID: "PORTFOLIO-001",
			StrategyID:  "STRATEGY-001",
			Symbol:      "TCS",
			Exchange:    "NSE",
			OrderType:   OrderTypeLimit,
			ProductType: ProductTypeDelivery,
			Side:        OrderSideBuy,
			Quantity:    25,
			Price:       3500.0,
			Status:      OrderStatusNew,
		},
		{
			ID:          "ERROR-TEST-002",
			PortfolioID: "PORTFOLIO-001",
			StrategyID:  "STRATEGY-001",
			Symbol:      "WIPRO",
			Exchange:    "NSE",
			OrderType:   OrderTypeLimit,
			ProductType: ProductTypeDelivery,
			Side:        OrderSideBuy,
			Quantity:    100,
			Price:       450.0,
			Status:      OrderStatusNew,
		},
		{
			ID:          "ERROR-TEST-003",
			PortfolioID: "PORTFOLIO-001",
			StrategyID:  "STRATEGY-001",
			Symbol:      "HDFCBANK",
			Exchange:    "NSE",
			OrderType:   OrderTypeLimit,
			ProductType: ProductTypeDelivery,
			Side:        OrderSideBuy,
			Quantity:    50,
			Price:       1600.0,
			Status:      OrderStatusNew,
		},
	}

	// Execute orders
	var successCount, errorCount int
	for _, order := range orders {
		_, err := engine.ExecuteOrder(ctx, order)
		if err != nil {
			// Error should be of type ExecutionError
			execErr, ok := err.(*ExecutionError)
			if !ok {
				t.Errorf("Expected ExecutionError, got %T", err)
			} else {
				// Verify error fields
				if execErr.Type == "" {
					t.Errorf("Error type is empty")
				}
				if execErr.Code == 0 {
					t.Errorf("Error code is 0")
				}
				if execErr.Message == "" {
					t.Errorf("Error message is empty")
				}
			}
			errorCount++
		} else {
			successCount++
		}
	}

	// Verify that some orders succeeded and some failed
	if successCount == 0 {
		t.Errorf("Expected some orders to succeed, but all failed")
	}
	if errorCount == 0 {
		t.Errorf("Expected some orders to fail, but all succeeded")
	}

	// Wait for any batched orders to process
	time.Sleep(500 * time.Millisecond)

	// Get metrics
	metrics := engine.GetMetrics()
	
	// Verify metrics
	if metrics["errorCount"].(int64) == 0 {
		t.Errorf("Expected non-zero error count in metrics")
	}
}

// setupRiskManagementTest sets up a risk management test
func setupRiskManagementTest(t *testing.T) (context.Context, *HighPerformanceOrderExecutionEngine, func()) {
	// Create logger
	logger := NewConsoleLogger()
	logger.SetPrefix("RiskManagementTest")

	// Create mock broker
	broker := NewMockBroker("XTS", logger)

	// Create order execution engine with strict risk limits
	engine := createTestExecutionEngineWithRiskLimits(broker, logger)

	// Create context
	ctx := context.Background()

	// Return cleanup function
	cleanup := func() {
		engine.Close()
	}

	return ctx, engine, cleanup
}

// testRiskManagement tests the risk management system
func testRiskManagement(t *testing.T, ctx context.Context, engine *HighPerformanceOrderExecutionEngine) {
	// Create valid order (within risk limits)
	validOrder := Order{
		ID:          "RISK-TEST-001",
		PortfolioID: "PORTFOLIO-001",
		StrategyID:  "STRATEGY-001",
		Symbol:      "SBIN",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    100,
		Price:       500.0, // Total value: 50,000
		Status:      OrderStatusNew,
	}

	// Create invalid order (exceeds position limit)
	invalidPositionOrder := Order{
		ID:          "RISK-TEST-002",
		PortfolioID: "PORTFOLIO-001",
		StrategyID:  "STRATEGY-001",
		Symbol:      "SBIN",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    1000,
		Price:       500.0, // Total value: 500,000 (exceeds position limit)
		Status:      OrderStatusNew,
	}

	// Create invalid order (exceeds order value limit)
	invalidValueOrder := Order{
		ID:          "RISK-TEST-003",
		PortfolioID: "PORTFOLIO-001",
		StrategyID:  "STRATEGY-001",
		Symbol:      "RELIANCE",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    500,
		Price:       2500.0, // Total value: 1,250,000 (exceeds order value limit)
		Status:      OrderStatusNew,
	}

	// Execute valid order
	_, err := engine.ExecuteOrder(ctx, validOrder)
	if err != nil {
		t.Errorf("Valid order should not be rejected: %v", err)
	}

	// Execute invalid position order
	_, err = engine.ExecuteOrder(ctx, invalidPositionOrder)
	if err == nil {
		t.Errorf("Order exceeding position limit should be rejected")
	} else {
		// Verify error type
		execErr, ok := err.(*ExecutionError)
		if !ok {
			t.Errorf("Expected ExecutionError, got %T", err)
		} else if execErr.Type != ErrorTypeRiskManagement {
			t.Errorf("Expected risk management error, got %s", execErr.Type)
		}
	}

	// Execute invalid value order
	_, err = engine.ExecuteOrder(ctx, invalidValueOrder)
	if err == nil {
		t.Errorf("Order exceeding value limit should be rejected")
	} else {
		// Verify error type
		execErr, ok := err.(*ExecutionError)
		if !ok {
			t.Errorf("Expected ExecutionError, got %T", err)
		} else if execErr.Type != ErrorTypeRiskManagement {
			t.Errorf("Expected risk management error, got %s", execErr.Type)
		}
	}
}

// setupBracketOrderTest sets up a bracket order test
func setupBracketOrderTest(t *testing.T) (context.Context, *HighPerformanceOrderExecutionEngine, func()) {
	// Create logger
	logger := NewConsoleLogger()
	logger.SetPrefix("BracketOrderTest")

	// Create mock broker
	broker := NewMockBroker("XTS", logger)

	// Create order execution engine
	engine := createTestExecutionEngine(broker, logger)

	// Create context
	ctx := context.Background()

	// Return cleanup function
	cleanup := func() {
		engine.Close()
	}

	return ctx, engine, cleanup
}

// testBracketOrder tests bracket order execution
func testBracketOrder(t *testing.T, ctx context.Context, engine *HighPerformanceOrderExecutionEngine) {
	// Create main order
	mainOrder := Order{
		ID:          "BRACKET-TEST-001",
		PortfolioID: "PORTFOLIO-001",
		StrategyID:  "STRATEGY-001",
		Symbol:      "ICICIBANK",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    100,
		Price:       900.0,
		Status:      OrderStatusNew,
	}

	// Create bracket order
	responses, err := engine.CreateBracketOrder(ctx, mainOrder, 950.0, 850.0)
	if err != nil {
		t.Fatalf("Failed to create bracket order: %v", err)
	}

	// Verify responses
	if len(responses) != 3 {
		t.Errorf("Expected 3 orders in bracket, got %d", len(responses))
	}

	// Wait for any batched orders to process
	time.Sleep(500 * time.Millisecond)

	// Verify order types
	var mainFound, tpFound, slFound bool
	for _, response := range responses {
		// Get order details
		details, err := engine.GetOrderDetails(ctx, response.Order.ID)
		if err != nil {
			t.Fatalf("Failed to get order details: %v", err)
		}

		// Check order type
		if details.Order.ID == mainOrder.ID {
			mainFound = true
		} else if details.Order.OrderType == OrderTypeLimit && details.Order.Side == OrderSideSell && details.Order.Price == 950.0 {
			tpFound = true
		} else if details.Order.OrderType == OrderTypeStopLoss && details.Order.Side == OrderSideSell && details.Order.Price == 850.0 {
			slFound = true
		}

(Content truncated due to size limit. Use line ranges to read in chunks)