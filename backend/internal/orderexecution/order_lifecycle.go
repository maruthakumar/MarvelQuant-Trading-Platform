package orderexecution

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// OrderLifecycleState represents the state of an order in its lifecycle
type OrderLifecycleState string

const (
	// Order lifecycle states
	OrderLifecycleStateCreated     OrderLifecycleState = "CREATED"
	OrderLifecycleStateValidated   OrderLifecycleState = "VALIDATED"
	OrderLifecycleStateSubmitted   OrderLifecycleState = "SUBMITTED"
	OrderLifecycleStateAcknowledged OrderLifecycleState = "ACKNOWLEDGED"
	OrderLifecycleStatePartiallyFilled OrderLifecycleState = "PARTIALLY_FILLED"
	OrderLifecycleStateCompleted   OrderLifecycleState = "COMPLETED"
	OrderLifecycleStateCancelling  OrderLifecycleState = "CANCELLING"
	OrderLifecycleStateCancelled   OrderLifecycleState = "CANCELLED"
	OrderLifecycleStateRejected    OrderLifecycleState = "REJECTED"
	OrderLifecycleStateExpired     OrderLifecycleState = "EXPIRED"
	OrderLifecycleStateFailed      OrderLifecycleState = "FAILED"
)

// OrderEvent represents an event in the lifecycle of an order
type OrderEvent struct {
	ID          string             `json:"id"`
	OrderID     string             `json:"orderId"`
	Type        string             `json:"type"`
	PreviousState OrderLifecycleState `json:"previousState"`
	CurrentState OrderLifecycleState `json:"currentState"`
	Timestamp   time.Time          `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// OrderLifecycle manages the lifecycle of an order
type OrderLifecycle struct {
	Order       Order              `json:"order"`
	CurrentState OrderLifecycleState `json:"currentState"`
	Events      []OrderEvent       `json:"events"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	ExpiresAt   *time.Time         `json:"expiresAt,omitempty"`
}

// OrderLifecycleManager manages order lifecycles
type OrderLifecycleManager struct {
	lifecycles map[string]*OrderLifecycle
	callbacks  map[OrderLifecycleState][]OrderLifecycleCallback
	mutex      sync.RWMutex
	logger     Logger
	errorHandler ErrorHandler
}

// OrderLifecycleCallback is a function that is called when an order changes state
type OrderLifecycleCallback func(lifecycle *OrderLifecycle, event OrderEvent)

// NewOrderLifecycleManager creates a new order lifecycle manager
func NewOrderLifecycleManager(logger Logger, errorHandler ErrorHandler) *OrderLifecycleManager {
	return &OrderLifecycleManager{
		lifecycles: make(map[string]*OrderLifecycle),
		callbacks:  make(map[OrderLifecycleState][]OrderLifecycleCallback),
		logger:     logger,
		errorHandler: errorHandler,
	}
}

// CreateLifecycle creates a new order lifecycle
func (m *OrderLifecycleManager) CreateLifecycle(order Order) (*OrderLifecycle, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if lifecycle already exists
	if _, exists := m.lifecycles[order.ID]; exists {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidOrder,
			fmt.Sprintf("Order lifecycle already exists for order ID %s", order.ID),
			nil,
			"OrderLifecycleManager",
		).WithOrderID(order.ID)
	}

	// Create new lifecycle
	now := time.Now()
	lifecycle := &OrderLifecycle{
		Order:       order,
		CurrentState: OrderLifecycleStateCreated,
		Events:      make([]OrderEvent, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Set expiration time if order has validity
	if order.Validity == "GTD" && order.Metadata.AdditionalParams != nil {
		if expiryStr, ok := order.Metadata.AdditionalParams["expiryTime"].(string); ok {
			expiry, err := time.Parse(time.RFC3339, expiryStr)
			if err == nil {
				lifecycle.ExpiresAt = &expiry
			}
		}
	}

	// Add initial event
	event := OrderEvent{
		ID:          generateEventID(),
		OrderID:     order.ID,
		Type:        "ORDER_CREATED",
		PreviousState: "",
		CurrentState: OrderLifecycleStateCreated,
		Timestamp:   now,
		Data:        map[string]interface{}{"order": order},
	}
	lifecycle.Events = append(lifecycle.Events, event)

	// Store lifecycle
	m.lifecycles[order.ID] = lifecycle

	// Trigger callbacks
	m.triggerCallbacks(lifecycle, event)

	m.logger.Info("Created order lifecycle",
		"orderId", order.ID,
		"state", string(lifecycle.CurrentState),
	)

	return lifecycle, nil
}

// GetLifecycle gets an order lifecycle
func (m *OrderLifecycleManager) GetLifecycle(orderID string) (*OrderLifecycle, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	lifecycle, exists := m.lifecycles[orderID]
	if !exists {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Order lifecycle not found for order ID %s", orderID),
			nil,
			"OrderLifecycleManager",
		)
	}

	return lifecycle, nil
}

