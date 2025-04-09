package position

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// PositionService defines the interface for position-related operations
type PositionService interface {
	CreatePositionFromOrder(order *models.Order) (*models.Position, error)
	GetPositionByID(id string) (*models.Position, error)
	GetPositions(filter models.PositionFilter, page, limit int) ([]models.Position, int, error)
	UpdatePosition(position *models.Position) (*models.Position, error)
	ClosePosition(id string, exitPrice float64, exitQuantity int) (*models.Position, error)
	CalculatePnL(position *models.Position) (float64, error)
	CalculateGreeks(position *models.Position) (*models.Greeks, error)
	CalculateExposure(positions []models.Position) (float64, error)
	AggregatePositions(positions []models.Position, groupBy string) (map[string]models.AggregatedPosition, error)
}

// PositionServiceImpl implements the PositionService interface
type PositionServiceImpl struct {
	positionRepo repositories.PositionRepository
	orderRepo    repositories.OrderRepository
}

// NewPositionService creates a new PositionService
func NewPositionService(positionRepo repositories.PositionRepository, orderRepo repositories.OrderRepository) PositionService {
	return &PositionServiceImpl{
		positionRepo: positionRepo,
		orderRepo:    orderRepo,
	}
}

// CreatePositionFromOrder creates a new position from an executed order
func (s *PositionServiceImpl) CreatePositionFromOrder(order *models.Order) (*models.Position, error) {
	// Validate the order
	if order == nil {
		return nil, errors.New("order cannot be nil")
	}

	// Check if order is executed
	if order.Status != models.OrderStatusExecuted && order.Status != models.OrderStatusPartial {
		return nil, errors.New("only executed or partially executed orders can create positions")
	}

	// Check if position already exists for this order
	existingPositions, _, err := s.positionRepo.GetAll(models.PositionFilter{OrderID: order.ID}, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(existingPositions) > 0 {
		return nil, errors.New("position already exists for this order")
	}

	// Create a new position
	position := &models.Position{
		UserID:         order.UserID,
		OrderID:        order.ID,
		Symbol:         order.Symbol,
		Exchange:       order.Exchange,
		Direction:      convertOrderDirectionToPositionDirection(order.Direction),
		EntryPrice:     order.ExecutionPrice,
		Quantity:       order.FilledQuantity,
		Status:         models.PositionStatusOpen,
		ProductType:    order.ProductType,
		InstrumentType: order.InstrumentType,
		OptionType:     order.OptionType,
		StrikePrice:    order.StrikePrice,
		Expiry:         order.Expiry,
		PortfolioID:    order.PortfolioID,
		StrategyID:     order.StrategyID,
		LegID:          order.LegID,
		Tags:           order.Tags,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Calculate initial P&L and Greeks
	initialPnL, _ := s.CalculatePnL(position)
	position.UnrealizedPnL = initialPnL

	initialGreeks, _ := s.CalculateGreeks(position)
	position.Greeks = *initialGreeks

	// Validate the position
	if err := position.Validate(); err != nil {
		return nil, err
	}

	// Create the position
	createdPosition, err := s.positionRepo.Create(position)
	if err != nil {
		return nil, err
	}

	return createdPosition, nil
}

// GetPositionByID retrieves a position by ID
func (s *PositionServiceImpl) GetPositionByID(id string) (*models.Position, error) {
	if id == "" {
		return nil, errors.New("position ID is required")
	}

	position, err := s.positionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update P&L and Greeks
	currentPnL, _ := s.CalculatePnL(position)
	position.UnrealizedPnL = currentPnL

	currentGreeks, _ := s.CalculateGreeks(position)
	position.Greeks = *currentGreeks

	return position, nil
}

// GetPositions retrieves positions with filtering and pagination
func (s *PositionServiceImpl) GetPositions(filter models.PositionFilter, page, limit int) ([]models.Position, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100 // Maximum limit to prevent excessive queries
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get positions with pagination
	positions, total, err := s.positionRepo.GetAll(filter, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Update P&L and Greeks for each position
	for i := range positions {
		currentPnL, _ := s.CalculatePnL(&positions[i])
		positions[i].UnrealizedPnL = currentPnL

		currentGreeks, _ := s.CalculateGreeks(&positions[i])
		positions[i].Greeks = *currentGreeks
	}

	return positions, total, nil
}

// UpdatePosition updates an existing position
func (s *PositionServiceImpl) UpdatePosition(position *models.Position) (*models.Position, error) {
	// Validate the position
	if err := position.Validate(); err != nil {
		return nil, err
	}

	// Check if position exists
	existingPosition, err := s.positionRepo.GetByID(position.ID)
	if err != nil {
		return nil, errors.New("position not found")
	}

	// Check if position can be updated
	if existingPosition.Status == models.PositionStatusClosed {
		return nil, errors.New("closed positions cannot be updated")
	}

	// Preserve certain fields from the existing position
	position.CreatedAt = existingPosition.CreatedAt
	position.UpdatedAt = time.Now()

	// Calculate P&L and Greeks
	currentPnL, _ := s.CalculatePnL(position)
	position.UnrealizedPnL = currentPnL

	currentGreeks, _ := s.CalculateGreeks(position)
	position.Greeks = *currentGreeks

	// Update the position
	updatedPosition, err := s.positionRepo.Update(position)
	if err != nil {
		return nil, err
	}

	return updatedPosition, nil
}

// ClosePosition closes an existing position
func (s *PositionServiceImpl) ClosePosition(id string, exitPrice float64, exitQuantity int) (*models.Position, error) {
	if id == "" {
		return nil, errors.New("position ID is required")
	}

	// Check if position exists
	existingPosition, err := s.positionRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("position not found")
	}

	// Check if position can be closed
	if existingPosition.Status == models.PositionStatusClosed {
		return nil, errors.New("position is already closed")
	}

	// Validate exit parameters
	if exitPrice <= 0 {
		return nil, errors.New("exit price must be greater than zero")
	}
	if exitQuantity <= 0 {
		return nil, errors.New("exit quantity must be greater than zero")
	}
	if exitQuantity > existingPosition.Quantity {
		return nil, errors.New("exit quantity cannot exceed position quantity")
	}

	// Update position
	existingPosition.ExitPrice = exitPrice
	existingPosition.ExitQuantity = exitQuantity
	existingPosition.UpdatedAt = time.Now()

	// Calculate realized P&L
	var realizedPnL float64
	if existingPosition.Direction == models.PositionDirectionLong {
		realizedPnL = (exitPrice - existingPosition.EntryPrice) * float64(exitQuantity)
	} else {
		realizedPnL = (existingPosition.EntryPrice - exitPrice) * float64(exitQuantity)
	}
	existingPosition.RealizedPnL = realizedPnL

	// Update position status
	if exitQuantity == existingPosition.Quantity {
		existingPosition.Status = models.PositionStatusClosed
	} else {
		existingPosition.Status = models.PositionStatusPartial
	}

	// Calculate unrealized P&L for remaining quantity
	if existingPosition.Status == models.PositionStatusPartial {
		remainingQuantity := existingPosition.Quantity - exitQuantity
		currentPnL, _ := s.CalculatePnL(existingPosition)
		existingPosition.UnrealizedPnL = currentPnL * float64(remainingQuantity) / float64(existingPosition.Quantity)
	} else {
		existingPosition.UnrealizedPnL = 0
	}

	// Update Greeks
	currentGreeks, _ := s.CalculateGreeks(existingPosition)
	existingPosition.Greeks = *currentGreeks

	// Update the position
	updatedPosition, err := s.positionRepo.Update(existingPosition)
	if err != nil {
		return nil, err
	}

	return updatedPosition, nil
}

// CalculatePnL calculates the P&L for a position
func (s *PositionServiceImpl) CalculatePnL(position *models.Position) (float64, error) {
	if position == nil {
		return 0, errors.New("position cannot be nil")
	}

	// For closed positions, return the realized P&L
	if position.Status == models.PositionStatusClosed {
		return position.RealizedPnL, nil
	}

	// For open or partially closed positions, calculate the unrealized P&L
	// In a real implementation, this would fetch the current market price
	// For now, we'll use a simple calculation based on a mock current price
	// that's 1% higher than the entry price (this is just for demonstration)
	mockCurrentPrice := position.EntryPrice * 1.01

	var unrealizedPnL float64
	remainingQuantity := position.Quantity - position.ExitQuantity
	if position.Direction == models.PositionDirectionLong {
		unrealizedPnL = (mockCurrentPrice - position.EntryPrice) * float64(remainingQuantity)
	} else {
		unrealizedPnL = (position.EntryPrice - mockCurrentPrice) * float64(remainingQuantity)
	}

	return unrealizedPnL, nil
}

// CalculateGreeks calculates the Greeks for a position
func (s *PositionServiceImpl) CalculateGreeks(position *models.Position) (*models.Greeks, error) {
	if position == nil {
		return nil, errors.New("position cannot be nil")
	}

	// In a real implementation, this would use option pricing models
	// For now, we'll return mock values based on the position
	greeks := &models.Greeks{
		Delta: 0.5,
		Gamma: 0.05,
		Theta: -0.1,
		Vega:  0.2,
	}

	// Adjust based on position type and direction
	if position.InstrumentType == models.InstrumentTypeOption {
		if position.OptionType == models.OptionTypeCall {
			if position.Direction == models.PositionDirectionLong {
				greeks.Delta = 0.6
			} else {
				greeks.Delta = -0.6
			}
		} else {
			if position.Direction == models.PositionDirectionLong {
				greeks.Delta = -0.4
			} else {
				greeks.Delta = 0.4
			}
		}
	} else {
		// For non-options, set appropriate values
		greeks.Delta = position.Direction == models.PositionDirectionLong ? 1.0 : -1.0
		greeks.Gamma = 0
		greeks.Theta = 0
		greeks.Vega = 0
	}

	// Scale by position size
	remainingQuantity := position.Quantity - position.ExitQuantity
	greeks.Delta *= float64(remainingQuantity)
	greeks.Gamma *= float64(remainingQuantity)
	greeks.Theta *= float64(remainingQuantity)
	greeks.Vega *= float64(remainingQuantity)

	return greeks, nil
}

// CalculateExposure calculates the total exposure for a list of positions
func (s *PositionServiceImpl) CalculateExposure(positions []models.Position) (float64, error) {
	if len(positions) == 0 {
		return 0, nil
	}

	var totalExposure float64
	for _, position := range positions {
		// Skip closed positions
		if position.Status == models.PositionStatusClosed {
			continue
		}

		// Calculate exposure based on position value
		remainingQuantity := position.Quantity - position.ExitQuantity
		positionValue := position.EntryPrice * float64(remainingQuantity)
		
		// For short positions, exposure is still positive
		totalExposure += positionValue
	}

	return totalExposure, nil
}

// AggregatePositions aggregates positions by the specified grouping
func (s *PositionServiceImpl) AggregatePositions(positions []models.Position, groupBy string) (map[string]models.AggregatedPosition, error) {
	if len(positions) == 0 {
		return make(map[string]models.AggregatedPosition), nil
	}

	// Validate groupBy parameter
	validGroupings := map[string]bool{
		"symbol":         true,
		"instrumentType": true,
		"productType":    true,
		"strategy":       true,
		"portfolio":      true,
	}
	if !validGroupings[groupBy] {
		return nil, errors.New("invalid groupBy parameter")
	}

	// Aggregate positions
	aggregated := make(map[string]models.AggregatedPosition)
	for _, position := range positions {
		// Skip closed positions
		if position.Status == models.PositionStatusClosed {
			continue
		}

		// Determine the key based on groupBy
		var key string
		switch groupBy {
		case "symbol":
			key = position.Symbol
		case "instrumentType":
			key = string(position.InstrumentType)
		case "productType":
			key = string(position.ProductType)
		case "strategy":
			key = position.StrategyID
		case "portfolio":
			key = position.PortfolioID
		}

		// If key is empty, use "unknown"
		if key == "" {
			key = "unknown"
		}

		// Get or create the aggregated position
		agg, exists := aggregated[key]
		if !exists {
			agg = models.AggregatedPosition{
				Key:           key,
				GroupBy:       groupBy,
				TotalQuantity: 0,
				NetQuantity:   0,
				TotalValue:    0,
				NetValue:      0,
				PnL:           0,
				Greeks: models.Greeks{
					Delta: 0,
					Gamma: 0,
					Theta: 0,
					Vega:  0,
				},
				PositionCount: 0,
			}
		}

		// Calculate remaining quantity
		remainingQuantity := position.Quantity - position.ExitQuantity
		
		// Update aggregated position
		agg.PositionCount++
		agg.TotalQuantity += remainingQuantity
		
		// For net calculations, consider direction
		if position.Direction == models.PositionDirectionLong {
			agg.NetQuantity += remainingQuantity
		} else {
			agg.NetQuantity -= remainingQuantity
		}
		
		// Calculate position value
		positionValue := position.EntryPrice * float64(remainingQuantity)
		agg.TotalValue += positionValue
		
		// For net value, consider direction
		if position.Direction == models.PositionDirectionLong {
			agg.NetValue += positionValue
		} else {
			agg.NetValue -= positionValue
		}
		
		// Add P&L
		agg.PnL += position.UnrealizedPnL + position.RealizedPnL
		
		// Add Greeks
		agg.Greeks.Delta += position.Greeks.Delta
		agg.Greeks.Gamma += position.Greeks.Gamma
		agg.Greeks.Theta += position.Greeks.Theta
		agg.Greeks.Vega += position.Greeks.Vega
		
		// Update the map
		aggregated[key] = agg
	}

	return aggregated, nil
}

// Helper function to convert order direction to position direction
func convertOrderDirectionToPositionDirection(direction models.OrderDirection) models.PositionDirection {
	if direction == models.OrderDirectionBuy {
		return models.PositionDirectionLong
	}
	return models.PositionDirectionShort
}
