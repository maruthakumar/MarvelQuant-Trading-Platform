package services

import (
	"errors"
	"time"
	"github.com/google/uuid"
	"trading_platform/backend/internal/models"
)

// MarketSimulationService handles operations related to market simulation
type MarketSimulationService struct {
	// Dependencies would be injected here in a real implementation
	// For example: database connection, market data service, etc.
}

// NewMarketSimulationService creates a new instance of MarketSimulationService
func NewMarketSimulationService() *MarketSimulationService {
	return &MarketSimulationService{}
}

// GetCurrentMarketPrice retrieves the current market price for a symbol
func (s *MarketSimulationService) GetCurrentMarketPrice(symbol string) (*models.MarketDataSnapshot, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	
	// In a real implementation, we would retrieve the current market price
	// from a market data provider or from a cache
	
	// For now, return a mock market data snapshot
	var price float64
	if symbol == "AAPL" {
		price = 150.25
	} else if symbol == "MSFT" {
		price = 280.50
	} else if symbol == "GOOGL" {
		price = 2100.75
	} else {
		price = 100.00
	}
	
	return &models.MarketDataSnapshot{
		ID:          uuid.New().String(),
		Symbol:      symbol,
		Timestamp:   time.Now(),
		Open:        price - 1.0,
		High:        price + 0.5,
		Low:         price - 1.5,
		Close:       price,
		Volume:      1000000,
		Bid:         price - 0.05,
		Ask:         price + 0.05,
		BidSize:     500,
		AskSize:     700,
		Timeframe:   "1m",
		Source:      "SIMULATION",
		IsSimulated: true,
	}, nil
}

// GetHistoricalMarketData retrieves historical market data for a symbol
func (s *MarketSimulationService) GetHistoricalMarketData(symbol string, startDate, endDate time.Time, timeframe string) ([]models.MarketDataSnapshot, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	
	if timeframe == "" {
		return nil, errors.New("timeframe is required")
	}
	
	// In a real implementation, we would retrieve historical market data
	// from a market data provider or from a database
	
	// For now, return mock historical market data
	var snapshots []models.MarketDataSnapshot
	
	// Generate mock data points
	currentTime := startDate
	var basePrice float64
	if symbol == "AAPL" {
		basePrice = 145.00
	} else if symbol == "MSFT" {
		basePrice = 275.00
	} else if symbol == "GOOGL" {
		basePrice = 2050.00
	} else {
		basePrice = 100.00
	}
	
	// Determine time increment based on timeframe
	var increment time.Duration
	switch timeframe {
	case "1m":
		increment = 1 * time.Minute
	case "5m":
		increment = 5 * time.Minute
	case "15m":
		increment = 15 * time.Minute
	case "1h":
		increment = 1 * time.Hour
	case "1d":
		increment = 24 * time.Hour
	default:
		increment = 1 * time.Hour
	}
	
	// Generate data points
	for currentTime.Before(endDate) {
		// Simulate price movement
		priceChange := (float64(currentTime.Nanosecond() % 100) / 1000.0) - 0.05 // Random price change between -0.05 and 0.05
		price := basePrice + priceChange
		
		// Create snapshot
		snapshot := models.MarketDataSnapshot{
			ID:          uuid.New().String(),
			Symbol:      symbol,
			Timestamp:   currentTime,
			Open:        price - 0.2,
			High:        price + 0.3,
			Low:         price - 0.4,
			Close:       price,
			Volume:      500000 + int64(currentTime.Nanosecond()%500000),
			Bid:         price - 0.05,
			Ask:         price + 0.05,
			BidSize:     400 + (currentTime.Nanosecond() % 200),
			AskSize:     600 + (currentTime.Nanosecond() % 200),
			Timeframe:   timeframe,
			Source:      "HISTORICAL",
			IsSimulated: true,
		}
		
		snapshots = append(snapshots, snapshot)
		
		// Increment time
		currentTime = currentTime.Add(increment)
		
		// Slightly adjust base price for next iteration to simulate a trend
		basePrice += priceChange * 0.1
	}
	
	return snapshots, nil
}