// TransitionState transitions an order to a new state
func (m *OrderLifecycleManager) TransitionState(orderID string, newState OrderLifecycleState, eventType string, data map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get lifecycle
	lifecycle, exists := m.lifecycles[orderID]
	if !exists {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Order lifecycle not found for order ID %s", orderID),
			nil,
			"OrderLifecycleManager",
		)
	}

	// Check if transition is valid
	if !isValidTransition(lifecycle.CurrentState, newState) {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidOrder,
			fmt.Sprintf("Invalid state transition from %s to %s", lifecycle.CurrentState, newState),
			nil,
			"OrderLifecycleManager",
		).WithOrderID(orderID)
	}

	// Create event
	now := time.Now()
	event := OrderEvent{
		ID:          generateEventID(),
		OrderID:     orderID,
		Type:        eventType,
		PreviousState: lifecycle.CurrentState,
		CurrentState: newState,
		Timestamp:   now,
		Data:        data,
	}

	// Update lifecycle
	previousState := lifecycle.CurrentState
	lifecycle.CurrentState = newState
	lifecycle.UpdatedAt = now
	lifecycle.Events = append(lifecycle.Events, event)

	// Update order status based on lifecycle state
	lifecycle.Order.Status = mapLifecycleStateToOrderStatus(newState)
	lifecycle.Order.UpdatedAt = now

	// Trigger callbacks
	m.triggerCallbacks(lifecycle, event)

	m.logger.Info("Transitioned order state",
		"orderId", orderID,
		"previousState", string(previousState),
		"newState", string(newState),
		"eventType", eventType,
	)

	return nil
}

// RegisterCallback registers a callback for a state transition
func (m *OrderLifecycleManager) RegisterCallback(state OrderLifecycleState, callback OrderLifecycleCallback) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.callbacks[state]; !exists {
		m.callbacks[state] = make([]OrderLifecycleCallback, 0)
	}

	m.callbacks[state] = append(m.callbacks[state], callback)
}

// triggerCallbacks triggers callbacks for a state transition
func (m *OrderLifecycleManager) triggerCallbacks(lifecycle *OrderLifecycle, event OrderEvent) {
	if callbacks, exists := m.callbacks[event.CurrentState]; exists {
		for _, callback := range callbacks {
			go func(cb OrderLifecycleCallback, lc *OrderLifecycle, ev OrderEvent) {
				defer func() {
					if r := recover(); r != nil {
						m.logger.Error("Panic in order lifecycle callback",
							"orderId", lc.Order.ID,
							"state", string(lc.CurrentState),
							"panic", r,
						)
					}
				}()
				cb(lc, ev)
			}(callback, lifecycle, event)
		}
	}
}

// GetOrderEvents gets all events for an order
func (m *OrderLifecycleManager) GetOrderEvents(orderID string) ([]OrderEvent, error) {
	lifecycle, err := m.GetLifecycle(orderID)
	if err != nil {
		return nil, err
	}

	return lifecycle.Events, nil
}

// GetActiveLifecycles gets all active order lifecycles
func (m *OrderLifecycleManager) GetActiveLifecycles() []*OrderLifecycle {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	activeLifecycles := make([]*OrderLifecycle, 0)
	for _, lifecycle := range m.lifecycles {
		if isActiveState(lifecycle.CurrentState) {
			activeLifecycles = append(activeLifecycles, lifecycle)
		}
	}

	return activeLifecycles
}

