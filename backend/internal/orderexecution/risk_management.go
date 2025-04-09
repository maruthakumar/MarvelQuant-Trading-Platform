package orderexecution

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// RiskLimitType represents the type of risk limit
type RiskLimitType string

const (
	// Risk limit types
	RiskLimitTypeOrderValue      RiskLimitType = "ORDER_VALUE"
	RiskLimitTypePositionSize    RiskLimitType = "POSITION_SIZE"
	RiskLimitTypeDrawdown        RiskLimitType = "DRAWDOWN"
	RiskLimitTypeLeverage        RiskLimitType = "LEVERAGE"
	RiskLimitTypeConcentration   RiskLimitType = "CONCENTRATION"
	RiskLimitTypeOrderRate       RiskLimitType = "ORDER_RATE"
	RiskLimitTypeExposure        RiskLimitType = "EXPOSURE"
	RiskLimitTypeVaR             RiskLimitType = "VAR"
	RiskLimitTypeStressTest      RiskLimitType = "STRESS_TEST"
)

// RiskLevel represents the level of risk
type RiskLevel string

const (
	// Risk levels
	RiskLevelLow     RiskLevel = "LOW"
	RiskLevelMedium  RiskLevel = "MEDIUM"
	RiskLevelHigh    RiskLevel = "HIGH"
	RiskLevelExtreme RiskLevel = "EXTREME"
)

// RiskLimit represents a risk limit
type RiskLimit struct {
	Type        RiskLimitType `json:"type"`
	Value       float64       `json:"value"`
	Level       RiskLevel     `json:"level"`
	Description string        `json:"description"`
	Enabled     bool          `json:"enabled"`
}

