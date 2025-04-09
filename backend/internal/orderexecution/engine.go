package orderexecution

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

// OrderType represents the type of order
type OrderType string

// Order types
const (
	Market  OrderType = "MARKET"
	Limit   OrderType = "LIMIT"
	StopLoss OrderType = "SL"
	StopLossMarket OrderType = "SL-M"
)

// TransactionType represents buy or sell
type TransactionType string

// Transaction types
const (
	Buy  TransactionType = "BUY"
	Sell TransactionType = "SELL"
)

// OrderStatus represents the status of an order
type OrderStatus string

// Order statuses
const (
	Pending    OrderStatus = "PENDING"
	Open       OrderStatus = "OPEN"
	Executed   OrderStatus = "EXECUTED"
	Cancelled  OrderStatus = "CANCELLED"
	Rejected   OrderStatus = "REJECTED"
	PartiallyExecuted OrderStatus = "PARTIALLY_EXECUTED"
)

// ProductType represents the product type for the order
type ProductType string

// Product types
const (
	Intraday ProductType = "MIS"
	Normal   ProductType = "NRML"
	CashAndCarry ProductType = "CNC"
)

// ValidityType represents how long the order is valid
type ValidityType string

// Validity types
const (
	Day ValidityType = "DAY"
	IOC ValidityType = "IOC"
	GTC ValidityType = "GTC"
)