// CheckExpiredOrders checks for and processes expired orders
func (m *OrderLifecycleManager) CheckExpiredOrders() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for orderID, lifecycle := range m.lifecycles {
		if lifecycle.ExpiresAt != nil && now.After(*lifecycle.ExpiresAt) && isActiveState(lifecycle.CurrentState) {
			// Create event
			event := OrderEvent{
				ID:          generateEventID(),
				OrderID:     orderID,
				Type:        "ORDER_EXPIRED",
				PreviousState: lifecycle.CurrentState,
				CurrentState: OrderLifecycleStateExpired,
				Timestamp:   now,
				Data:        map[string]interface{}{"expiryTime": lifecycle.ExpiresAt},
			}

			// Update lifecycle
			lifecycle.CurrentState = OrderLifecycleStateExpired
			lifecycle.UpdatedAt = now
			lifecycle.Events = append(lifecycle.Events, event)

			// Update order status
			lifecycle.Order.Status = OrderStatusCancelled
			lifecycle.Order.UpdatedAt = now

			// Trigger callbacks
			m.triggerCallbacks(lifecycle, event)

			m.logger.Info("Order expired",
				"orderId", orderID,
				"expiryTime", lifecycle.ExpiresAt,
			)
		}
	}
}

// StartExpiryChecker starts a goroutine to periodically check for expired orders
func (m *OrderLifecycleManager) StartExpiryChecker(ctx context.Context, checkInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.CheckExpiredOrders()
			}
		}
	}()

	m.logger.Info("Started order expiry checker",
		"checkInterval", checkInterval,
	)
}

// isValidTransition checks if a state transition is valid
func isValidTransition(from, to OrderLifecycleState) bool {
	// Define valid transitions
	validTransitions := map[OrderLifecycleState][]OrderLifecycleState{
		OrderLifecycleStateCreated: {
			OrderLifecycleStateValidated,
			OrderLifecycleStateRejected,
			OrderLifecycleStateFailed,
		},
		OrderLifecycleStateValidated: {
			OrderLifecycleStateSubmitted,
			OrderLifecycleStateRejected,
			OrderLifecycleStateFailed,
		},
		OrderLifecycleStateSubmitted: {
			OrderLifecycleStateAcknowledged,
			OrderLifecycleStateRejected,
			OrderLifecycleStateFailed,
			OrderLifecycleStateCancelling,
		},
		OrderLifecycleStateAcknowledged: {
			OrderLifecycleStatePartiallyFilled,
			OrderLifecycleStateCompleted,
			OrderLifecycleStateCancelling,
			OrderLifecycleStateCancelled,
			OrderLifecycleStateRejected,
			OrderLifecycleStateFailed,
			OrderLifecycleStateExpired,
		},
		OrderLifecycleStatePartiallyFilled: {
			OrderLifecycleStateCompleted,
			OrderLifecycleStateCancelling,
			OrderLifecycleStateCancelled,
			OrderLifecycleStateFailed,
			OrderLifecycleStateExpired,
		},
		OrderLifecycleStateCancelling: {
			OrderLifecycleStateCancelled,
			OrderLifecycleStateFailed,
		},
	}

	// Check if transition is valid
	if validStates, exists := validTransitions[from]; exists {
		for _, validState := range validStates {
			if validState == to {
				return true
			}
		}
	}

	return false
}

// isActiveState checks if a state is considered active
func isActiveState(state OrderLifecycleState) bool {
	activeStates := map[OrderLifecycleState]bool{
		OrderLifecycleStateCreated:     true,
		OrderLifecycleStateValidated:   true,
		OrderLifecycleStateSubmitted:   true,
		OrderLifecycleStateAcknowledged: true,
		OrderLifecycleStatePartiallyFilled: true,
	}

	return activeStates[state]
}