// SimulateOrderExecution simulates the execution of an order
func (s *MarketSimulationService) SimulateOrderExecution(order *models.SimulationOrder, marketSettings *models.MarketSettings) error {
	if order == nil {
		return errors.New("order is required")
	}
	
	if marketSettings == nil {
		return errors.New("market settings are required")
	}
	
	// Get current market price
	marketData, err := s.GetCurrentMarketPrice(order.Symbol)
	if err != nil {
		return err
	}
	
	// Determine execution price based on order type
	var executionPrice float64
	if order.OrderType == "MARKET" {
		// For market orders, use current market price
		if order.Side == "BUY" {
			executionPrice = marketData.Ask // Buy at ask price
		} else {
			executionPrice = marketData.Bid // Sell at bid price
		}
	} else if order.OrderType == "LIMIT" {
		// For limit orders, check if the limit price is satisfied
		if order.Side == "BUY" {
			if order.Price >= marketData.Ask {
				executionPrice = marketData.Ask // Buy at ask price, but no higher than limit price
				if executionPrice > order.Price {
					executionPrice = order.Price
				}
			} else {
				// Limit price not satisfied
				return nil
			}
		} else {
			if order.Price <= marketData.Bid {
				executionPrice = marketData.Bid // Sell at bid price, but no lower than limit price
				if executionPrice < order.Price {
					executionPrice = order.Price
				}
			} else {
				// Limit price not satisfied
				return nil
			}
		}
	} else if order.OrderType == "STOP" {
		// For stop orders, check if the stop price is triggered
		if order.Side == "BUY" {
			if marketData.Ask >= order.StopPrice {
				executionPrice = marketData.Ask // Buy at ask price once triggered
			} else {
				// Stop price not triggered
				return nil
			}
		} else {
			if marketData.Bid <= order.StopPrice {
				executionPrice = marketData.Bid // Sell at bid price once triggered
			} else {
				// Stop price not triggered
				return nil
			}
		}
	} else if order.OrderType == "STOP_LIMIT" {
		// For stop-limit orders, check if the stop price is triggered
		// and then apply limit price logic
		if order.Side == "BUY" {
			if marketData.Ask >= order.StopPrice {
				if order.Price >= marketData.Ask {
					executionPrice = marketData.Ask // Buy at ask price, but no higher than limit price
					if executionPrice > order.Price {
						executionPrice = order.Price
					}
				} else {
					// Limit price not satisfied
					return nil
				}
			} else {
				// Stop price not triggered
				return nil
			}
		} else {
			if marketData.Bid <= order.StopPrice {
				if order.Price <= marketData.Bid {
					executionPrice = marketData.Bid // Sell at bid price, but no lower than limit price
					if executionPrice < order.Price {
						executionPrice = order.Price
					}
				} else {
					// Limit price not satisfied
					return nil
				}
			} else {
				// Stop price not triggered
				return nil
			}
		}
	} else {
		return errors.New("unsupported order type")
	}
	
	// Apply slippage
	var slippageAmount float64
	if marketSettings.SlippageModel == "FIXED" {
		slippageAmount = marketSettings.SlippageValue
	} else if marketSettings.SlippageModel == "PERCENTAGE" {
		slippageAmount = executionPrice * marketSettings.SlippageValue
	} else if marketSettings.SlippageModel == "VARIABLE" {
		// Simulate variable slippage based on order size and market conditions
		volumeFactor := float64(order.Quantity) / float64(marketData.Volume) * 10000 // Scale factor
		if volumeFactor > 1.0 {
			volumeFactor = 1.0
		}
		slippageAmount = executionPrice * marketSettings.SlippageValue * (1.0 + volumeFactor)
	} else {
		slippageAmount = 0
	}
	
	// Apply slippage to execution price
	if order.Side == "BUY" {
		executionPrice += slippageAmount
	} else {
		executionPrice -= slippageAmount
	}
	
	// Calculate commission
	var commissionAmount float64
	if marketSettings.CommissionModel == "FIXED" {
		commissionAmount = marketSettings.CommissionValue
	} else if marketSettings.CommissionModel == "PERCENTAGE" {
		commissionAmount = executionPrice * float64(order.Quantity) * marketSettings.CommissionValue
	} else if marketSettings.CommissionModel == "TIERED" {
		// Simulate tiered commission based on order size
		if order.Quantity <= 100 {
			commissionAmount = executionPrice * float64(order.Quantity) * 0.002 // 0.2% for small orders
		} else if order.Quantity <= 1000 {
			commissionAmount = executionPrice * float64(order.Quantity) * 0.001 // 0.1% for medium orders
		} else {
			commissionAmount = executionPrice * float64(order.Quantity) * 0.0005 // 0.05% for large orders
		}
	} else {
		commissionAmount = 0
	}
	
	// Simulate latency
	var latencyMs int
	if marketSettings.LatencyModel == "FIXED" {
		latencyMs = marketSettings.LatencyValue
	} else if marketSettings.LatencyModel == "VARIABLE" {
		// Simulate variable latency based on market conditions
		baseLatency := marketSettings.LatencyValue
		variableFactor := time.Now().Nanosecond() % 100 // Random factor between 0 and 99
		latencyMs = baseLatency + (variableFactor / 2)
	} else if marketSettings.LatencyModel == "REALISTIC" {
		// Simulate realistic latency with occasional spikes
		baseLatency := marketSettings.LatencyValue
		if time.Now().Nanosecond()%100 < 5 { // 5% chance of latency spike
			latencyMs = baseLatency * 5
		} else {
			variableFactor := time.Now().Nanosecond() % 50 // Random factor between 0 and 49
			latencyMs = baseLatency + variableFactor
		}
	} else {
		latencyMs = 0
	}
	
	// Update order with execution details
	order.Status = "FILLED"
	order.FilledQty = order.Quantity
	order.AvgFillPrice = executionPrice
	order.UpdatedAt = time.Now()
	order.SimulatedFillPrice = executionPrice
	order.SimulatedFillTime = time.Now()
	order.SlippageAmount = slippageAmount
	order.LatencyMs = latencyMs
	order.CommissionAmount = commissionAmount
	
	return nil
}

