package strategy

import (
	"errors"
	"sync"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services/order"
)

// StrategyExecutionEngine handles the execution of trading strategies
type StrategyExecutionEngine struct {
	strategyService StrategyService
	orderService    order.OrderService
	activeStrategies map[string]bool
	mutex           sync.RWMutex
}

// NewStrategyExecutionEngine creates a new StrategyExecutionEngine
func NewStrategyExecutionEngine(strategyService StrategyService, orderService order.OrderService) *StrategyExecutionEngine {
	return &StrategyExecutionEngine{
		strategyService: strategyService,
		orderService:    orderService,
		activeStrategies: make(map[string]bool),
		mutex:           sync.RWMutex{},
	}
}

// StartEngine starts the strategy execution engine
func (e *StrategyExecutionEngine) StartEngine() error {
	// Start the scheduler
	go e.runScheduler()
	
	// Start the monitoring service
	go e.monitorActiveStrategies()
	
	return nil
}

// ExecuteStrategy executes a specific strategy
func (e *StrategyExecutionEngine) ExecuteStrategy(strategyID string) error {
	// Get the strategy
	strategy, err := e.strategyService.GetStrategyByID(strategyID)
	if err != nil {
		return err
	}
	
	// Check if strategy is already active
	e.mutex.RLock()
	if e.activeStrategies[strategyID] {
		e.mutex.RUnlock()
		return errors.New("strategy is already being executed")
	}
	e.mutex.RUnlock()
	
	// Mark strategy as active
	e.mutex.Lock()
	e.activeStrategies[strategyID] = true
	e.mutex.Unlock()
	
	// Execute the strategy in a goroutine
	go func() {
		defer func() {
			// Mark strategy as inactive when done
			e.mutex.Lock()
			delete(e.activeStrategies, strategyID)
			e.mutex.Unlock()
		}()
		
		// Update strategy status to active
		strategy.Status = models.StrategyStatusActive
		strategy.LastExecutedAt = time.Now()
		e.strategyService.UpdateStrategy(strategy)
		
		// Process entry conditions
		e.processEntryConditions(strategy)
		
		// Process exit conditions for existing positions
		e.processExitConditions(strategy)
	}()
	
	return nil
}

// StopStrategy stops the execution of a strategy
func (e *StrategyExecutionEngine) StopStrategy(strategyID string) error {
	// Check if strategy is active
	e.mutex.RLock()
	if !e.activeStrategies[strategyID] {
		e.mutex.RUnlock()
		return errors.New("strategy is not active")
	}
	e.mutex.RUnlock()
	
	// Get the strategy
	strategy, err := e.strategyService.GetStrategyByID(strategyID)
	if err != nil {
		return err
	}
	
	// Update strategy status to stopped
	strategy.Status = models.StrategyStatusStopped
	strategy.UpdatedAt = time.Now()
	_, err = e.strategyService.UpdateStrategy(strategy)
	if err != nil {
		return err
	}
	
	// Mark strategy as inactive
	e.mutex.Lock()
	delete(e.activeStrategies, strategyID)
	e.mutex.Unlock()
	
	return nil
}

// runScheduler runs the strategy scheduler
func (e *StrategyExecutionEngine) runScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			e.checkScheduledStrategies()
		}
	}
}

// checkScheduledStrategies checks for strategies that need to be executed based on their schedule
func (e *StrategyExecutionEngine) checkScheduledStrategies() {
	// This is a simplified implementation
	// In a real system, you would query the database for strategies with schedules
	// that match the current time
	
	// For now, we'll just log that we're checking
	// log.Println("Checking scheduled strategies")
}

// monitorActiveStrategies monitors active strategies for stop conditions
func (e *StrategyExecutionEngine) monitorActiveStrategies() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			e.checkActiveStrategies()
		}
	}
}

// checkActiveStrategies checks active strategies for stop conditions
func (e *StrategyExecutionEngine) checkActiveStrategies() {
	// Get a copy of the active strategies map
	e.mutex.RLock()
	activeStrategiesCopy := make(map[string]bool)
	for id, active := range e.activeStrategies {
		activeStrategiesCopy[id] = active
	}
	e.mutex.RUnlock()
	
	// Check each active strategy
	for strategyID := range activeStrategiesCopy {
		// Get the strategy
		strategy, err := e.strategyService.GetStrategyByID(strategyID)
		if err != nil {
			continue
		}
		
		// Check if strategy should be stopped
		if strategy.Status != models.StrategyStatusActive {
			e.mutex.Lock()
			delete(e.activeStrategies, strategyID)
			e.mutex.Unlock()
			continue
		}
		
		// Check risk parameters
		e.checkRiskParameters(strategy)
	}
}

// processEntryConditions processes the entry conditions of a strategy
func (e *StrategyExecutionEngine) processEntryConditions(strategy *models.Strategy) {
	// This is a simplified implementation
	// In a real system, you would evaluate each condition against market data
	// and create orders when conditions are met
	
	// For demonstration purposes, we'll create a sample order
	if len(strategy.EntryConditions) > 0 && len(strategy.Instruments) > 0 {
		// Create a sample order
		order := &models.Order{
			UserID:      strategy.UserID,
			StrategyID:  strategy.ID,
			Symbol:      strategy.Instruments[0],
			OrderType:   models.OrderTypeLimit,
			Side:        models.OrderSideBuy,
			Quantity:    1,
			Price:       100.0, // Sample price
			Status:      models.OrderStatusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		// Submit the order
		_, err := e.orderService.CreateOrder(order)
		if err != nil {
			// Handle error
			// log.Printf("Error creating order: %v", err)
		}
	}
}

// processExitConditions processes the exit conditions of a strategy
func (e *StrategyExecutionEngine) processExitConditions(strategy *models.Strategy) {
	// This is a simplified implementation
	// In a real system, you would evaluate each condition against market data
	// and create exit orders when conditions are met
	
	// For demonstration purposes, we'll just log that we're processing exit conditions
	// log.Printf("Processing exit conditions for strategy %s", strategy.ID)
}

// checkRiskParameters checks if a strategy has exceeded its risk parameters
func (e *StrategyExecutionEngine) checkRiskParameters(strategy *models.Strategy) {
	// This is a simplified implementation
	// In a real system, you would calculate current positions, P&L, etc.
	// and stop the strategy if risk parameters are exceeded
	
	// For demonstration purposes, we'll just log that we're checking risk parameters
	// log.Printf("Checking risk parameters for strategy %s", strategy.ID)
}