// mapLifecycleStateToOrderStatus maps a lifecycle state to an order status
func mapLifecycleStateToOrderStatus(state OrderLifecycleState) OrderStatus {
	stateToStatus := map[OrderLifecycleState]OrderStatus{
		OrderLifecycleStateCreated:     OrderStatusPending,
		OrderLifecycleStateValidated:   OrderStatusPending,
		OrderLifecycleStateSubmitted:   OrderStatusSubmitted,
		OrderLifecycleStateAcknowledged: OrderStatusOpen,
		OrderLifecycleStatePartiallyFilled: OrderStatusPartial,
		OrderLifecycleStateCompleted:   OrderStatusCompleted,
		OrderLifecycleStateCancelling:  OrderStatusOpen,
		OrderLifecycleStateCancelled:   OrderStatusCancelled,
		OrderLifecycleStateRejected:    OrderStatusRejected,
		OrderLifecycleStateExpired:     OrderStatusCancelled,
		OrderLifecycleStateFailed:      OrderStatusFailed,
	}

	return stateToStatus[state]
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// OrderDependency represents a dependency between orders
type OrderDependency struct {
	ID            string    `json:"id"`
	ParentOrderID string    `json:"parentOrderId"`
	ChildOrderID  string    `json:"childOrderId"`
	Type          string    `json:"type"`
	Condition     string    `json:"condition"`
	CreatedAt     time.Time `json:"createdAt"`
}

// OrderDependencyManager manages dependencies between orders
type OrderDependencyManager struct {
	dependencies map[string][]OrderDependency // parentOrderID -> dependencies
	childToParent map[string]string          // childOrderID -> parentOrderID
	mutex        sync.RWMutex
	lifecycleManager *OrderLifecycleManager
	logger       Logger
}

// NewOrderDependencyManager creates a new order dependency manager
func NewOrderDependencyManager(lifecycleManager *OrderLifecycleManager, logger Logger) *OrderDependencyManager {
	manager := &OrderDependencyManager{
		dependencies: make(map[string][]OrderDependency),
		childToParent: make(map[string]string),
		lifecycleManager: lifecycleManager,
		logger:       logger,
	}

	// Register callbacks for parent order state changes
	lifecycleManager.RegisterCallback(OrderLifecycleStateCompleted, manager.handleParentOrderCompleted)
	lifecycleManager.RegisterCallback(OrderLifecycleStateCancelled, manager.handleParentOrderCancelled)
	lifecycleManager.RegisterCallback(OrderLifecycleStateRejected, manager.handleParentOrderRejected)
	lifecycleManager.RegisterCallback(OrderLifecycleStateFailed, manager.handleParentOrderFailed)

	return manager
}

// CreateDependency creates a new order dependency
func (m *OrderDependencyManager) CreateDependency(parentOrderID, childOrderID, dependencyType, condition string) (*OrderDependency, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if parent order exists
	_, err := m.lifecycleManager.GetLifecycle(parentOrderID)
	if err != nil {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Parent order not found: %s", parentOrderID),
			err,
			"OrderDependencyManager",
		)
	}

	// Check if child order exists
	_, err = m.lifecycleManager.GetLifecycle(childOrderID)
	if err != nil {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeOrderNotFound,
			fmt.Sprintf("Child order not found: %s", childOrderID),
			err,
			"OrderDependencyManager",
		)
	}

	// Check if dependency already exists
	if m.dependencyExists(parentOrderID, childOrderID) {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Dependency already exists between parent %s and child %s", parentOrderID, childOrderID),
			nil,
			"OrderDependencyManager",
		)
	}

	// Create dependency
	dependency := OrderDependency{
		ID:            fmt.Sprintf("dep_%d", time.Now().UnixNano()),
		ParentOrderID: parentOrderID,
		ChildOrderID:  childOrderID,
		Type:          dependencyType,
		Condition:     condition,
		CreatedAt:     time.Now(),
	}

	// Store dependency
	if _, exists := m.dependencies[parentOrderID]; !exists {
		m.dependencies[parentOrderID] = make([]OrderDependency, 0)
	}
	m.dependencies[parentOrderID] = append(m.dependencies[parentOrderID], dependency)
	m.childToParent[childOrderID] = parentOrderID

	m.logger.Info("Created order dependency",
		"dependencyId", dependency.ID,
		"parentOrderId", parentOrderID,
		"childOrderId", childOrderID,
		"type", dependencyType,
		"condition", condition,
	)

	return &dep
(Content truncated due to size limit. Use line ranges to read in chunks)