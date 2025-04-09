package orderexecution

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// MockBrokerType represents the type of mock broker
type MockBrokerType string

const (
	// Mock broker types
	MockBrokerTypeXTS     MockBrokerType = "XTS"
	MockBrokerTypeZerodha MockBrokerType = "ZERODHA"
	MockBrokerTypeGeneric MockBrokerType = "GENERIC"
)

// MockBrokerConfig represents the configuration for a mock broker
type MockBrokerConfig struct {
	Type                 MockBrokerType `json:"type"`
	Name                 string         `json:"name"`
	SimulateLatency      bool           `json:"simulateLatency"`
	BaseLatencyMs        int            `json:"baseLatencyMs"`
	LatencyVariationMs   int            `json:"latencyVariationMs"`
	SimulateErrors       bool           `json:"simulateErrors"`
	ErrorRate            float64        `json:"errorRate"`
	SimulatePartialFills bool           `json:"simulatePartialFills"`
	PartialFillRate      float64        `json:"partialFillRate"`
	SimulateRejections   bool           `json:"simulateRejections"`
	RejectionRate        float64        `json:"rejectionRate"`
}

// MockBrokerConnector implements a mock broker connector for testing
type MockBrokerConnector struct {
	config         MockBrokerConfig
	orders         map[string]Order // brokerOrderID -> order
	orderStatus    map[string]OrderStatus
	fills          map[string][]OrderFill
	connected      bool
	logger         Logger
	errorHandler   ErrorHandler
	mutex          sync.RWMutex
	statusUpdates  chan OrderStatusUpdate
	fillUpdates    chan OrderFill
	simulationDone chan bool
}

// OrderStatusUpdate represents an update to an order's status
type OrderStatusUpdate struct {
	BrokerOrderID string      `json:"brokerOrderId"`
	Status        OrderStatus `json:"status"`
	Message       string      `json:"message"`
	Timestamp     time.Time   `json:"timestamp"`
}

