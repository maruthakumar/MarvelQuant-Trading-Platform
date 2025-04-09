package orderexecution

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

// TestOrderExecutionEngine tests the order execution engine
func TestOrderExecutionEngine(t *testing.T) {
	// Create a mock broker adapter
	mockBroker := NewMockBrokerAdapter()
	
	// Create a smart router
	smartRouter := NewDefaultSmartRouter(BestPrice)
	
	// Create an order execution engine
	engine := NewOrderExecutionEngine(smartRouter)
	
	// Register the mock broker
	engine.RegisterBroker("MOCK", mockBroker)
	
	// Create an order request
	request := &OrderRequest{
		Symbol:          "RELIANCE-EQ",
		Quantity:        100,
		Price:           2500.0,
		OrderType:       Limit,
		TransactionType: Buy,
		Validity:        Day,
		Exchange:        "NSE",
		Product:         Normal,
	}
	
	// Execute the order
	ctx := context.Background()
	response, err := engine.ExecuteOrder(ctx, request)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error executing order: %v", err)
		return
	}
	
	// Check response
	if !response.Status {
		t.Errorf("Order execution failed: %s", response.Error)
		return
	}
	
	// Check order ID
	if response.Order.ID == "" {
		t.Errorf("Order ID is empty")
		return
	}
	
	// Check order status
	if response.Order.Status != Open {
		t.Errorf("Expected order status %s, got %s", Open, response.Order.Status)
		return
	}
	
	// Test order modification
	modifyRequest := &OrderRequest{
		Symbol:          "RELIANCE-EQ",
		Quantity:        50,
		Price:           2550.0,
		OrderType:       Limit,
		TransactionType: Buy,
		Validity:        Day,
		Exchange:        "NSE",
		Product:         Normal,
	}
	
	modifyResponse, err := engine.ModifyOrder(ctx, response.Order.ID, modifyRequest)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error modifying order: %v", err)
		return
	}
	
	// Check response
	if !modifyResponse.Status {
		t.Errorf("Order modification failed: %s", modifyResponse.Error)
		return
	}
	
	// Check modified order
	if modifyResponse.Order.Quantity != 50 {
		t.Errorf("Expected modified quantity 50, got %d", modifyResponse.Order.Quantity)
		return
	}
	
	// Test order cancellation
	cancelResponse, err := engine.CancelOrder(ctx, response.Order.ID)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error cancelling order: %v", err)
		return
	}
	
	// Check response
	if !cancelResponse.Status {
		t.Errorf("Order cancellation failed: %s", cancelResponse.Error)
		return
	}
	
	// Check cancelled order
	if cancelResponse.Order.Status != Cancelled {
		t.Errorf("Expected order status %s, got %s", Cancelled, cancelResponse.Order.Status)
		return
	}
	
	fmt.Println("Order execution engine tests passed")
}

// TestSmartRouter tests the smart router
func TestSmartRouter(t *testing.T) {
	// Create mock broker adapters
	mockBroker1 := NewMockBrokerAdapter()
	mockBroker2 := NewMockBrokerAdapter()
	
	// Create a smart router
	smartRouter := NewDefaultSmartRouter(BestPrice)
	
	// Register the mock brokers
	smartRouter.RegisterBroker("MOCK1", mockBroker1)
	smartRouter.RegisterBroker("MOCK2", mockBroker2)
	
	// Create an order request
	request := &OrderRequest{
		Symbol:          "RELIANCE-EQ",
		Quantity:        100,
		Price:           2500.0,
		OrderType:       Limit,
		TransactionType: Buy,
		Validity:        Day,
		Exchange:        "NSE",
		Product:         Normal,
	}
	
	// Route the order
	ctx := context.Background()
	broker, err := smartRouter.RouteOrder(ctx, request)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error routing order: %v", err)
		return
	}
	
	// Check broker
	if broker == nil {
		t.Errorf("Broker is nil")
		return
	}
	
	// Test different routing strategies
	strategies := []SmartRoutingStrategy{
		BestPrice,
		LowestLatency,
		HighestFillRate,
		LowestCost,
		RoundRobin,
		VolumeWeighted,
	}
	
	for _, strategy := range strategies {
		smartRouter.defaultStrategy = strategy
		broker, err := smartRouter.RouteOrder(ctx, request)
		
		// Check for errors
		if err != nil {
			t.Errorf("Error routing order with strategy %s: %v", strategy, err)
			continue
		}
		
		// Check broker
		if broker == nil {
			t.Errorf("Broker is nil with strategy %s", strategy)
			continue
		}
	}
	
	fmt.Println("Smart router tests passed")
}

