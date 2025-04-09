package orderexecution

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// OrderStatus constants defined in engine.go

// OrderMonitoringSystem is responsible for monitoring order status
type OrderMonitoringSystem struct {
	engine            *OrderExecutionEngine
	monitoredOrders   map[string]*MonitoredOrder
	mutex             sync.RWMutex
	isRunning         bool
	stopChan          chan struct{}
	wg                sync.WaitGroup
	updateInterval    time.Duration
	statusCallbacks   map[string][]OrderStatusCallback
	callbacksMutex    sync.RWMutex
	alertThresholds   map[string]AlertThreshold
	alertsMutex       sync.RWMutex
	alertCallbacks    []AlertCallback
	alertCallbackMutex sync.RWMutex
}

// MonitoredOrder represents an order being monitored
type MonitoredOrder struct {
	Order           *Order
	LastUpdated     time.Time
	MonitoringStart time.Time
	ExpectedTime    time.Duration
	Alerts          []Alert
}

// OrderStatusCallback is a function that gets called when an order status changes
type OrderStatusCallback func(order *Order, previousStatus OrderStatus)

// AlertType represents the type of alert
type AlertType string

const (
	// AlertTypeDelayed indicates an order is taking longer than expected
	AlertTypeDelayed AlertType = "DELAYED"
	
	// AlertTypeRejected indicates an order was rejected
	AlertTypeRejected AlertType = "REJECTED"
	
	// AlertTypePartialFill indicates an order was partially filled
	AlertTypePartialFill AlertType = "PARTIAL_FILL"
	
	// AlertTypePriceDeviation indicates a significant price deviation
	AlertTypePriceDeviation AlertType = "PRICE_DEVIATION"
)

// Alert represents an alert for an order
type Alert struct {
	Type        AlertType
	OrderID     string
	Message     string
	CreatedAt   time.Time
	Acknowledged bool
}

// AlertThreshold defines thresholds for generating alerts
type AlertThreshold struct {
	DelayThreshold      time.Duration
	PriceDeviationPct   float64
	PartialFillDuration time.Duration
}

// AlertCallback is a function that gets called when an alert is generated
type AlertCallback func(alert Alert)

// NewOrderMonitoringSystem creates a new order monitoring system
func NewOrderMonitoringSystem(engine *OrderExecutionEngine, updateInterval time.Duration) *OrderMonitoringSystem {
	return &OrderMonitoringSystem{
		engine:          engine,
		monitoredOrders: make(map[string]*MonitoredOrder),
		stopChan:        make(chan struct{}),
		updateInterval:  updateInterval,
		statusCallbacks: make(map[string][]OrderStatusCallback),
		alertThresholds: make(map[string]AlertThreshold),
	}
}

// Start starts the order monitoring system
func (s *OrderMonitoringSystem) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.isRunning {
		return
	}
	
	s.isRunning = true
	s.stopChan = make(chan struct{})
	
	// Start the monitoring goroutine
	s.wg.Add(1)
	go s.monitorOrders()
	
	log.Println("Order monitoring system started")
}

// Stop stops the order monitoring system
func (s *OrderMonitoringSystem) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if !s.isRunning {
		return
	}
	
	close(s.stopChan)
	s.wg.Wait()
	s.isRunning = false
	
	log.Println("Order monitoring system stopped")
}

// monitorOrders periodically checks the status of monitored orders
func (s *OrderMonitoringSystem) monitorOrders() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(s.updateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.updateOrderStatuses()
			s.checkForAlerts()
		}
	}
}

// updateOrderStatuses updates the status of all monitored orders
func (s *OrderMonitoringSystem) updateOrderStatuses() {
	// Get a snapshot of monitored order IDs
	s.mutex.RLock()
	orderIDs := make([]string, 0, len(s.monitoredOrders))
	for id := range s.monitoredOrders {
		orderIDs = append(orderIDs, id)
	}
	s.mutex.RUnlock()
	
	// Update each order
	for _, id := range orderIDs {
		s.updateOrderStatus(id)
	}
}