// OrderFill represents a fill for an order
type OrderFill struct {
	BrokerOrderID string    `json:"brokerOrderId"`
	FillID        string    `json:"fillId"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	Timestamp     time.Time `json:"timestamp"`
}

// NewMockBrokerConnector creates a new mock broker connector
func NewMockBrokerConnector(config MockBrokerConfig, logger Logger, errorHandler ErrorHandler) *MockBrokerConnector {
	return &MockBrokerConnector{
		config:         config,
		orders:         make(map[string]Order),
		orderStatus:    make(map[string]OrderStatus),
		fills:          make(map[string][]OrderFill),
		connected:      false,
		logger:         logger,
		errorHandler:   errorHandler,
		statusUpdates:  make(chan OrderStatusUpdate, 100),
		fillUpdates:    make(chan OrderFill, 100),
		simulationDone: make(chan bool),
	}
}

// Connect connects to the mock broker
func (m *MockBrokerConnector) Connect(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Simulate connection latency
	if m.config.SimulateLatency {
		time.Sleep(time.Duration(m.config.BaseLatencyMs) * time.Millisecond)
	}

	// Simulate connection error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return NewExecutionError(
				ErrorTypeConnection,
				ErrorSeverityError,
				ErrCodeConnectionFailed,
				fmt.Sprintf("Failed to connect to %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			)
		}
	}

	m.connected = true
	m.logger.Info("Connected to mock broker",
		"broker", m.config.Name,
		"type", string(m.config.Type),
	)

	return nil
}

// Disconnect disconnects from the mock broker
func (m *MockBrokerConnector) Disconnect(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Simulate disconnection latency
	if m.config.SimulateLatency {
		time.Sleep(time.Duration(m.config.BaseLatencyMs) * time.Millisecond)
	}

	// Simulate disconnection error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return NewExecutionError(
				ErrorTypeConnection,
				ErrorSeverityError,
				ErrCodeDisconnectionFailed,
				fmt.Sprintf("Failed to disconnect from %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			)
		}
	}

	m.connected = false
	close(m.statusUpdates)
	close(m.fillUpdates)
	close(m.simulationDone)

	m.logger.Info("Disconnected from mock broker",
		"broker", m.config.Name,
		"type", string(m.config.Type),
	)

	return nil
}

// PlaceOrder places an order with the mock broker
func (m *MockBrokerConnector) PlaceOrder(ctx context.Context, order Order) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if connected
	if !m.connected {
		return "", NewExecutionError(
			ErrorTypeConnection,
			ErrorSeverityError,
			ErrCodeNotConnected,
			fmt.Sprintf("Not connected to %s broker", m.config.Name),
			nil,
			"MockBrokerConnector",
		)
	}

	// Simulate order placement latency
	if m.config.SimulateLatency {
		latency := m.config.BaseLatencyMs + randInt(-m.config.LatencyVariationMs, m.config.LatencyVariationMs)
		if latency > 0 {
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}

	// Simulate order placement error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return "", NewExecutionError(
				ErrorTypeExecution,
				ErrorSeverityError,
				ErrCodeOrderPlacementFailed,
				fmt.Sprintf("Failed to place order with %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			).WithOrderID(order.ID)
		}
	}

	// Simulate order rejection
	if m.config.SimulateRejections && m.config.RejectionRate > 0 {
		if randFloat() < m.config.RejectionRate {
			return "", NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				ErrCodeOrderRejected,
				fmt.Sprintf("Order rejected by %s broker: Invalid parameters", m.config.Name),
				nil,
				"MockBrokerConnector",
			).WithOrderID(order.ID)
		}
	}

	// Generate broker order ID
	brokerOrderID := fmt.Sprintf("%s_%s_%d", strings.ToLower(string(m.config.Type)), order.ID, time.Now().UnixNano())

	// Store order
	orderCopy := order
	orderCopy.BrokerOrderID = brokerOrderID
	orderCopy.Status = OrderStatusSubmitted
	orderCopy.UpdatedAt = time.Now()
	m.orders[brokerOrderID] = orderCopy
	m.orderStatus[brokerOrderID] = OrderStatusSubmitted

	// Start order simulation
	go m.simulateOrderExecution(brokerOrderID, orderCopy)

	m.logger.Info("Placed order with mock broker",
		"broker", m.config.Name,
		"type", string(m.config.Type),
		"orderId", order.ID,
		"brokerOrderId", brokerOrderID,
		"symbol", order.Symbol,
		"side", string(order.Side),
		"quantity", order.Quantity,
		"price", order.Price,
	)

	return brokerOrderID, nil
}

// CancelOrder cancels an order with the mock broker
func (m *MockBrokerConnector) CancelOrder(ctx context.Context, brokerOrderID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if connected
	if !m.connected {
		return NewExecutionError(
			ErrorTypeConnection,
			ErrorSeverityError,
			ErrCodeNotConnected,
			fmt.Sprintf("Not connected to %s broker", m.config.Name),
			nil,
			"MockBrokerConnector",
		)
	}

	// Simulate order cancellation latency
	if m.config.SimulateLatency {
		latency := m.config.BaseLatencyMs + randInt(-m.config.LatencyVariationMs, m.config.LatencyVariationMs)
		if latency > 0 {
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}

	// Simulate order cancellation error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return NewExecutionError(
				ErrorTypeExecution,
				ErrorSeverityError,
				ErrCodeOrderCancellationFailed,
				fmt.Sprintf("Failed to cancel order with %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			)
		}
	}

	// Check if order exists
	order, exists := m.orders[brokerOrderID]
	if !exists {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Order not found: %s", brokerOrderID),
			nil,
			"MockBrokerConnector",
		)
	}

	// Check if order can be cancelled
	status := m.orderStatus[brokerOrderID]
	if status == OrderStatusCompleted || status == OrderStatusCancelled || status == OrderStatusRejected {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidOrder,
			fmt.Sprintf("Order cannot be cancelled in state %s", status),
			nil,
			"MockBrokerConnector",
		)
	}

	// Update order status
	m.orderStatus[brokerOrderID] = OrderStatusCancelled
	order.Status = OrderStatusCancelled
	order.UpdatedAt = time.Now()
	m.orders[brokerOrderID] = order

	// Send status update
	m.statusUpdates <- OrderStatusUpdate{
		BrokerOrderID: brokerOrderID,
		Status:        OrderStatusCancelled,
		Message:       "Order cancelled",
		Timestamp:     time.Now(),
	}

	m.logger.Info("Cancelled order with mock broker",
		"broker", m.config.Name,
		"type", string(m.config.Type),
		"orderId", order.ID,
		"brokerOrderId", brokerOrderID,
	)

	return nil
}

// ModifyOrder modifies an order with the mock broker
func (m *MockBrokerConnector) ModifyOrder(ctx context.Context, brokerOrderID string, price float64, quantity int, triggerPrice float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if connected
	if !m.connected {
		return NewExecutionError(
			ErrorTypeConnection,
			ErrorSeverityError,
			ErrCodeNotConnected,
			fmt.Sprintf("Not connected to %s broker", m.config.Name),
			nil,
			"MockBrokerConnector",
		)
	}

	// Simulate order modification latency
	if m.config.SimulateLatency {
		latency := m.config.BaseLatencyMs + randInt(-m.config.LatencyVariationMs, m.config.LatencyVariationMs)
		if latency > 0 {
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}

	// Simulate order modification error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return NewExecutionError(
				ErrorTypeExecution,
				ErrorSeverityError,
				ErrCodeOrderModificationFailed,
				fmt.Sprintf("Failed to modify order with %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			)
		}
	}

	// Check if order exists
	order, exists := m.orders[brokerOrderID]
	if !exists {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Order not found: %s", brokerOrderID),
			nil,
			"MockBrokerConnector",
		)
	}

	// Check if order can be modified
	status := m.orderStatus[brokerOrderID]
	if status == OrderStatusCompleted || status == OrderStatusCancelled || status == OrderStatusRejected {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidOrder,
			fmt.Sprintf("Order cannot be modified in state %s", status),
			nil,
			"MockBrokerConnector",
		)
	}

	// Update order
	order.Price = price
	order.Quantity = quantity
	order.TriggerPrice = triggerPrice
	order.UpdatedAt = time.Now()
	m.orders[brokerOrderID] = order

	m.logger.Info("Modified order with mock broker",
		"broker", m.config.Name,
		"type", string(m.config.Type),
		"orderId", order.ID,
		"brokerOrderId", brokerOrderID,
		"price", price,
		"quantity", quantity,
		"triggerPrice", triggerPrice,
	)

	return nil
}

// GetOrderStatus gets the status of an order from the mock broker
func (m *MockBrokerConnector) GetOrderStatus(ctx context.Context, brokerOrderID string) (OrderStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if connected
	if !m.connected {
		return "", NewExecutionError(
			ErrorTypeConnection,
			ErrorSeverityError,
			ErrCodeNotConnected,
			fmt.Sprintf("Not connected to %s broker", m.config.Name),
			nil,
			"MockBrokerConnector",
		)
	}

	// Simulate order status latency
	if m.config.SimulateLatency {
		latency := m.config.BaseLatencyMs + randInt(-m.config.LatencyVariationMs, m.config.LatencyVariationMs)
		if latency > 0 {
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}

	// Simulate order status error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return "", NewExecutionError(
				ErrorTypeExecution,
				ErrorSeverityError,
				ErrCodeOrderStatusFailed,
				fmt.Sprintf("Failed to get order status from %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			)
		}
	}

	// Check if order exists
	status, exists := m.orderStatus[brokerOrderID]
	if !exists {
		return "", NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Order not found: %s", brokerOrderID),
			nil,
			"MockBrokerConnector",
		)
	}

	return status, nil
}

// GetOrderFills gets the fills for an order from the mock broker
func (m *MockBrokerConnector) GetOrderFills(ctx context.Context, brokerOrderID string) ([]OrderFill, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if connected
	if !m.connected {
		return nil, NewExecutionError(
			ErrorTypeConnection,
			ErrorSeverityError,
			ErrCodeNotConnected,
			fmt.Sprintf("Not connected to %s broker", m.config.Name),
			nil,
			"MockBrokerConnector",
		)
	}

	// Simulate order fills latency
	if m.config.SimulateLatency {
		latency := m.config.BaseLatencyMs + randInt(-m.config.LatencyVariationMs, m.config.LatencyVariationMs)
		if latency > 0 {
			time.Sleep(time.Duration(latency) * time.Millisecond)
		}
	}

	// Simulate order fills error
	if m.config.SimulateErrors && m.config.ErrorRate > 0 {
		if randFloat() < m.config.ErrorRate {
			return nil, NewExecutionError(
				ErrorTypeExecution,
				ErrorSeverityError,
				ErrCodeOrderFillsFailed,
				fmt.Sprintf("Failed to get order fills from %s broker", m.config.Name),
				nil,
				"MockBrokerConnector",
			)
		}
	}

	// Check if order exists
	fills, exists := m.fills[brokerOrderID]
	if !exists {
		return []OrderFill{}, nil
	}

	return fills, nil
}

// GetStatusUpdateChannel gets the channel for order status updates
func (m *MockBrokerConnector) GetStatusUpdateChannel() <-chan OrderStatusUpdate {
	return m.statusUpdates
}

// GetFillUpdateChannel gets the channel for order fill updates
func (m *MockBrokerConnector) GetFillUpdateChannel() <-chan OrderFill {
	return m.fillUpdates
}

// simulateOrderExecution simulates the execution of an order
func (m *MockBrokerConnector) simulateOrderExecution(brokerOrderID string, order Order) {
	// Simulate order acknowledgement
	time.Sleep(time.Duration(100+randInt(0, 200)) * time.Millisecond)

	m.mutex.Lock()
	m.orderStatus[brokerOrderID] = OrderStatusOpen
	m.statusUpdates <- OrderStatusUpdate{
		BrokerOrderID: brokerOrderID,
		Status:        OrderStatusOpen,
		Message:       "Order acknowledged",
		Timestamp:     time.Now(),
	}
	m.mutex.Unlock()

	// Simulate partial fills if enabled
	if m.config.SimulatePartialFills && m.config.PartialFillRate > 0 {
		// Determine number of partial fills
		maxPartialFills := order.Quantity / 2
		if maxPartialFills > 5 {
			maxPartialFills = 5
		}
		if maxPartialFills > 0 {
			numPartialFills := randInt(1, maxPartialFills)
			remainingQuantity := order.Quantity

			for i := 0; i < numPartialFills && remainingQuantity > 0; i++ {
				// Simulate fill latency
				time.Sleep(time.Duration(500+randInt(0, 1000)) * time.Millisecond)

				// Calculate fill quantity
				fillQuantity := randInt(1, remainingQuantity/2)
				if i == numPartialFills-1 || fillQuantity == 0 {
					fillQuantity = remainingQuantity
				}
				remainingQuantity -= fillQuantity

				// Calculate fill price with some slippage
				slippageBps := randInt(-10, 10)
				fillPrice := order.Price * (1.0 + float64(slippageBps)/10000.0)

				// Create fill
				fill := OrderFill{
					BrokerOrderID: brokerOrderID,
					FillID:        fmt.Sprintf("fill_%s_%d", brokerOrderID, i+1),
					Quantity:      fillQuantity,
					Price:         fillPrice,
					Timestamp:     time.Now(),
				}

				m.mutex.Lock()
				// Add fill to order fills
				if _, exists := m.fills[brokerOrderID]; !exists {
					m.fills[brokerOrderID] = make([]OrderFill, 0)
				}
				m.fills[brokerOrderID] = append(m.fills[brokerOrderID
(Content truncated due to size limit. Use line ranges to read in chunks)