// TestOrderPlacementService tests the order placement service
func TestOrderPlacementService(t *testing.T) {
	// Create a mock broker adapter
	mockBroker := NewMockBrokerAdapter()
	
	// Create a smart router
	smartRouter := NewDefaultSmartRouter(BestPrice)
	
	// Create an order execution engine
	engine := NewOrderExecutionEngine(smartRouter)
	
	// Register the mock broker
	engine.RegisterBroker("MOCK", mockBroker)
	
	// Create a strategy manager
	strategyManager := NewStrategyManager()
	
	// Create an order placement service
	placementService := NewOrderPlacementService(engine, strategyManager, 100, 5, 100)
	
	// Start the service
	placementService.Start()
	defer placementService.Stop()
	
	// Create an order request
	request := &OrderRequest{
		Symbol:          "RELIANCE-EQ",
		Quantity:        100,
		Price:           2500.0,
		OrderType:       Limit,
		TransactionType: Buy,
		Validity:        Day,
		Exchange:        "NSE",
		Product:         Normal,
	}
	
	// Test synchronous order placement
	ctx := context.Background()
	response, err := placementService.PlaceOrderSync(ctx, request)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error placing order synchronously: %v", err)
		return
	}
	
	// Check response
	if !response.Status {
		t.Errorf("Order placement failed: %s", response.Error)
		return
	}
	
	// Test asynchronous order placement
	orderPlaced := make(chan bool)
	err = placementService.PlaceOrder(request, func(response *OrderResponse, err error) {
		if err != nil {
			t.Errorf("Error in async order callback: %v", err)
			orderPlaced <- false
			return
		}
		
		if !response.Status {
			t.Errorf("Async order placement failed: %s", response.Error)
			orderPlaced <- false
			return
		}
		
		orderPlaced <- true
	})
	
	// Check for errors
	if err != nil {
		t.Errorf("Error placing order asynchronously: %v", err)
		return
	}
	
	// Wait for order to be placed
	select {
	case success := <-orderPlaced:
		if !success {
			return
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout waiting for async order placement")
		return
	}
	
	// Test batch order placement
	requests := []*OrderRequest{
		{
			Symbol:          "RELIANCE-EQ",
			Quantity:        100,
			Price:           2500.0,
			OrderType:       Limit,
			TransactionType: Buy,
			Validity:        Day,
			Exchange:        "NSE",
			Product:         Normal,
		},
		{
			Symbol:          "INFY-EQ",
			Quantity:        200,
			Price:           1500.0,
			OrderType:       Limit,
			TransactionType: Buy,
			Validity:        Day,
			Exchange:        "NSE",
			Product:         Normal,
		},
	}
	
	responses, err := placementService.PlaceOrderBatchSync(ctx, requests)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error placing batch orders: %v", err)
		return
	}
	
	// Check responses
	if len(responses) != len(requests) {
		t.Errorf("Expected %d responses, got %d", len(requests), len(responses))
		return
	}
	
	for i, response := range responses {
		if !response.Status {
			t.Errorf("Batch order %d failed: %s", i, response.Error)
			return
		}
	}
	
	fmt.Println("Order placement service tests passed")
}

// TestOrderMonitoringSystem tests the order monitoring system
func TestOrderMonitoringSystem(t *testing.T) {
	// Create a mock broker adapter
	mockBroker := NewMockBrokerAdapter()
	
	// Create a smart router
	smartRouter := NewDefaultSmartRouter(BestPrice)
	
	// Create an order execution engine
	engine := NewOrderExecutionEngine(smartRouter)
	
	// Register the mock broker
	engine.RegisterBroker("MOCK", mockBroker)
	
	// Create an order monitoring system
	monitoringSystem := NewOrderMonitoringSystem(engine, 100*time.Millisecond)
	
	// Start the monitoring system
	monitoringSystem.Start()
	defer monitoringSystem.Stop()
	
	// Create and execute an order
	request := &OrderRequest{
		Symbol:          "RELIANCE-EQ",
		Quantity:        100,
		Price:           2500.0,
		OrderType:       Limit,
		TransactionType: Buy,
		Validity:        Day,
		Exchange:        "NSE",
		Product:         Normal,
	}
	
	ctx := context.Background()
	response, err := engine.ExecuteOrder(ctx, request)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error executing order: %v", err)
		return
	}
	
	// Start monitoring the order
	err = monitoringSystem.StartMonitoring(response.Order.ID, 5*time.Second)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error starting order monitoring: %v", err)
		return
	}
	
	// Register a status callback
	statusChanged := make(chan bool)
	monitoringSystem.RegisterStatusCallback(response.Order.ID, func(order *Order, previousStatus OrderStatus) {
		fmt.Printf("Order status changed from %s to %s\n", previousStatus, order.Status)
		statusChanged <- true
	})
	
	// Register an alert callback
	alertReceived := make(chan bool)
	monitoringSystem.RegisterAlertCallback(func(alert Alert) {
		fmt.Printf("Alert received: %s - %s\n", alert.Type, alert.Message)
		alertReceived <- true
	})
	
	// Set alert threshold
	monitoringSystem.SetAlertThreshold(response.Order.ID, AlertThreshold{
		DelayThreshold:      2 * time.Second,
		PriceDeviationPct:   5.0,
		PartialFillDuration: 2 * time.Second,
	})
	
	// Simulate order status change
	mockBroker.SimulateOrderStatusChange(response.Order.ID, PartiallyExecuted)
	
	// Wait for status change notification
	select {
	case <-statusChanged:
		// Status change detected
	case <-time.After(3 * time.Second):
		t.Errorf("Timeout waiting for status change notification")
		return
	}
	
	// Wait for alert
	select {
	case <-alertReceived:
		// Alert received
	case <-time.After(3 * time.Second):
		t.Errorf("Timeout waiting for alert")
		return
	}
	
	// Get monitored order
	monitoredOrder, err := monitoringSystem.GetMonitoredOrder(response.Order.ID)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error getting monitored order: %v", err)
		return
	}
	
	// Check monitored order
	if monitoredOrder.Order.Status != PartiallyExecuted {
		t.Errorf("Expected order status %s, got %s", PartiallyExecuted, monitoredOrder.Order.Status)
		return
	}
	
	// Check alerts
	alerts, err := monitoringSystem.GetAlerts(response.Order.ID)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error getting alerts: %v", err)
		return
	}
	
	// Check alert count
	if len(alerts) == 0 {
		t.Errorf("Expected at least one alert")
		return
	}
	
	// Acknowledge alert
	err = monitoringSystem.AcknowledgeAlert(response.Order.ID, alerts[0].Type)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error acknowledging alert: %v", err)
		return
	}
	
	// Stop monitoring
	monitoringSystem.StopMonitoring(response.Order.ID)
	
	fmt.Println("Order monitoring system tests passed")
}