// updateOrderStatus updates the status of a single order
func (s *OrderMonitoringSystem) updateOrderStatus(orderID string) {
	// Get the current order from the engine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := s.engine.SyncOrderStatus(ctx, orderID)
	if err != nil {
		log.Printf("Error syncing order status for %s: %v", orderID, err)
		return
	}
	
	// Get the updated order
	updatedOrder, err := s.engine.GetOrder(orderID)
	if err != nil {
		log.Printf("Error getting updated order %s: %v", orderID, err)
		return
	}
	
	// Get the monitored order
	s.mutex.Lock()
	monitoredOrder, exists := s.monitoredOrders[orderID]
	if !exists {
		s.mutex.Unlock()
		return
	}
	
	// Check if status has changed
	previousStatus := monitoredOrder.Order.Status
	statusChanged := previousStatus != updatedOrder.Status
	
	// Update the monitored order
	monitoredOrder.Order = updatedOrder
	monitoredOrder.LastUpdated = time.Now()
	s.mutex.Unlock()
	
	// Notify callbacks if status changed
	if statusChanged {
		s.notifyStatusCallbacks(updatedOrder, previousStatus)
	}
	
	// Check if we should stop monitoring this order
	if isTerminalStatus(updatedOrder.Status) {
		s.StopMonitoring(orderID)
	}
}

// isTerminalStatus checks if an order status is terminal
func isTerminalStatus(status OrderStatus) bool {
	return status == Executed || status == Cancelled || status == Rejected
}

// notifyStatusCallbacks notifies all registered callbacks about a status change
func (s *OrderMonitoringSystem) notifyStatusCallbacks(order *Order, previousStatus OrderStatus) {
	s.callbacksMutex.RLock()
	defer s.callbacksMutex.RUnlock()
	
	// Call order-specific callbacks
	if callbacks, exists := s.statusCallbacks[order.ID]; exists {
		for _, callback := range callbacks {
			go callback(order, previousStatus)
		}
	}
}

// checkForAlerts checks for alert conditions in monitored orders
func (s *OrderMonitoringSystem) checkForAlerts() {
	now := time.Now()
	
	// Get a snapshot of monitored orders
	s.mutex.RLock()
	orders := make(map[string]*MonitoredOrder)
	for id, order := range s.monitoredOrders {
		orders[id] = order
	}
	s.mutex.RUnlock()
	
	// Check each order for alert conditions
	for id, monitoredOrder := range orders {
		// Skip orders that are already in a terminal state
		if isTerminalStatus(monitoredOrder.Order.Status) {
			continue
		}
		
		// Get alert threshold for this order
		s.alertsMutex.RLock()
		threshold, exists := s.alertThresholds[id]
		if !exists {
			// Use default threshold
			threshold = AlertThreshold{
				DelayThreshold:      30 * time.Second,
				PriceDeviationPct:   5.0,
				PartialFillDuration: 60 * time.Second,
			}
		}
		s.alertsMutex.RUnlock()
		
		// Check for delay
		if threshold.DelayThreshold > 0 {
			elapsed := now.Sub(monitoredOrder.MonitoringStart)
			if elapsed > monitoredOrder.ExpectedTime+threshold.DelayThreshold {
				s.generateAlert(Alert{
					Type:        AlertTypeDelayed,
					OrderID:     id,
					Message:     fmt.Sprintf("Order %s is delayed. Expected: %v, Elapsed: %v", id, monitoredOrder.ExpectedTime, elapsed),
					CreatedAt:   now,
					Acknowledged: false,
				})
			}
		}
		
		// Check for partial fill that's taking too long
		if threshold.PartialFillDuration > 0 && monitoredOrder.Order.Status == PartiallyExecuted {
			elapsed := now.Sub(monitoredOrder.LastUpdated)
			if elapsed > threshold.PartialFillDuration {
				s.generateAlert(Alert{
					Type:        AlertTypePartialFill,
					OrderID:     id,
					Message:     fmt.Sprintf("Order %s is partially filled for too long. Filled: %d/%d, Elapsed: %v", 
						id, monitoredOrder.Order.FilledQuantity, monitoredOrder.Order.Quantity, elapsed),
					CreatedAt:   now,
					Acknowledged: false,
				})
			}
		}
		
		// Check for price deviation (for limit orders)
		if threshold.PriceDeviationPct > 0 && monitoredOrder.Order.OrderType == Limit && monitoredOrder.Order.FilledQuantity > 0 {
			expectedPrice := monitoredOrder.Order.Price
			actualPrice := monitoredOrder.Order.AveragePrice
			deviation := ((actualPrice - expectedPrice) / expectedPrice) * 100
			
			if deviation < 0 {
				deviation = -deviation // Get absolute value
			}
			
			if deviation > threshold.PriceDeviationPct {
				s.generateAlert(Alert{
					Type:        AlertTypePriceDeviation,
					OrderID:     id,
					Message:     fmt.Sprintf("Order %s has price deviation of %.2f%%. Expected: %.2f, Actual: %.2f", 
						id, deviation, expectedPrice, actualPrice),
					CreatedAt:   now,
					Acknowledged: false,
				})
			}
		}
		
		// Check for rejected orders
		if monitoredOrder.Order.Status == Rejected {
			s.generateAlert(Alert{
				Type:        AlertTypeRejected,
				OrderID:     id,
				Message:     fmt.Sprintf("Order %s was rejected. Reason: %s", id, monitoredOrder.Order.Message),
				CreatedAt:   now,
				Acknowledged: false,
			})
		}
	}
}