// SimulateMarketMovement simulates market movement for a symbol
func (s *MarketSimulationService) SimulateMarketMovement(symbol string, timeframe string, duration time.Duration) ([]models.MarketDataSnapshot, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	
	if timeframe == "" {
		return nil, errors.New("timeframe is required")
	}
	
	// Get current market data
	currentData, err := s.GetCurrentMarketPrice(symbol)
	if err != nil {
		return nil, err
	}
	
	// Determine time increment based on timeframe
	var increment time.Duration
	switch timeframe {
	case "1m":
		increment = 1 * time.Minute
	case "5m":
		increment = 5 * time.Minute
	case "15m":
		increment = 15 * time.Minute
	case "1h":
		increment = 1 * time.Hour
	case "1d":
		increment = 24 * time.Hour
	default:
		increment = 1 * time.Hour
	}
	
	// Calculate number of data points
	numPoints := int(duration / increment)
	if numPoints <= 0 {
		numPoints = 1
	}
	
	// Generate simulated market data
	var snapshots []models.MarketDataSnapshot
	currentPrice := currentData.Close
	currentTime := time.Now()
	
	for i := 0; i < numPoints; i++ {
		// Simulate price movement
		priceChange := (float64(time.Now().Nanosecond()%200) / 1000.0) - 0.1 // Random price change between -0.1 and 0.1
		currentPrice += priceChange
		
		// Ensure price doesn't go negative
		if currentPrice <= 0 {
			currentPrice = 0.01
		}
		
		// Create snapshot
		snapshot := models.MarketDataSnapshot{
			ID:          uuid.New().String(),
			Symbol:      symbol,
			Timestamp:   currentTime,
			Open:        currentPrice - (priceChange * 0.5),
			High:        currentPrice + (priceChange * 0.2),
			Low:         currentPrice - (priceChange * 0.7),
			Close:       currentPrice,
			Volume:      500000 + int64(time.Now().Nanosecond()%500000),
			Bid:         currentPrice - 0.05,
			Ask:         currentPrice + 0.05,
			BidSize:     400 + (time.Now().Nanosecond() % 200),
			AskSize:     600 + (time.Now().Nanosecond() % 200),
			Timeframe:   timeframe,
			Source:      "SIMULATION",
			IsSimulated: true,
		}
		
		snapshots = append(snapshots, snapshot)
		
		// Increment time
		currentTime = currentTime.Add(increment)
	}
	
	return snapshots, nil
}