// TestExecutionAlgorithms tests the execution algorithms
func TestExecutionAlgorithms(t *testing.T) {
	// Create a mock broker adapter
	mockBroker := NewMockBrokerAdapter()
	
	// Create a smart router
	smartRouter := NewDefaultSmartRouter(BestPrice)
	
	// Create an order execution engine
	engine := NewOrderExecutionEngine(smartRouter)
	
	// Register the mock broker
	engine.RegisterBroker("MOCK", mockBroker)
	
	// Create an algorithm manager
	algorithmManager := NewAlgorithmManager()
	
	// Register algorithms
	povAlgorithm := NewPOVAlgorithm(10.0, 5*time.Minute, func(symbol string) (int, error) {
		// Mock market volume function
		return 10000, nil
	})
	algorithmManager.RegisterAlgorithm(povAlgorithm)
	
	impactMinimizationAlgorithm := NewImpactMinimizationAlgorithm(5*time.Minute, 5, 0.2, func(symbol string) (float64, error) {
		// Mock market volatility function
		return 0.5, nil
	})
	algorithmManager.RegisterAlgorithm(impactMinimizationAlgorithm)
	
	adaptiveAlgorithm := NewAdaptiveAlgorithm(5*time.Minute, 5, 30*time.Second, func(symbol string) (float64, float64, error) {
		// Mock market data function (price, volume)
		return 2500.0, 10000.0, nil
	})
	algorithmManager.RegisterAlgorithm(adaptiveAlgorithm)
	
	// Create an order request
	request := &OrderRequest{
		Symbol:          "RELIANCE-EQ",
		Quantity:        1000,
		Price:           2500.0,
		OrderType:       Limit,
		TransactionType: Buy,
		Validity:        Day,
		Exchange:        "NSE",
		Product:         Normal,
	}
	
	// Test each algorithm
	algorithms := []string{"POV", "ImpactMinimization", "Adaptive"}
	
	for _, algorithmName := range algorithms {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		
		response, err := algorithmManager.ExecuteWithAlgorithm(ctx, engine, algorithmName, request)
		
		// Check for errors
		if err != nil {
			t.Errorf("Error executing with algorithm %s: %v", algorithmName, err)
			continue
		}
		
		// Check response
		if !response.Status {
			t.Errorf("Algorithm %s execution failed: %s", algorithmName, response.Error)
			continue
		}
		
		fmt.Printf("Algorithm %s executed successfully\n", algorithmName)
	}
	
	fmt.Println("Execution algorithms tests passed")
}

// MockBrokerAdapter is a mock implementation of the BrokerAdapter interface
type MockBrokerAdapter struct {
	orders map[string]*Order
	mutex  sync.RWMutex
}

// NewMockBrokerAdapter creates a new mock broker adapter
func NewMockBrokerAdapter() *MockBrokerAdapter {
	return &MockBrokerAdapter{
		orders: make(map[string]*Order),
	}
}