// generateAlert generates an alert and notifies callbacks
func (s *OrderMonitoringSystem) generateAlert(alert Alert) {
	// Add the alert to the monitored order
	s.mutex.Lock()
	if monitoredOrder, exists := s.monitoredOrders[alert.OrderID]; exists {
		// Check if we already have a similar alert
		for _, existingAlert := range monitoredOrder.Alerts {
			if existingAlert.Type == alert.Type && !existingAlert.Acknowledged {
				// Skip duplicate alerts
				s.mutex.Unlock()
				return
			}
		}
		
		monitoredOrder.Alerts = append(monitoredOrder.Alerts, alert)
	}
	s.mutex.Unlock()
	
	// Notify alert callbacks
	s.alertCallbackMutex.RLock()
	defer s.alertCallbackMutex.RUnlock()
	
	for _, callback := range s.alertCallbacks {
		go callback(alert)
	}
	
	log.Printf("Alert generated: %s - %s", alert.Type, alert.Message)
}

// StartMonitoring starts monitoring an order
func (s *OrderMonitoringSystem) StartMonitoring(orderID string, expectedTime time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if !s.isRunning {
		return errors.New("order monitoring system is not running")
	}
	
	// Check if we're already monitoring this order
	if _, exists := s.monitoredOrders[orderID]; exists {
		return nil // Already monitoring
	}
	
	// Get the order from the engine
	order, err := s.engine.GetOrder(orderID)
	if err != nil {
		return err
	}
	
	// Start monitoring the order
	s.monitoredOrders[orderID] = &MonitoredOrder{
		Order:           order,
		LastUpdated:     time.Now(),
		MonitoringStart: time.Now(),
		ExpectedTime:    expectedTime,
		Alerts:          make([]Alert, 0),
	}
	
	log.Printf("Started monitoring order %s", orderID)
	return nil
}

// StopMonitoring stops monitoring an order
func (s *OrderMonitoringSystem) StopMonitoring(orderID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	delete(s.monitoredOrders, orderID)
	
	log.Printf("Stopped monitoring order %s", orderID)
}

