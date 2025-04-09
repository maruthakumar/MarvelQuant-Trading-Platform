package orderexecution

import (
	"context"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the Logger interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

// MockErrorHandler is a mock implementation of the ErrorHandler interface
type MockErrorHandler struct {
	mock.Mock
}

func (m *MockErrorHandler) HandleError(ctx context.Context, err error) (bool, error) {
	args := m.Called(ctx, err)
	return args.Bool(0), args.Error(1)
}

// MockBrokerConnector is a mock implementation of the BrokerConnector interface
type MockBrokerConnector struct {
	mock.Mock
}

func (m *MockBrokerConnector) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBrokerConnector) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBrokerConnector) PlaceOrder(ctx context.Context, order Order) (string, error) {
	args := m.Called(ctx, order)
	return args.String(0), args.Error(1)
}

func (m *MockBrokerConnector) CancelOrder(ctx context.Context, brokerOrderID string) error {
	args := m.Called(ctx, brokerOrderID)
	return args.Error(0)
}

func (m *MockBrokerConnector) ModifyOrder(ctx context.Context, brokerOrderID string, price float64, quantity int, triggerPrice float64) error {
	args := m.Called(ctx, brokerOrderID, price, quantity, triggerPrice)
	return args.Error(0)
}

func (m *MockBrokerConnector) GetOrderStatus(ctx context.Context, brokerOrderID string) (OrderStatus, error) {
	args := m.Called(ctx, brokerOrderID)
	return args.Get(0).(OrderStatus), args.Error(1)
}

// MockOrderRouter is a mock implementation of the OrderRouter interface
type MockOrderRouter struct {
	mock.Mock
}

func (m *MockOrderRouter) RouteOrder(order Order) (BrokerConnector, error) {
	args := m.Called(order)
	return args.Get(0).(BrokerConnector), args.Error(1)
}

// MockRiskManager is a mock implementation of the RiskManager interface
type MockRiskManager struct {
	mock.Mock
}

func (m *MockRiskManager) ValidateOrder(order Order, portfolio Portfolio, strategy Strategy) error {
	args := m.Called(order, portfolio, strategy)
	return args.Error(0)
}

func (m *MockRiskManager) CheckPositionLimits(order Order, portfolio Portfolio) error {
	args := m.Called(order, portfolio)
	return args.Error(0)
}

func (m *MockRiskManager) CheckMarginRequirements(order Order, portfolio Portfolio) error {
	args := m.Called(order, portfolio)
	return args.Error(0)
}

func (m *MockRiskManager) CheckRiskParameters(order Order, strategy Strategy) error {
	args := m.Called(order, strategy)
	return args.Error(0)
}

func (m *MockRiskManager) CheckRateLimits(order Order) error {
	args := m.Called(order)
	return args.Error(0)
}

// TestErrorHandling tests the error handling functionality
func TestErrorHandling(t *testing.T) {
	// Create mock logger
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Create error handler
	errorHandler := NewDefaultErrorHandler(3, 100*time.Millisecond, mockLogger)

	// Create execution error
	executionError := NewExecutionError(
		ErrorTypeExecution,
		ErrorSeverityError,
		ErrCodeExecutionFailed,
		"Test execution error",
		nil,
		"TestErrorHandling",
	)

	// Test HandleError
	ctx := context.Background()
	shouldRetry, err := errorHandler.HandleError(ctx, executionError)
	assert.True(t, shouldRetry)
	assert.Error(t, err)
	assert.Equal(t, executionError, err)

	// Test retry exhaustion
	for i := 0; i < 3; i++ {
		shouldRetry, err = errorHandler.HandleError(ctx, executionError)
	}
	assert.False(t, shouldRetry)
	assert.Error(t, err)

	// Test circuit breaker
	cb := NewCircuitBreaker("test", 3, 100*time.Millisecond, 1)
	
	// Test successful execution
	err = cb.Execute(func() error {
		return nil
	})
	assert.NoError(t, err)

	// Test failed execution
	testErr := NewExecutionError(
		ErrorTypeExecution,
		ErrorSeverityError,
		ErrCodeExecutionFailed,
		"Test circuit breaker error",
		nil,
		"TestErrorHandling",
	)
	
	// Trigger circuit breaker
	for i := 0; i < 4; i++ {
		err = cb.Execute(func() error {
			return testErr
		})
	}
	
	// Circuit should be open now
	assert.Error(t, err)
	assert.Equal(t, CircuitStateOpen, cb.State())

	// Verify mock expectations
	mockLogger.AssertExpectations(t)
}