// PlaceOrder places an order with the mock broker
func (a *MockBrokerAdapter) PlaceOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Create a new order
	orderID := fmt.Sprintf("MOCK-%d", len(a.orders)+1)
	order := &Order{
		ID:              orderID,
		Symbol:          request.Symbol,
		Quantity:        request.Quantity,
		Price:           request.Price,
		OrderType:       request.OrderType,
		TransactionType: request.TransactionType,
		Status:          Open,
		FilledQuantity:  0,
		AveragePrice:    0,
		PlacedAt:        time.Now(),
		UpdatedAt:       time.Now(),
		Validity:        request.Validity,
		TriggerPrice:    request.TriggerPrice,
		Exchange:        request.Exchange,
		Product:         request.Product,
		BrokerOrderID:   orderID,
		StrategyID:      request.StrategyID,
		Tags:            request.Tags,
	}
	
	// Store the order
	a.orders[orderID] = order
	
	return &OrderResponse{
		Order:  order,
		Status: true,
	}, nil
}

// ModifyOrder modifies an existing order with the mock broker
func (a *MockBrokerAdapter) ModifyOrder(ctx context.Context, orderID string, request *OrderRequest) (*OrderResponse, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Find the order
	order, exists := a.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	
	// Update the order
	order.Quantity = request.Quantity
	order.Price = request.Price
	order.OrderType = request.OrderType
	order.Validity = request.Validity
	order.TriggerPrice = request.TriggerPrice
	order.UpdatedAt = time.Now()
	
	return &OrderResponse{
		Order:  order,
		Status: true,
	}, nil
}

// CancelOrder cancels an existing order with the mock broker
func (a *MockBrokerAdapter) CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Find the order
	order, exists := a.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	
	// Cancel the order
	order.Status = Cancelled
	order.UpdatedAt = time.Now()
	
	return &OrderResponse{
		Order:  order,
		Status: true,
	}, nil
}

// GetOrderStatus gets the status of an order from the mock broker
func (a *MockBrokerAdapter) GetOrderStatus(ctx context.Context, orderID string) (*Order, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	// Find the order
	order, exists := a.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	
	return order, nil
}

// GetOrders gets all orders from the mock broker
func (a *MockBrokerAdapter) GetOrders(ctx context.Context) ([]*Order, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	// Create a slice of orders
	orders := make([]*Order, 0, len(a.orders))
	for _, order := range a.orders {
		orders = append(orders, order)
	}
	
	return orders, nil
}

// SimulateOrderStatusChange simulates a change in order status
func (a *MockBrokerAdapter) SimulateOrderStatusChange(orderID string, newStatus OrderStatus) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Find the order
	order, exists := a.orders[orderID]
	if !exists {
		log.Printf("Order not found: %s", orderID)
		return
	}
	
	// Update the order status
	order.Status = newStatus
	order.UpdatedAt = time.Now()
	
	// Update filled quantity for partially executed orders
	if newStatus == PartiallyExecuted {
		order.FilledQuantity = order.Quantity / 2
		order.AveragePrice = order.Price
	} else if newStatus == Executed {
		order.FilledQuantity = order.Quantity
		order.AveragePrice = order.Price
	}
}

// TestBrokerIntegration tests the broker integration
func TestBrokerIntegration(t *testing.T) {
	// Create a broker factory
	brokerFactory := NewBrokerFactory()
	
	// List available brokers
	brokers := brokerFactory.ListBrokers()
	fmt.Printf("Available brokers: %v\n", brokers)
	
	// Create a mock broker
	config := map[string]string{
		"baseURL":       "https://api.xts.com",
		"apiKey":        "mock-api-key",
		"apiSecret":     "mock-api-secret",
		"clientCode":    "mock-client-code",
		"isInteractive": "true",
	}
	
	broker, err := brokerFactory.CreateBroker("XTS", config)
	
	// Check for errors
	if err != nil {
		t.Errorf("Error creating broker: %v", err)
		return
	}
	
	// Check broker
	if broker == nil {
		t.Errorf("Broker is nil")
		return
	}
	
	fmt.Println("Broker integration tests passed")
}

// RunAllTests runs all tests
func RunAllTests() {
	// Create a testing.T instance
	t := &testing.T{}
	
	// Run tests
	fmt.Println("Running order execution engine tests...")
	TestOrderExecutionEngine(t)
	
	fmt.Println("\nRunning smart router tests...")
	TestSmartRouter(t)
	
	fmt.Println("\nRunning order placement service tests...")
	TestOrderPlacementService(t)
	
	fmt.Println("\nRunning order monitoring system tests...")
	TestOrderMonitoringSystem(t)
	
	fmt.Println("\nRunning execution algorithms tests...")
	TestExecutionAlgorithms(t)
	
	fmt.Println("\nRunning broker integration tests...")
	TestBrokerIntegration(t)
	
	fmt.Println("\nAll tests completed")
}