// GetMonitoredOrder returns a monitored order
func (s *OrderMonitoringSystem) GetMonitoredOrder(orderID string) (*MonitoredOrder, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	monitoredOrder, exists := s.monitoredOrders[orderID]
	if !exists {
		return nil, errors.New("order not being monitored")
	}
	
	return monitoredOrder, nil
}

// GetAllMonitoredOrders returns all monitored orders
func (s *OrderMonitoringSystem) GetAllMonitoredOrders() map[string]*MonitoredOrder {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// Create a copy to avoid concurrent access issues
	ordersCopy := make(map[string]*MonitoredOrder)
	for id, order := range s.monitoredOrders {
		ordersCopy[id] = &MonitoredOrder{
			Order:           order.Order,
			LastUpdated:     order.LastUpdated,
			MonitoringStart: order.MonitoringStart,
			ExpectedTime:    order.ExpectedTime,
			Alerts:          append([]Alert{}, order.Alerts...),
		}
	}
	
	return ordersCopy
}

// RegisterStatusCallback registers a callback for order status changes
func (s *OrderMonitoringSystem) RegisterStatusCallback(orderID string, callback OrderStatusCallback) {
	s.callbacksMutex.Lock()
	defer s.callbacksMutex.Unlock()
	
	if _, exists := s.statusCallbacks[orderID]; !exists {
		s.statusCallbacks[orderID] = make([]OrderStatusCallback, 0)
	}
	
	s.statusCallbacks[orderID] = append(s.statusCallbacks[orderID], callback)
}

// RegisterAlertCallback registers a callback for alerts
func (s *OrderMonitoringSystem) RegisterAlertCallback(callback AlertCallback) {
	s.alertCallbackMutex.Lock()
	defer s.alertCallbackMutex.Unlock()
	
	s.alertCallbacks = append(s.alertCallbacks, callback)
}

// SetAlertThreshold sets the alert threshold for an order
func (s *OrderMonitoringSystem) SetAlertThreshold(orderID string, threshold AlertThreshold) {
	s.alertsMutex.Lock()
	defer s.alertsMutex.Unlock()
	
	s.alertThresholds[orderID] = threshold
}

// AcknowledgeAlert acknowledges an alert
func (s *OrderMonitoringSystem) AcknowledgeAlert(orderID string, alertType AlertType) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	monitoredOrder, exists := s.monitoredOrders[orderID]
	if !exists {
		return errors.New("order not being monitored")
	}
	
	// Find and acknowledge the alert
	for i := range monitoredOrder.Alerts {
		if monitoredOrder.Alerts[i].Type == alertType && !monitoredOrder.Alerts[i].Acknowledged {
			monitoredOrder.Alerts[i].Acknowledged = true
			return nil
		}
	}
	
	return errors.New("alert not found or already acknowledged")
}

// GetAlerts returns all alerts for an order
func (s *OrderMonitoringSystem) GetAlerts(orderID string) ([]Alert, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	monitoredOrder, exists := s.monitoredOrders[orderID]
	if !exists {
		return nil, errors.New("order not being monitored")
	}
	
	// Create a copy of the alerts
	alertsCopy := make([]Alert, len(monitoredOrder.Alerts))
	copy(alertsCopy, monitoredOrder.Alerts)
	
	return alertsCopy, nil
}

// GetAllAlerts returns all alerts for all orders
func (s *OrderMonitoringSystem) GetAllAlerts() []Alert {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var allAlerts []Alert
	
	for _, monitoredOrder := range s.monitoredOrders {
		allAlerts = append(allAlerts, monitoredOrder.Alerts...)
	}
	
	return allAlerts
}

// GetActiveAlerts returns all active (unacknowledged) alerts
func (s *OrderMonitoringSystem) GetActiveAlerts() []Alert {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var activeAlerts []Alert
	
	for _, monitoredOrder := range s.monitoredOrders {
		for _, alert := range monitoredOrder.Alerts {
			if !alert.Acknowledged {
				activeAlerts = append(activeAlerts, alert)
			}
		}
	}
	
	return activeAlerts
}