// RiskProfile represents a risk profile
type RiskProfile struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Limits      map[RiskLimitType]RiskLimit `json:"limits"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

// EnhancedRiskManager implements an enhanced risk management system
type EnhancedRiskManager struct {
	defaultRiskManager *DefaultRiskManager
	riskProfiles       map[string]RiskProfile
	portfolioPositions map[string]map[string]Position // portfolioID -> symbol -> position
	strategyPositions  map[string]map[string]Position // strategyID -> symbol -> position
	orderHistory       map[string][]Order             // portfolioID -> orders
	mutex              sync.RWMutex
	logger             Logger
	errorHandler       ErrorHandler
}

// NewEnhancedRiskManager creates a new enhanced risk manager
func NewEnhancedRiskManager(
	defaultRiskManager *DefaultRiskManager,
	logger Logger,
	errorHandler ErrorHandler,
) *EnhancedRiskManager {
	return &EnhancedRiskManager{
		defaultRiskManager: defaultRiskManager,
		riskProfiles:       make(map[string]RiskProfile),
		portfolioPositions: make(map[string]map[string]Position),
		strategyPositions:  make(map[string]map[string]Position),
		orderHistory:       make(map[string][]Order),
		logger:             logger,
		errorHandler:       errorHandler,
	}
}

// CreateRiskProfile creates a new risk profile
func (r *EnhancedRiskManager) CreateRiskProfile(profile RiskProfile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Validate profile
	if profile.ID == "" {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			"Risk profile ID cannot be empty",
			nil,
			"EnhancedRiskManager",
		)
	}

	// Check if profile already exists
	if _, exists := r.riskProfiles[profile.ID]; exists {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Risk profile with ID %s already exists", profile.ID),
			nil,
			"EnhancedRiskManager",
		)
	}

	// Set timestamps
	now := time.Now()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	// Store profile
	r.riskProfiles[profile.ID] = profile

	r.logger.Info("Created risk profile",
		"profileId", profile.ID,
		"name", profile.Name,
	)

	return nil
}

// GetRiskProfile gets a risk profile by ID
func (r *EnhancedRiskManager) GetRiskProfile(profileID string) (RiskProfile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	profile, exists := r.riskProfiles[profileID]
	if !exists {
		return RiskProfile{}, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Risk profile with ID %s not found", profileID),
			nil,
			"EnhancedRiskManager",
		)
	}

	return profile, nil
}

// UpdateRiskProfile updates a risk profile
func (r *EnhancedRiskManager) UpdateRiskProfile(profile RiskProfile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if profile exists
	existingProfile, exists := r.riskProfiles[profile.ID]
	if !exists {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Risk profile with ID %s not found", profile.ID),
			nil,
			"EnhancedRiskManager",
		)
	}

	// Preserve creation time
	profile.CreatedAt = existingProfile.CreatedAt
	profile.UpdatedAt = time.Now()

	// Update profile
	r.riskProfiles[profile.ID] = profile

	r.logger.Info("Updated risk profile",
		"profileId", profile.ID,
		"name", profile.Name,
	)

	return nil
}

// DeleteRiskProfile deletes a risk profile
func (r *EnhancedRiskManager) DeleteRiskProfile(profileID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if profile exists
	if _, exists := r.riskProfiles[profileID]; !exists {
		return NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Risk profile with ID %s not found", profileID),
			nil,
			"EnhancedRiskManager",
		)
	}

	// Delete profile
	delete(r.riskProfiles, profileID)

	r.logger.Info("Deleted risk profile",
		"profileId", profileID,
	)

	return nil
}

// UpdatePosition updates a position
func (r *EnhancedRiskManager) UpdatePosition(portfolioID string, position Position) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Initialize portfolio positions if needed
	if _, exists := r.portfolioPositions[portfolioID]; !exists {
		r.portfolioPositions[portfolioID] = make(map[string]Position)
	}

	// Update position
	r.portfolioPositions[portfolioID][position.Symbol] = position

	r.logger.Info("Updated position",
		"portfolioId", portfolioID,
		"symbol", position.Symbol,
		"quantity", position.Quantity,
		"averagePrice", position.AveragePrice,
		"pnl", position.PnL,
	)

	return nil
}

// GetPosition gets a position
func (r *EnhancedRiskManager) GetPosition(portfolioID string, symbol string) (Position, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Check if portfolio exists
	portfolioPositions, exists := r.portfolioPositions[portfolioID]
	if !exists {
		return Position{}, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Portfolio with ID %s not found", portfolioID),
			nil,
			"EnhancedRiskManager",
		)
	}

	// Check if position exists
	position, exists := portfolioPositions[symbol]
	if !exists {
		return Position{}, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Position for symbol %s not found in portfolio %s", symbol, portfolioID),
			nil,
			"EnhancedRiskManager",
		)
	}

	return position, nil
}

// GetPortfolioPositions gets all positions for a portfolio
func (r *EnhancedRiskManager) GetPortfolioPositions(portfolioID string) ([]Position, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Check if portfolio exists
	portfolioPositions, exists := r.portfolioPositions[portfolioID]
	if !exists {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Portfolio with ID %s not found", portfolioID),
			nil,
			"EnhancedRiskManager",
		)
	}

	// Convert map to slice
	positions := make([]Position, 0, len(portfolioPositions))
	for _, position := range portfolioPositions {
		positions = append(positions, position)
	}

	return positions, nil
}

// RecordOrder records an order in the order history
func (r *EnhancedRiskManager) RecordOrder(order Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Initialize order history if needed
	if _, exists := r.orderHistory[order.PortfolioID]; !exists {
		r.orderHistory[order.PortfolioID] = make([]Order, 0)
	}

	// Add order to history
	r.orderHistory[order.PortfolioID] = append(r.orderHistory[order.PortfolioID], order)

	r.logger.Info("Recorded order",
		"orderId", order.ID,
		"portfolioId", order.PortfolioID,
		"strategyId", order.StrategyID,
		"symbol", order.Symbol,
		"side", order.Side,
		"quantity", order.Quantity,
		"price", order.Price,
		"status", order.Status,
	)

	return nil
}

// GetOrderHistory gets the order history for a portfolio
func (r *EnhancedRiskManager) GetOrderHistory(portfolioID string) ([]Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Check if portfolio exists
	orderHistory, exists := r.orderHistory[portfolioID]
	if !exists {
		return nil, NewExecutionError(
			ErrorTypeValidation,
			ErrorSeverityError,
			ErrCodeInvalidParameter,
			fmt.Sprintf("Portfolio with ID %s not found", portfolioID),
			nil,
			"EnhancedRiskManager",
		)
	}

	return orderHistory, nil
}

// ValidateOrder validates an order against risk parameters
func (r *EnhancedRiskManager) ValidateOrder(order Order, portfolio Portfolio, strategy Strategy) error {
	// First, use the default risk manager for basic validation
	if err := r.defaultRiskManager.ValidateOrder(order, portfolio, strategy); err != nil {
		return err
	}

	// Get risk profile for the strategy
	var riskProfile RiskProfile
	var err error

	// Check if strategy has a risk profile
	if riskProfileID, ok := strategy.RiskParams["riskProfileID"].(string); ok {
		riskProfile, err = r.GetRiskProfile(riskProfileID)
		if err != nil {
			r.logger.Warn("Risk profile not found, using default validation",
				"strategyId", strategy.ID,
				"riskProfileId", riskProfileID,
			)
		}
	}

	// Perform enhanced validation
	if err := r.validateOrderValue(order, riskProfile); err != nil {
		return err
	}

	if err := r.validatePositionSize(order, portfolio, riskProfile); err != nil {
		return err
	}

	if err := r.validateDrawdown(order, portfolio, riskProfile); err != nil {
		return err
	}

	if err := r.validateLeverage(order, portfolio, riskProfile); err != nil {
		return err
	}

	if err := r.validateConcentration(order, portfolio, riskProfile); err != nil {
		return err
	}

	if err := r.validateOrderRate(order, portfolio, riskProfile); err != nil {
		return err
	}

	if err := r.validateExposure(order, portfolio, strategy, riskProfile); err != nil {
		return err
	}

	if err := r.validateVaR(order, portfolio, strategy, riskProfile); err != nil {
		return err
	}

	if err := r.validateStressTest(order, portfolio, strategy, riskProfile); err != nil {
		return err
	}

	return nil
}

// validateOrderValue validates the order value
func (r *EnhancedRiskManager) validateOrderValue(order Order, riskProfile RiskProfile) error {
	// Calculate order value
	orderValue := float64(order.Quantity) * order.Price

	// Check if order value limit is enabled
	if limit, ok := riskProfile.Limits[RiskLimitTypeOrderValue]; ok && limit.Enabled {
		if orderValue > limit.Value {
			return NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				ErrCodeInvalidOrder,
				fmt.Sprintf("Order value %.2f exceeds limit %.2f", orderValue, limit.Value),
				nil,
				"EnhancedRiskManager",
			).WithOrderID(order.ID).WithPortfolioID(order.PortfolioID).WithStrategyID(order.StrategyID)
		}
	}

	return nil
}

// validatePositionSize validates the position size
func (r *EnhancedRiskManager) validatePositionSize(order Order, portfolio Portfolio, riskProfile RiskProfile) error {
	// Get current position
	position, err := r.GetPosition(portfolio.ID, order.Symbol)
	if err != nil {
		// No existing position, use order quantity as new position size
		position = Position{
			Symbol:   order.Symbol,
			Quantity: 0,
		}
	}

	// Calculate new position size
	newQuantity := position.Quantity
	if order.Side == OrderSideBuy {
		newQuantity += order.Quantity
	} else {
		newQuantity -= order.Quantity
	}

	// Check if position size limit is enabled
	if limit, ok := riskProfile.Limits[RiskLimitTypePositionSize]; ok && limit.Enabled {
		if math.Abs(float64(newQuantity)) > limit.Value {
			return NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				ErrCodePositionLimitExceeded,
				fmt.Sprintf("Position size %d exceeds limit %.0f", newQuantity, limit.Value),
				nil,
				"EnhancedRiskManager",
			).WithOrderID(order.ID).WithPortfolioID(order.PortfolioID).WithStrategyID(order.StrategyID)
		}
	}

	return nil
}

// validateDrawdown validates the drawdown
func (r *EnhancedRiskManager) validateDrawdown(order Order, portfolio Portfolio, riskProfile RiskProfile) error {
	// Check if drawdown limit is enabled
	if limit, ok := riskProfile.Limits[RiskLimitTypeDrawdown]; ok && limit.Enabled {
		// In a real implementation, this would calculate the potential drawdown
		// based on historical data and risk models
		
		// For now, we'll use a simplified approach
		// This is a placeholder for actual drawdown calculation
		potentialDrawdown := 0.0
		
		if potentialDrawdown > limit.Value {
			return NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				ErrCodeInvalidOrder,
				fmt.Sprintf("Potential drawdown %.2f%% exceeds limit %.2f%%", potentialDrawdown*100, limit.Value*100),
				nil,
				"EnhancedRiskManager",
			).WithOrderID(order.ID).WithPortfolioID(order.PortfolioID).WithStrategyID(order.StrategyID)
		}
	}

	return nil
}

// validateLeverage validates the leverage
func (r *EnhancedRiskManager) validateLeverage(order Order, portfolio Portfolio, riskProfile RiskProfile) error {
	// Check if leverage limit is enabled
	if limit, ok := riskProfile.Limits[RiskLimitTypeLeverage]; ok && limit.Enabled {
		// In a real implementation, this would calculate the leverage
		// based on margin requirements and account equity
		
		// For now, we'll use a simplified approach
		// This is a placeholder for actual leverage calculation
		leverage := 1.0
		
		if leverage > limit.Value {
			return NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				ErrCodeInvalidOrder,
				fmt.Sprintf("Leverage %.2fx exceeds limit %.2fx", leverage, limit.Value),
				nil,
				"EnhancedRiskManager",
			).WithOrderID(order.ID).WithPortfolioID(order.PortfolioID).WithStrategyID(order.StrategyID)
		}
	}

	return nil
}

// validateConcentration validates the concentration
func (r *EnhancedRiskManager) validateConcentration(order Order, portfolio Portfolio, riskProfile RiskProfile) error {
	// Check if concentration limit is enabled
	if limit, ok := riskProfile.Limits[RiskLimitTypeConcentration]; ok && limit.Enabled {
		// In a real implementation, this would calculate the concentration
		// based on portfolio value and position value
		
		// For now, we'll use a simplified approach
		// This is a placeholder for actual concentration calculation
		concentration := 0.0
		
		if concentration > limit.Value {
			return NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				ErrCodeInvalidOrder,
				fmt.Sprintf("Concentration %.2f%% exceeds limit %.2f%%", concentration*100, limit.Value*100),
				nil,
				"EnhancedRiskManager",
			).WithOrderID(order.ID).WithPortfolioID(order.PortfolioID).WithStrategyID(order.StrategyID)
		}
	}

	return nil
}

// validateOrderRate validates the order rate
func (r *EnhancedRiskManager) validateOrderRate(order Order, portfolio Portfolio, riskProfile RiskProfile) error {
	// Check if order rate limit is enabled
	if limit, ok := riskProfile.Limits[RiskLimitTypeOrderRate]; ok && limit.Enabled {
		// Get order history
		orderHistory, err := r.GetOrderHistory(portfolio.ID)
		if err != nil {
			// No order history, first order is always valid
			return nil
		}
		
		// Count orders in the last minute
		now := time.Now()
		orderCount := 0
		for _, historyOrder := range orderHistory {
			if now.Sub(historyOrder.CreatedAt) <= time.Minute {
				orderCount++
			}
		}
		
		if float64(orderCount) >= limit.Value {
			return NewExecutionError(
				ErrorTypeValidation,
				ErrorSeverityError,
				Er
(Content truncated due to size limit. Use line ranges to read in chunks)