// Order represents an order in the system
type Order struct {
	ID              string          `json:"id"`
	Symbol          string          `json:"symbol"`
	Quantity        int             `json:"quantity"`
	Price           float64         `json:"price"`
	OrderType       OrderType       `json:"orderType"`
	TransactionType TransactionType `json:"transactionType"`
	Status          OrderStatus     `json:"status"`
	FilledQuantity  int             `json:"filledQuantity"`
	AveragePrice    float64         `json:"averagePrice"`
	PlacedAt        time.Time       `json:"placedAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	Validity        ValidityType    `json:"validity"`
	TriggerPrice    float64         `json:"triggerPrice,omitempty"`
	Exchange        string          `json:"exchange"`
	Product         ProductType     `json:"product"`
	Message         string          `json:"message,omitempty"`
	BrokerOrderID   string          `json:"brokerOrderID,omitempty"`
	ParentOrderID   string          `json:"parentOrderID,omitempty"`
	StrategyID      string          `json:"strategyID,omitempty"`
	Tags            []string        `json:"tags,omitempty"`
}

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Symbol          string          `json:"symbol"`
	Quantity        int             `json:"quantity"`
	Price           float64         `json:"price"`
	OrderType       OrderType       `json:"orderType"`
	TransactionType TransactionType `json:"transactionType"`
	Validity        ValidityType    `json:"validity"`
	TriggerPrice    float64         `json:"triggerPrice,omitempty"`
	Exchange        string          `json:"exchange"`
	Product         ProductType     `json:"product"`
	StrategyID      string          `json:"strategyID,omitempty"`
	Tags            []string        `json:"tags,omitempty"`
}

// OrderResponse represents a response after placing an order
type OrderResponse struct {
	Order  *Order `json:"order"`
	Error  string `json:"error,omitempty"`
	Status bool   `json:"status"`
}

// OrderUpdateCallback is a function that gets called when an order is updated
type OrderUpdateCallback func(order *Order)

// BrokerAdapter defines the interface for broker-specific implementations
type BrokerAdapter interface {
	PlaceOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error)
	ModifyOrder(ctx context.Context, orderID string, request *OrderRequest) (*OrderResponse, error)
	CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error)
	GetOrderStatus(ctx context.Context, orderID string) (*Order, error)
	GetOrders(ctx context.Context) ([]*Order, error)
}

// SmartRouter is responsible for routing orders to the appropriate broker
type SmartRouter interface {
	RouteOrder(ctx context.Context, request *OrderRequest) (BrokerAdapter, error)
}

// OrderExecutionEngine is the main engine for executing orders
type OrderExecutionEngine struct {
	brokers       map[string]BrokerAdapter
	smartRouter   SmartRouter
	orders        map[string]*Order
	ordersMutex   sync.RWMutex
	callbacks     []OrderUpdateCallback
	callbackMutex sync.RWMutex
}

// NewOrderExecutionEngine creates a new order execution engine
func NewOrderExecutionEngine(smartRouter SmartRouter) *OrderExecutionEngine {
	return &OrderExecutionEngine{
		brokers:     make(map[string]BrokerAdapter),
		smartRouter: smartRouter,
		orders:      make(map[string]*Order),
	}
}

// RegisterBroker registers a broker adapter with the engine
func (e *OrderExecutionEngine) RegisterBroker(name string, broker BrokerAdapter) {
	e.brokers[name] = broker
}

// RegisterCallback registers a callback for order updates
func (e *OrderExecutionEngine) RegisterCallback(callback OrderUpdateCallback) {
	e.callbackMutex.Lock()
	defer e.callbackMutex.Unlock()
	e.callbacks = append(e.callbacks, callback)
}

// ExecuteOrder executes an order using the smart router
func (e *OrderExecutionEngine) ExecuteOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	// Use smart router to determine the best broker for this order
	broker, err := e.smartRouter.RouteOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	// Place the order with the selected broker
	response, err := broker.PlaceOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	// Store the order in our local cache
	if response.Status && response.Order != nil {
		e.ordersMutex.Lock()
		e.orders[response.Order.ID] = response.Order
		e.ordersMutex.Unlock()

		// Notify callbacks
		e.notifyOrderUpdate(response.Order)
	}

	return response, nil
}

// ExecuteOrderBatch executes multiple orders in a batch
func (e *OrderExecutionEngine) ExecuteOrderBatch(ctx context.Context, requests []*OrderRequest) ([]*OrderResponse, error) {
	responses := make([]*OrderResponse, len(requests))
	var wg sync.WaitGroup
	
	// Create a channel to limit concurrent order executions
	// This helps with high-volume order execution by controlling the rate
	semaphore := make(chan struct{}, 50) // Allow up to 50 concurrent order executions
	
	for i, request := range requests {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		
		go func(idx int, req *OrderRequest) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore
			
			response, err := e.ExecuteOrder(ctx, req)
			if err != nil {
				responses[idx] = &OrderResponse{
					Status: false,
					Error:  err.Error(),
				}
				return
			}
			
			responses[idx] = response
		}(i, request)
	}
	
	wg.Wait()
	return responses, nil
}

// ModifyOrder modifies an existing order
func (e *OrderExecutionEngine) ModifyOrder(ctx context.Context, orderID string, request *OrderRequest) (*OrderResponse, error) {
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		return nil, errors.New("order not found")
	}
	
	// Find the broker that placed this order
	broker, exists := e.brokers[order.Exchange]
	if !exists {
		return nil, errors.New("broker not found for this order")
	}
	
	// Modify the order with the broker
	response, err := broker.ModifyOrder(ctx, orderID, request)
	if err != nil {
		return nil, err
	}
	
	// Update our local cache
	if response.Status && response.Order != nil {
		e.ordersMutex.Lock()
		e.orders[response.Order.ID] = response.Order
		e.ordersMutex.Unlock()
		
		// Notify callbacks
		e.notifyOrderUpdate(response.Order)
	}
	
	return response, nil
}

// CancelOrder cancels an existing order
func (e *OrderExecutionEngine) CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		return nil, errors.New("order not found")
	}
	
	// Find the broker that placed this order
	broker, exists := e.brokers[order.Exchange]
	if !exists {
		return nil, errors.New("broker not found for this order")
	}
	
	// Cancel the order with the broker
	response, err := broker.CancelOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// Update our local cache
	if response.Status && response.Order != nil {
		e.ordersMutex.Lock()
		e.orders[response.Order.ID] = response.Order
		e.ordersMutex.Unlock()
		
		// Notify callbacks
		e.notifyOrderUpdate(response.Order)
	}
	
	return response, nil
}

// GetOrder retrieves an order by ID
func (e *OrderExecutionEngine) GetOrder(orderID string) (*Order, error) {
	e.ordersMutex.RLock()
	defer e.ordersMutex.RUnlock()
	
	order, exists := e.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	
	return order, nil
}

// GetOrders retrieves all orders
func (e *OrderExecutionEngine) GetOrders() []*Order {
	e.ordersMutex.RLock()
	defer e.ordersMutex.RUnlock()
	
	orders := make([]*Order, 0, len(e.orders))
	for _, order := range e.orders {
		orders = append(orders, order)
	}
	
	return orders
}

// SyncOrderStatus synchronizes the order status with the broker
func (e *OrderExecutionEngine) SyncOrderStatus(ctx context.Context, orderID string) error {
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		return errors.New("order not found")
	}
	
	// Find the broker that placed this order
	broker, exists := e.brokers[order.Exchange]
	if !exists {
		return errors.New("broker not found for this order")
	}
	
	// Get the latest order status from the broker
	updatedOrder, err := broker.GetOrderStatus(ctx, orderID)
	if err != nil {
		return err
	}
	
	// Update our local cache
	e.ordersMutex.Lock()
	e.orders[orderID] = updatedOrder
	e.ordersMutex.Unlock()
	
	// Notify callbacks
	e.notifyOrderUpdate(updatedOrder)
	
	return nil
}

// SyncAllOrders synchronizes all orders with their respective brokers
func (e *OrderExecutionEngine) SyncAllOrders(ctx context.Context) error {
	var wg sync.WaitGroup
	var errMutex sync.Mutex
	var syncErrors []error
	
	// Get a snapshot of all order IDs
	e.ordersMutex.RLock()
	orderIDs := make([]string, 0, len(e.orders))
	for id := range e.orders {
		orderIDs = append(orderIDs, id)
	}
	e.ordersMutex.RUnlock()
	
	// Create a channel to limit concurrent sync operations
	semaphore := make(chan struct{}, 20) // Allow up to 20 concurrent sync operations
	
	for _, id := range orderIDs {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		
		go func(orderID string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore
			
			err := e.SyncOrderStatus(ctx, orderID)
			if err != nil {
				errMutex.Lock()
				syncErrors = append(syncErrors, err)
				errMutex.Unlock()
				log.Printf("Error syncing order %s: %v", orderID, err)
			}
		}(id)
	}
	
	wg.Wait()
	
	if len(syncErrors) > 0 {
		return errors.New("one or more errors occurred while syncing orders")
	}
	
	return nil
}

// notifyOrderUpdate notifies all registered callbacks about an order update
func (e *OrderExecutionEngine) notifyOrderUpdate(order *Order) {
	e.callbackMutex.RLock()
	defer e.callbackMutex.RUnlock()
	
	for _, callback := range e.callbacks {
		go callback(order)
	}
}