// SimulateMarketEvent simulates a market event (e.g., earnings announcement, economic data release)
func (s *MarketSimulationService) SimulateMarketEvent(symbol string, eventType string, impactMagnitude float64) (*models.MarketDataSnapshot, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	
	if eventType == "" {
		return nil, errors.New("event type is required")
	}
	
	// Get current market data
	currentData, err := s.GetCurrentMarketPrice(symbol)
	if err != nil {
		return nil, err
	}
	
	// Determine price impact based on event type and magnitude
	var priceImpact float64
	if eventType == "EARNINGS_BEAT" {
		priceImpact = currentData.Close * impactMagnitude * 0.05 // Positive impact
	} else if eventType == "EARNINGS_MISS" {
		priceImpact = -currentData.Close * impactMagnitude * 0.05 // Negative impact
	} else if eventType == "ECONOMIC_DATA_POSITIVE" {
		priceImpact = currentData.Close * impactMagnitude * 0.02 // Positive impact
	} else if eventType == "ECONOMIC_DATA_NEGATIVE" {
		priceImpact = -currentData.Close * impactMagnitude * 0.02 // Negative impact
	} else if eventType == "MERGER_ANNOUNCEMENT" {
		priceImpact = currentData.Close * impactMagnitude * 0.1 // Positive impact
	} else if eventType == "REGULATORY_ISSUE" {
		priceImpact = -currentData.Close * impactMagnitude * 0.08 // Negative impact
	} else {
		priceImpact = 0
	}
	
	// Apply price impact
	newPrice := currentData.Close + priceImpact
	
	// Ensure price doesn't go negative
	if newPrice <= 0 {
		newPrice = 0.01
	}
	
	// Create new market data snapshot
	snapshot := models.MarketDataSnapshot{
		ID:          uuid.New().String(),
		Symbol:      symbol,
		Timestamp:   time.Now(),
		Open:        currentData.Close, // Open at previous close
		High:        math.Max(currentData.Close, newPrice),
		Low:         math.Min(currentData.Close, newPrice),
		Close:       newPrice,
		Volume:      currentData.Volume * 3, // Increased volume due to event
		Bid:         newPrice - 0.1,
		Ask:         newPrice + 0.1,
		BidSize:     currentData.BidSize * 2, // Increased liquidity
		AskSize:     currentData.AskSize * 2, // Increased liquidity
		Timeframe:   "1m",
		Source:      "EVENT_SIMULATION",
		IsSimulated: true,
	}
	
	return &snapshot, nil
}

// GetMarketDepth retrieves the market depth (order book) for a symbol
func (s *MarketSimulationService) GetMarketDepth(symbol string, levels int) (map[string]interface{}, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	
	if levels <= 0 {
		levels = 5 // Default to 5 levels
	}
	
	// Get current market price
	marketData, err := s.GetCurrentMarketPrice(symbol)
	if err != nil {
		return nil, err
	}
	
	// Generate simulated order book
	var bids []map[string]interface{}
	var asks []map[string]interface{}
	
	// Generate bid levels
	bidPrice := marketData.Bid
	bidSize := marketData.BidSize
	for i := 0; i < levels; i++ {
		bids = append(bids, map[string]interface{}{
			"price": bidPrice,
			"size":  bidSize,
		})
		
		// Decrease price and size for next level
		bidPrice -= 0.01 + (float64(i) * 0.01)
		bidSize = int(float64(bidSize) * 0.8)
	}
	
	// Generate ask levels
	askPrice := marketData.Ask
	askSize := marketData.AskSize
	for i := 0; i < levels; i++ {
		asks = append(asks, map[string]interface{}{
			"price": askPrice,
			"size":  askSize,
		})
		
		// Increase price and size for next level
		askPrice += 0.01 + (float64(i) * 0.01)
		askSize = int(float64(askSize) * 0.8)
	}
	
	return map[string]interface{}{
		"symbol":    symbol,
		"timestamp": time.Now(),
		"bids":      bids,
		"asks":      asks,
		"spread":    marketData.Ask - marketData.Bid,
	}, nil
}