// TestRiskManagement tests the risk management functionality
func TestRiskManagement(t *testing.T) {
	// Create mock logger
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Create mock error handler
	mockErrorHandler := new(MockErrorHandler)
	mockErrorHandler.On("HandleError", mock.Anything, mock.Anything).Return(false, nil)

	// Create default risk manager
	defaultRiskManager := NewDefaultRiskManager(1000, 10000, 100)

	// Create enhanced risk manager
	riskManager := NewEnhancedRiskManager(defaultRiskManager, mockLogger, mockErrorHandler)

	// Create risk profile
	profile := RiskProfile{
		ID:          "test-profile",
		Name:        "Test Profile",
		Description: "Test risk profile",
		Limits: map[RiskLimitType]RiskLimit{
			RiskLimitTypeOrderValue: {
				Type:        RiskLimitTypeOrderValue,
				Value:       10000,
				Level:       RiskLevelMedium,
				Description: "Maximum order value",
				Enabled:     true,
			},
			RiskLimitTypePositionSize: {
				Type:        RiskLimitTypePositionSize,
				Value:       100,
				Level:       RiskLevelMedium,
				Description: "Maximum position size",
				Enabled:     true,
			},
		},
	}

	// Test CreateRiskProfile
	err := riskManager.CreateRiskProfile(profile)
	assert.NoError(t, err)

	// Test GetRiskProfile
	retrievedProfile, err := riskManager.GetRiskProfile("test-profile")
	assert.NoError(t, err)
	assert.Equal(t, profile.ID, retrievedProfile.ID)
	assert.Equal(t, profile.Name, retrievedProfile.Name)

	// Test UpdateRiskProfile
	profile.Name = "Updated Test Profile"
	err = riskManager.UpdateRiskProfile(profile)
	assert.NoError(t, err)

	// Test GetRiskProfile after update
	retrievedProfile, err = riskManager.GetRiskProfile("test-profile")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test Profile", retrievedProfile.Name)

	// Test ValidateOrder with valid order
	order := Order{
		ID:          "test-order",
		PortfolioID: "test-portfolio",
		StrategyID:  "test-strategy",
		Symbol:      "AAPL",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    10,
		Price:       150.0,
	}

	portfolio := Portfolio{
		ID: "test-portfolio",
	}

	strategy := Strategy{
		ID: "test-strategy",
		RiskParams: map[string]interface{}{
			"riskProfileID": "test-profile",
		},
	}

	err = riskManager.ValidateOrder(order, portfolio, strategy)
	assert.NoError(t, err)

	// Test ValidateOrder with invalid order (exceeds order value limit)
	invalidOrder := Order{
		ID:          "invalid-order",
		PortfolioID: "test-portfolio",
		StrategyID:  "test-strategy",
		Symbol:      "AAPL",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    1000,
		Price:       150.0, // Total value: 150,000 > limit 10,000
	}

	err = riskManager.ValidateOrder(invalidOrder, portfolio, strategy)
	assert.Error(t, err)
	executionError, ok := err.(*ExecutionError)
	assert.True(t, ok)
	assert.Equal(t, ErrCodeInvalidOrder, executionError.Code)

	// Verify mock expectations
	mockLogger.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

// TestOrderLifecycle tests the order lifecycle management functionality
func TestOrderLifecycle(t *testing.T) {
	// Create mock logger
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Create mock error handler
	mockErrorHandler := new(MockErrorHandler)
	mockErrorHandler.On("HandleError", mock.Anything, mock.Anything).Return(false, nil)

	// Create lifecycle manager
	lifecycleManager := NewOrderLifecycleManager(mockLogger, mockErrorHandler)

	// Create test order
	order := Order{
		ID:          "test-order",
		PortfolioID: "test-portfolio",
		StrategyID:  "test-strategy",
		Symbol:      "AAPL",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    10,
		Price:       150.0,
		Status:      OrderStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test CreateLifecycle
	lifecycle, err := lifecycleManager.CreateLifecycle(order)
	assert.NoError(t, err)
	assert.Equal(t, OrderLifecycleStateCreated, lifecycle.CurrentState)
	assert.Equal(t, 1, len(lifecycle.Events))

	// Test GetLifecycle
	retrievedLifecycle, err := lifecycleManager.GetLifecycle("test-order")
	assert.NoError(t, err)
	assert.Equal(t, lifecycle.CurrentState, retrievedLifecycle.CurrentState)

	// Test TransitionState
	err = lifecycleManager.TransitionState(
		"test-order",
		OrderLifecycleStateValidated,
		"VALIDATION_PASSED",
		nil,
	)
	assert.NoError(t, err)

	// Test GetLifecycle after transition
	retrievedLifecycle, err = lifecycleManager.GetLifecycle("test-order")
	assert.NoError(t, err)
	assert.Equal(t, OrderLifecycleStateValidated, retrievedLifecycle.CurrentState)
	assert.Equal(t, 2, len(retrievedLifecycle.Events))

	// Test invalid transition
	err = lifecycleManager.TransitionState(
		"test-order",
		OrderLifecycleStateCompleted,
		"INVALID_TRANSITION",
		nil,
	)
	assert.Error(t, err)
	executionError, ok := err.(*ExecutionError)
	assert.True(t, ok)
	assert.Equal(t, ErrCodeInvalidOrder, executionError.Code)

	// Test valid transitions through the lifecycle
	validTransitions := []struct {
		state     OrderLifecycleState
		eventType string
	}{
		{OrderLifecycleStateSubmitted, "ORDER_SUBMITTED"},
		{OrderLifecycleStateAcknowledged, "ORDER_ACKNOWLEDGED"},
		{OrderLifecycleStatePartiallyFilled, "ORDER_PARTIALLY_FILLED"},
		{OrderLifecycleStateCompleted, "ORDER_COMPLETED"},
	}

	for _, transition := range validTransitions {
		err = lifecycleManager.TransitionState(
			"test-order",
			transition.state,
			transition.eventType,
			nil,
		)
		assert.NoError(t, err)

		retrievedLifecycle, err = lifecycleManager.GetLifecycle("test-order")
		assert.NoError(t, err)
		assert.Equal(t, transition.state, retrievedLifecycle.CurrentState)
	}

	// Test GetOrderEvents
	events, err := lifecycleManager.GetOrderEvents("test-order")
	assert.NoError(t, err)
	assert.Equal(t, 6, len(events)) // Initial + 5 transitions

	// Verify mock expectations
	mockLogger.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

// TestOrderDependencies tests the order dependency management functionality
func TestOrderDependencies(t *testing.T) {
	// Create mock logger
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Create mock error handler
	mockErrorHandler := new(MockErrorHandler)
	mockErrorHandler.On("HandleError", mock.Anything, mock.Anything).Return(false, nil)

	// Create lifecycle manager
	lifecycleManager := NewOrderLifecycleManager(mockLogger, mockErrorHandler)

	// Create dependency manager
	dependencyManager := NewOrderDependencyManager(lifecycleManager, mockLogger)

	// Create parent order
	parentOrder := Order{
		ID:          "parent-order",
		PortfolioID: "test-portfolio",
		StrategyID:  "test-strategy",
		Symbol:      "AAPL",
		Exchange:    "NSE",
		OrderType:   OrderTypeLimit,
		ProductType: ProductTypeDelivery,
		Side:        OrderSideBuy,
		Quantity:    10,
		Price:       150.0,
		Status:      OrderStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create child order
	childOrder := Order{
		ID:            "child-order",
		PortfolioID:   "test-portfolio",
		StrategyID:    "test-strategy",
		Symbol:        "AAPL",
		Exchange:      "NSE",
		OrderType:     OrderTypeLimit,
		ProductType:   ProductTypeDelivery,
		Side:          OrderSideSell,
		Quantity:      10,
		Price:         160.0,
		Status:        OrderStatusPending,
		ParentOrderID: "parent-order",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create lifecycles
	_, err := lifecycleManager.CreateLifecycle(parentOrder)
	assert.NoError(t, err)

	_, err = lifecycleManager.CreateLifecycle(childOrder)
	assert.NoError(t, err)

	// Test CreateDependency
	dependency, err := dependencyManager.CreateDependency(
		"parent-order",
		"child-order",
		"OTO", // One-Triggers-Other
		"",
	)
	assert.NoError(t, err)
	assert.Equal(t, "parent-order", dependency.ParentOrderID)
	assert.Equal(t, "child-order", dependency.ChildOrderID)
	assert.Equal(t, "OTO", dependency.Type)

	// Test GetDependencies
	dependencies, err := dependencyManager.GetDependencies("parent-order")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dependencies))
	assert.Equal(t, "child-order", dependencies[0].ChildOrderID)

	// Test GetParentOrder
	parentID, err := dependencyManager.GetParentOrder("child-order")
	assert.NoError(t, err)
	assert.Equal(t, "parent-order", parentID)

	// Test parent order completion triggering child order
	// Register a callback to verify the child order state change
	childOrderSubmitted := false
	lifecycleManager.RegisterCallback(OrderLifecycleStateSubmitted, func(lifecycle *OrderLifecycle, event OrderEvent) {
		if lifecycle.Order.ID == "child-order" {
			childOrderSubmitted = true
		}
	})

	// Transition parent order to completed
	err = lifecycleManager.TransitionState(
		"parent-order",
		OrderLifecycleStateValidated,
		"VALIDATION_PASSED",
		nil,
	)
	assert.NoError(t, err)

	err = lifecycleManager.TransitionState(
		"parent-order",
		OrderLifecycleStateSubmitted,
		"ORDER_SUBMITTED",
		nil,
	)
	assert.NoError(t, err)

	err = lifecycleManager.TransitionState(
		"parent-order",
		OrderLifecycleStateAcknowledged,
		"ORDER_ACKNOWLEDGED",
		nil,
	)
	assert.NoError(t, err)

	err = lifecycleManager.TransitionState(
		"parent-order",
		OrderLifecycleStateCompleted,
		"ORDER_COMPLETED",
		nil,
	)
	assert.NoError(t, err)

	// Give time for the callback to execute
	time.Sleep(100 * time.Millisecond)

	// Check if child order was submitted
	childLifecycle, err := lifecycleManager.GetLifecycle("child-order")
	assert.NoError(t, err)
	assert.Equal(t, OrderLifecycleStateSubmitted, childLifecycle.CurrentState)
	assert.True(t, childOrderSubmitted)

	// Test DeleteDependency
	err = dependencyManager.DeleteDependency(dependency.ID)
	assert.NoError(t, err)

	// Test GetDependencies after deletion
	dependencies, err = dependencyManager.GetDependencies("parent-order")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(dependencies))

	// Verify mock expectations
	mockLogger.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

// TestEnhancedOrderExecutionEngine tests the enhanced order execution engine
func TestEnhancedOrderExecutionEngine(t *testing.T) {
	// Create mock components
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	mockErrorHandler := new(MockErrorHandler)
	mockErrorHandler.On("HandleError", mock.Anything, mock.Anything).R
(Content truncated due to size limit. Use line ranges to read in chunks)