// SimulateMarketCondition simulates a specific market condition
func (s *MarketSimulationService) SimulateMarketCondition(symbol string, condition string, duration time.Duration) ([]models.MarketDataSnapshot, error) {
	if symbol == "" {
		return nil, errors.New("symbol is required")
	}
	
	if condition == "" {
		return nil, errors.New("condition is required")
	}
	
	// Get current market data
	currentData, err := s.GetCurrentMarketPrice(symbol)
	if err != nil {
		return nil, err
	}
	
	// Determine simulation parameters based on condition
	var volatilityFactor float64
	var trendFactor float64
	var volumeFactor float64
	
	if condition == "HIGH_VOLATILITY" {
		volatilityFactor = 3.0
		trendFactor = 0.0
		volumeFactor = 2.0
	} else if condition == "LOW_VOLATILITY" {
		volatilityFactor = 0.3
		trendFactor = 0.0
		volumeFactor = 0.7
	} else if condition == "UPTREND" {
		volatilityFactor = 1.0
		trendFactor = 0.002 // 0.2% upward trend per minute
		volumeFactor = 1.2
	} else if condition == "DOWNTREND" {
		volatilityFactor = 1.0
		trendFactor = -0.002 // 0.2% downward trend per minute
		volumeFactor = 1.2
	} else if condition == "SIDEWAYS" {
		volatilityFactor = 0.5
		trendFactor = 0.0
		volumeFactor = 0.8
	} else if condition == "MARKET_CRASH" {
		volatilityFactor = 5.0
		trendFactor = -0.01 // 1% downward trend per minute
		volumeFactor = 3.0
	} else if condition == "MARKET_RALLY" {
		volatilityFactor = 3.0
		trendFactor = 0.008 // 0.8% upward trend per minute
		volumeFactor = 2.5
	} else {
		return nil, errors.New("unsupported market condition")
	}
	
	// Simulate market data
	var snapshots []models.MarketDataSnapshot
	currentPrice := currentData.Close
	currentTime := time.Now()
	
	// Determine number of data points (1 per minute)
	numPoints := int(duration.Minutes())
	if numPoints <= 0 {
		numPoints = 1
	}
	
	for i := 0; i < numPoints; i++ {
		// Apply trend
		currentPrice *= (1.0 + trendFactor)
		
		// Apply volatility
		volatility := (float64(time.Now().Nanosecond()%200) / 1000.0) - 0.1 // Random between -0.1 and 0.1
		volatility *= volatilityFactor
		currentPrice *= (1.0 + volatility)
		
		// Ensure price doesn't go negative
		if currentPrice <= 0 {
			currentPrice = 0.01
		}
		
		// Calculate volume
		volume := int64(float64(currentData.Volume) * volumeFactor * (1.0 + (float64(time.Now().Nanosecond()%100) / 100.0)))
		
		// Create snapshot
		snapshot := models.MarketDataSnapshot{
			ID:          uuid.New().String(),
			Symbol:      symbol,
			Timestamp:   currentTime,
			Open:        currentPrice * (1.0 - (volatility * 0.2)),
			High:        currentPrice * (1.0 + (volatility * 0.5)),
			Low:         currentPrice * (1.0 - (volatility * 0.5)),
			Close:       currentPrice,
			Volume:      volume,
			Bid:         currentPrice * 0.998,
			Ask:         currentPrice * 1.002,
			BidSize:     int(float64(currentData.BidSize) * volumeFactor),
			AskSize:     int(float64(currentData.AskSize) * volumeFactor),
			Timeframe:   "1m",
			Source:      "CONDITION_SIMULATION",
			IsSimulated: true,
		}
		
		snapshots = append(snapshots, snapshot)
		
		// Increment time
		currentTime = currentTime.Add(1 * time.Minute)
	}
	
	return snapshots, nil
}
