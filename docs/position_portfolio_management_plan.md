# Position and Portfolio Management System Implementation Plan

## Overview

This document outlines the implementation plan for the Position and Portfolio Management System, a core component of the Trading Platform. This system is responsible for tracking positions, managing portfolios, calculating performance metrics, and providing real-time portfolio analytics.

## Architecture

The Position and Portfolio Management System follows a layered architecture:

1. **Position Tracker**: Tracks individual positions across all instruments and brokers
2. **Portfolio Manager**: Manages portfolios consisting of multiple positions
3. **Performance Calculator**: Calculates performance metrics for positions and portfolios
4. **Portfolio Analytics Engine**: Provides advanced analytics for portfolios
5. **Reporting System**: Generates reports for positions and portfolios

## Implementation Components

### 1. Position Tracker

```go
// position_tracker.go
package portfolio

import (
    "sync"
    "time"
    "github.com/trading-platform/backend/internal/broker/common"
    "github.com/trading-platform/backend/internal/api"
    "github.com/trading-platform/backend/internal/database"
)

// Position represents a trading position
type Position struct {
    UserID            string
    BrokerID          string
    ExchangeSegment   string
    TradingSymbol     string
    ProductType       string
    Quantity          int
    BuyQuantity       int
    SellQuantity      int
    NetQuantity       int
    AveragePrice      float64
    LastPrice         float64
    RealizedProfit    float64
    UnrealizedProfit  float64
    Value             float64
    LastUpdated       time.Time
}

// PositionTracker tracks individual positions across all instruments and brokers
type PositionTracker struct {
    brokerManager *api.BrokerManager
    db            *database.Database
    positions     map[string]map[string]*Position // userID -> positionKey -> Position
    mu            sync.RWMutex
    logger        *Logger
}

// NewPositionTracker creates a new position tracker
func NewPositionTracker(brokerManager *api.BrokerManager, db *database.Database) *PositionTracker {
    return &PositionTracker{
        brokerManager: brokerManager,
        db:            db,
        positions:     make(map[string]map[string]*Position),
        mu:            sync.RWMutex{},
        logger:        NewLogger("position_tracker"),
    }
}

// GetPosition gets a position for a user and symbol
func (t *PositionTracker) GetPosition(userID, exchangeSegment, tradingSymbol string) (*Position, error) {
    positionKey := t.getPositionKey(exchangeSegment, tradingSymbol)
    
    // Check cache
    t.mu.RLock()
    userPositions, ok := t.positions[userID]
    if ok {
        position, ok := userPositions[positionKey]
        if ok {
            t.mu.RUnlock()
            return position, nil
        }
    }
    t.mu.RUnlock()
    
    // Get from database
    position, err := t.loadPosition(userID, exchangeSegment, tradingSymbol)
    if err != nil {
        return nil, err
    }
    
    // Cache position
    t.mu.Lock()
    if _, ok := t.positions[userID]; !ok {
        t.positions[userID] = make(map[string]*Position)
    }
    t.positions[userID][positionKey] = position
    t.mu.Unlock()
    
    return position, nil
}

// GetPositions gets all positions for a user
func (t *PositionTracker) GetPositions(userID string) ([]*Position, error) {
    // Check cache
    t.mu.RLock()
    userPositions, ok := t.positions[userID]
    t.mu.RUnlock()
    
    if ok {
        positions := make([]*Position, 0, len(userPositions))
        for _, position := range userPositions {
            positions = append(positions, position)
        }
        return positions, nil
    }
    
    // Get from database
    positions, err := t.loadPositions(userID)
    if err != nil {
        return nil, err
    }
    
    // Cache positions
    t.mu.Lock()
    t.positions[userID] = make(map[string]*Position)
    for _, position := range positions {
        positionKey := t.getPositionKey(position.ExchangeSegment, position.TradingSymbol)
        t.positions[userID][positionKey] = position
    }
    t.mu.Unlock()
    
    return positions, nil
}

// UpdatePosition updates a position based on a trade
func (t *PositionTracker) UpdatePosition(userID string, trade *common.Trade) error {
    // Get current position
    position, err := t.GetPosition(userID, trade.ExchangeSegment, trade.TradingSymbol)
    if err != nil {
        // Create new position if not found
        position = &Position{
            UserID:          userID,
            BrokerID:        trade.BrokerID,
            ExchangeSegment: trade.ExchangeSegment,
            TradingSymbol:   trade.TradingSymbol,
            ProductType:     trade.ProductType,
        }
    }
    
    // Update position based on trade
    if trade.TradeSide == "BUY" {
        position.BuyQuantity += trade.TradeQuantity
    } else {
        position.SellQuantity += trade.TradeQuantity
    }
    
    position.Quantity = position.BuyQuantity + position.SellQuantity
    position.NetQuantity = position.BuyQuantity - position.SellQuantity
    
    // Calculate average price and realized profit
    // This is a simplified calculation; in practice, you would use a more sophisticated algorithm
    if trade.TradeSide == "BUY" {
        if position.BuyQuantity > 0 {
            position.AveragePrice = ((position.AveragePrice * float64(position.BuyQuantity-trade.TradeQuantity)) + 
                                    (trade.TradePrice * float64(trade.TradeQuantity))) / 
                                    float64(position.BuyQuantity)
        }
    } else {
        if position.NetQuantity >= 0 {
            // Selling from long position
            costBasis := position.AveragePrice * float64(trade.TradeQuantity)
            saleProceeds := trade.TradePrice * float64(trade.TradeQuantity)
            position.RealizedProfit += saleProceeds - costBasis
        } else {
            // Short covering
            costBasis := position.AveragePrice * float64(trade.TradeQuantity)
            saleProceeds := trade.TradePrice * float64(trade.TradeQuantity)
            position.RealizedProfit += costBasis - saleProceeds
        }
    }
    
    // Update last price and unrealized profit
    position.LastPrice = trade.TradePrice
    if position.NetQuantity > 0 {
        // Long position
        position.UnrealizedProfit = float64(position.NetQuantity) * (position.LastPrice - position.AveragePrice)
    } else if position.NetQuantity < 0 {
        // Short position
        position.UnrealizedProfit = float64(-position.NetQuantity) * (position.AveragePrice - position.LastPrice)
    } else {
        // Flat position
        position.UnrealizedProfit = 0
    }
    
    position.Value = float64(abs(position.NetQuantity)) * position.LastPrice
    position.LastUpdated = time.Now()
    
    // Save position
    t.mu.Lock()
    if _, ok := t.positions[userID]; !ok {
        t.positions[userID] = make(map[string]*Position)
    }
    positionKey := t.getPositionKey(position.ExchangeSegment, position.TradingSymbol)
    t.positions[userID][positionKey] = position
    t.mu.Unlock()
    
    // Save to database
    if err := t.savePosition(position); err != nil {
        t.logger.Error("Failed to save position", "error", err, "user_id", userID, "symbol", position.TradingSymbol)
        return err
    }
    
    t.logger.Info("Position updated", "user_id", userID, "symbol", position.TradingSymbol, "net_quantity", position.NetQuantity)
    return nil
}

// UpdatePositionPrice updates the price of a position
func (t *PositionTracker) UpdatePositionPrice(userID, exchangeSegment, tradingSymbol string, lastPrice float64) error {
    // Get current position
    position, err := t.GetPosition(userID, exchangeSegment, tradingSymbol)
    if err != nil {
        return err
    }
    
    // Update last price and unrealized profit
    position.LastPrice = lastPrice
    if position.NetQuantity > 0 {
        // Long position
        position.UnrealizedProfit = float64(position.NetQuantity) * (position.LastPrice - position.AveragePrice)
    } else if position.NetQuantity < 0 {
        // Short position
        position.UnrealizedProfit = float64(-position.NetQuantity) * (position.AveragePrice - position.LastPrice)
    }
    
    position.Value = float64(abs(position.NetQuantity)) * position.LastPrice
    position.LastUpdated = time.Now()
    
    // Save position
    t.mu.Lock()
    if _, ok := t.positions[userID]; !ok {
        t.positions[userID] = make(map[string]*Position)
    }
    positionKey := t.getPositionKey(position.ExchangeSegment, position.TradingSymbol)
    t.positions[userID][positionKey] = position
    t.mu.Unlock()
    
    // Save to database
    if err := t.savePosition(position); err != nil {
        t.logger.Error("Failed to save position", "error", err, "user_id", userID, "symbol", position.TradingSymbol)
        return err
    }
    
    return nil
}

// getPositionKey generates a unique key for a position
func (t *PositionTracker) getPositionKey(exchangeSegment, tradingSymbol string) string {
    return exchangeSegment + ":" + tradingSymbol
}

// loadPosition loads a position from the database
func (t *PositionTracker) loadPosition(userID, exchangeSegment, tradingSymbol string) (*Position, error) {
    // Implementation
    return &Position{}, nil
}

// loadPositions loads all positions for a user from the database
func (t *PositionTracker) loadPositions(userID string) ([]*Position, error) {
    // Implementation
    return []*Position{}, nil
}

// savePosition saves a position to the database
func (t *PositionTracker) savePosition(position *Position) error {
    // Implementation
    return nil
}

// abs returns the absolute value of an integer
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
```

### 2. Portfolio Manager

```go
// portfolio_manager.go
package portfolio

import (
    "sync"
    "time"
    "github.com/trading-platform/backend/internal/database"
)

// Portfolio represents a trading portfolio
type Portfolio struct {
    ID              string
    UserID          string
    Name            string
    Description     string
    Positions       []*Position
    TotalValue      float64
    RealizedProfit  float64
    UnrealizedProfit float64
    LastUpdated     time.Time
}

// PortfolioManager manages portfolios consisting of multiple positions
type PortfolioManager struct {
    positionTracker *PositionTracker
    db              *database.Database
    portfolios      map[string]map[string]*Portfolio // userID -> portfolioID -> Portfolio
    mu              sync.RWMutex
    logger          *Logger
}

// NewPortfolioManager creates a new portfolio manager
func NewPortfolioManager(positionTracker *PositionTracker, db *database.Database) *PortfolioManager {
    return &PortfolioManager{
        positionTracker: positionTracker,
        db:              db,
        portfolios:      make(map[string]map[string]*Portfolio),
        mu:              sync.RWMutex{},
        logger:          NewLogger("portfolio_manager"),
    }
}

// GetPortfolio gets a portfolio by ID
func (m *PortfolioManager) GetPortfolio(userID, portfolioID string) (*Portfolio, error) {
    // Check cache
    m.mu.RLock()
    userPortfolios, ok := m.portfolios[userID]
    if ok {
        portfolio, ok := userPortfolios[portfolioID]
        if ok {
            m.mu.RUnlock()
            return portfolio, nil
        }
    }
    m.mu.RUnlock()
    
    // Get from database
    portfolio, err := m.loadPortfolio(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Cache portfolio
    m.mu.Lock()
    if _, ok := m.portfolios[userID]; !ok {
        m.portfolios[userID] = make(map[string]*Portfolio)
    }
    m.portfolios[userID][portfolioID] = portfolio
    m.mu.Unlock()
    
    return portfolio, nil
}

// GetPortfolios gets all portfolios for a user
func (m *PortfolioManager) GetPortfolios(userID string) ([]*Portfolio, error) {
    // Check cache
    m.mu.RLock()
    userPortfolios, ok := m.portfolios[userID]
    m.mu.RUnlock()
    
    if ok {
        portfolios := make([]*Portfolio, 0, len(userPortfolios))
        for _, portfolio := range userPortfolios {
            portfolios = append(portfolios, portfolio)
        }
        return portfolios, nil
    }
    
    // Get from database
    portfolios, err := m.loadPortfolios(userID)
    if err != nil {
        return nil, err
    }
    
    // Cache portfolios
    m.mu.Lock()
    m.portfolios[userID] = make(map[string]*Portfolio)
    for _, portfolio := range portfolios {
        m.portfolios[userID][portfolio.ID] = portfolio
    }
    m.mu.Unlock()
    
    return portfolios, nil
}

// CreatePortfolio creates a new portfolio
func (m *PortfolioManager) CreatePortfolio(userID, name, description string) (*Portfolio, error) {
    portfolio := &Portfolio{
        ID:          generateID(),
        UserID:      userID,
        Name:        name,
        Description: description,
        Positions:   make([]*Position, 0),
        LastUpdated: time.Now(),
    }
    
    // Save to database
    if err := m.savePortfolio(portfolio); err != nil {
        m.logger.Error("Failed to save portfolio", "error", err, "user_id", userID, "name", name)
        return nil, err
    }
    
    // Cache portfolio
    m.mu.Lock()
    if _, ok := m.portfolios[userID]; !ok {
        m.portfolios[userID] = make(map[string]*Portfolio)
    }
    m.portfolios[userID][portfolio.ID] = portfolio
    m.mu.Unlock()
    
    m.logger.Info("Portfolio created", "user_id", userID, "portfolio_id", portfolio.ID, "name", name)
    return portfolio, nil
}

// UpdatePortfolio updates a portfolio
func (m *PortfolioManager) UpdatePortfolio(userID, portfolioID, name, description string) (*Portfolio, error) {
    // Get portfolio
    portfolio, err := m.GetPortfolio(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Update portfolio
    portfolio.Name = name
    portfolio.Description = description
    portfolio.LastUpdated = time.Now()
    
    // Save to database
    if err := m.savePortfolio(portfolio); err != nil {
        m.logger.Error("Failed to save portfolio", "error", err, "user_id", userID, "portfolio_id", portfolioID)
        return nil, err
    }
    
    // Update cache
    m.mu.Lock()
    m.portfolios[userID][portfolioID] = portfolio
    m.mu.Unlock()
    
    m.logger.Info("Portfolio updated", "user_id", userID, "portfolio_id", portfolioID, "name", name)
    return portfolio, nil
}

// DeletePortfolio deletes a portfolio
func (m *PortfolioManager) DeletePortfolio(userID, portfolioID string) error {
    // Delete from database
    if err := m.deletePortfolio(userID, portfolioID); err != nil {
        m.logger.Error("Failed to delete portfolio", "error", err, "user_id", userID, "portfolio_id", portfolioID)
        return err
    }
    
    // Remove from cache
    m.mu.Lock()
    if userPortfolios, ok := m.portfolios[userID]; ok {
        delete(userPortfolios, portfolioID)
    }
    m.mu.Unlock()
    
    m.logger.Info("Portfolio deleted", "user_id", userID, "portfolio_id", portfolioID)
    return nil
}

// AddPositionToPortfolio adds a position to a portfolio
func (m *PortfolioManager) AddPositionToPortfolio(userID, portfolioID, exchangeSegment, tradingSymbol string) error {
    // Get portfolio
    portfolio, err := m.GetPortfolio(userID, portfolioID)
    if err != nil {
        return err
    }
    
    // Get position
    position, err := m.positionTracker.GetPosition(userID, exchangeSegment, tradingSymbol)
    if err != nil {
        return err
    }
    
    // Check if position already exists in portfolio
    for _, p := range portfolio.Positions {
        if p.ExchangeSegment == exchangeSegment && p.TradingSymbol == tradingSymbol {
            m.logger.Warn("Position already exists in portfolio", "user_id", userID, "portfolio_id", portfolioID, "symbol", tradingSymbol)
            return nil
        }
    }
    
    // Add position to portfolio
    portfolio.Positions = append(portfolio.Positions, position)
    
    // Update portfolio metrics
    m.updatePortfolioMetrics(portfolio)
    
    // Save to database
    if err := m.savePortfolio(portfolio); err != nil {
        m.logger.Error("Failed to save portfolio", "error", err, "user_id", userID, "portfolio_id", portfolioID)
        return err
    }
    
    // Update cache
    m.mu.Lock()
    m.portfolios[userID][portfolioID] = portfolio
    m.mu.Unlock()
    
    m.logger.Info("Position added to portfolio", "user_id", userID, "portfolio_id", portfolioID, "symbol", tradingSymbol)
    return nil
}

// RemovePositionFromPortfolio removes a position from a portfolio
func (m *PortfolioManager) RemovePositionFromPortfolio(userID, portfolioID, exchangeSegment, tradingSymbol string) error {
    // Get portfolio
    portfolio, err := m.GetPortfolio(userID, portfolioID)
    if err != nil {
        return err
    }
    
    // Remove position from portfolio
    for i, p := range portfolio.Positions {
        if p.ExchangeSegment == exchangeSegment && p.TradingSymbol == tradingSymbol {
            portfolio.Positions = append(portfolio.Positions[:i], portfolio.Positions[i+1:]...)
            break
        }
    }
    
    // Update portfolio metrics
    m.updatePortfolioMetrics(portfolio)
    
    // Save to database
    if err := m.savePortfolio(portfolio); err != nil {
        m.logger.Error("Failed to save portfolio", "error", err, "user_id", userID, "portfolio_id", portfolioID)
        return err
    }
    
    // Update cache
    m.mu.Lock()
    m.portfolios[userID][portfolioID] = portfolio
    m.mu.Unlock()
    
    m.logger.Info("Position removed from portfolio", "user_id", userID, "portfolio_id", portfolioID, "symbol", tradingSymbol)
    return nil
}

// UpdatePortfolioMetrics updates the metrics for a portfolio
func (m *PortfolioManager) UpdatePortfolioMetrics(userID, portfolioID string) error {
    // Get portfolio
    portfolio, err := m.GetPortfolio(userID, portfolioID)
    if err != nil {
        return err
    }
    
    // Update portfolio metrics
    m.updatePortfolioMetrics(portfolio)
    
    // Save to database
    if err := m.savePortfolio(portfolio); err != nil {
        m.logger.Error("Failed to save portfolio", "error", err, "user_id", userID, "portfolio_id", portfolioID)
        return err
    }
    
    // Update cache
    m.mu.Lock()
    m.portfolios[userID][portfolioID] = portfolio
    m.mu.Unlock()
    
    return nil
}

// updatePortfolioMetrics updates the metrics for a portfolio
func (m *PortfolioManager) updatePortfolioMetrics(portfolio *Portfolio) {
    totalValue := 0.0
    realizedProfit := 0.0
    unrealizedProfit := 0.0
    
    for _, position := range portfolio.Positions {
        totalValue += position.Value
        realizedProfit += position.RealizedProfit
        unrealizedProfit += position.UnrealizedProfit
    }
    
    portfolio.TotalValue = totalValue
    portfolio.RealizedProfit = realizedProfit
    portfolio.UnrealizedProfit = unrealizedProfit
    portfolio.LastUpdated = time.Now()
}

// loadPortfolio loads a portfolio from the database
func (m *PortfolioManager) loadPortfolio(userID, portfolioID string) (*Portfolio, error) {
    // Implementation
    return &Portfolio{}, nil
}

// loadPortfolios loads all portfolios for a user from the database
func (m *PortfolioManager) loadPortfolios(userID string) ([]*Portfolio, error) {
    // Implementation
    return []*Portfolio{}, nil
}

// savePortfolio saves a portfolio to the database
func (m *PortfolioManager) savePortfolio(portfolio *Portfolio) error {
    // Implementation
    return nil
}

// deletePortfolio deletes a portfolio from the database
func (m *PortfolioManager) deletePortfolio(userID, portfolioID string) error {
    // Implementation
    return nil
}

// generateID generates a unique ID
func generateID() string {
    return "portfolio_" + time.Now().Format("20060102150405")
}
```

### 3. Performance Calculator

```go
// performance_calculator.go
package portfolio

import (
    "math"
    "time"
)

// PerformanceMetrics represents performance metrics for a position or portfolio
type PerformanceMetrics struct {
    TotalReturn      float64
    PercentReturn    float64
    AnnualizedReturn float64
    Volatility       float64
    SharpeRatio      float64
    MaxDrawdown      float64
    WinRate          float64
    ProfitFactor     float64
}

// PerformanceCalculator calculates performance metrics for positions and portfolios
type PerformanceCalculator struct {
    positionTracker *PositionTracker
    portfolioManager *PortfolioManager
    db              *database.Database
    logger          *Logger
}

// NewPerformanceCalculator creates a new performance calculator
func NewPerformanceCalculator(positionTracker *PositionTracker, portfolioManager *PortfolioManager, db *database.Database) *PerformanceCalculator {
    return &PerformanceCalculator{
        positionTracker: positionTracker,
        portfolioManager: portfolioManager,
        db:              db,
        logger:          NewLogger("performance_calculator"),
    }
}

// CalculatePositionPerformance calculates performance metrics for a position
func (c *PerformanceCalculator) CalculatePositionPerformance(userID, exchangeSegment, tradingSymbol string, startDate, endDate time.Time) (*PerformanceMetrics, error) {
    // Get position
    position, err := c.positionTracker.GetPosition(userID, exchangeSegment, tradingSymbol)
    if err != nil {
        return nil, err
    }
    
    // Get historical data for position
    historicalData, err := c.getHistoricalPositionData(userID, exchangeSegment, tradingSymbol, startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    // Calculate performance metrics
    metrics := c.calculateMetrics(historicalData)
    
    return metrics, nil
}

// CalculatePortfolioPerformance calculates performance metrics for a portfolio
func (c *PerformanceCalculator) CalculatePortfolioPerformance(userID, portfolioID string, startDate, endDate time.Time) (*PerformanceMetrics, error) {
    // Get portfolio
    portfolio, err := c.portfolioManager.GetPortfolio(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Get historical data for portfolio
    historicalData, err := c.getHistoricalPortfolioData(userID, portfolioID, startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    // Calculate performance metrics
    metrics := c.calculateMetrics(historicalData)
    
    return metrics, nil
}

// calculateMetrics calculates performance metrics from historical data
func (c *PerformanceCalculator) calculateMetrics(historicalData []HistoricalDataPoint) *PerformanceMetrics {
    if len(historicalData) < 2 {
        return &PerformanceMetrics{}
    }
    
    // Calculate returns
    returns := make([]float64, len(historicalData)-1)
    for i := 1; i < len(historicalData); i++ {
        returns[i-1] = (historicalData[i].Value - historicalData[i-1].Value) / historicalData[i-1].Value
    }
    
    // Calculate total return
    initialValue := historicalData[0].Value
    finalValue := historicalData[len(historicalData)-1].Value
    totalReturn := finalValue - initialValue
    percentReturn := (finalValue / initialValue) - 1
    
    // Calculate annualized return
    startDate := historicalData[0].Date
    endDate := historicalData[len(historicalData)-1].Date
    years := endDate.Sub(startDate).Hours() / 24 / 365
    annualizedReturn := math.Pow(1+percentReturn, 1/years) - 1
    
    // Calculate volatility
    meanReturn := mean(returns)
    variance := 0.0
    for _, r := range returns {
        variance += math.Pow(r-meanReturn, 2)
    }
    variance /= float64(len(returns))
    volatility := math.Sqrt(variance) * math.Sqrt(252) // Annualized volatility
    
    // Calculate Sharpe ratio
    riskFreeRate := 0.02 // Assume 2% risk-free rate
    sharpeRatio := (annualizedReturn - riskFreeRate) / volatility
    
    // Calculate maximum drawdown
    maxDrawdown := 0.0
    peak := initialValue
    for _, point := range historicalData {
        if point.Value > peak {
            peak = point.Value
        }
        drawdown := (peak - point.Value) / peak
        if drawdown > maxDrawdown {
            maxDrawdown = drawdown
        }
    }
    
    // Calculate win rate and profit factor
    wins := 0
    totalProfit := 0.0
    totalLoss := 0.0
    for _, r := range returns {
        if r > 0 {
            wins++
            totalProfit += r
        } else {
            totalLoss -= r
        }
    }
    winRate := float64(wins) / float64(len(returns))
    profitFactor := 0.0
    if totalLoss > 0 {
        profitFactor = totalProfit / totalLoss
    }
    
    return &PerformanceMetrics{
        TotalReturn:      totalReturn,
        PercentReturn:    percentReturn,
        AnnualizedReturn: annualizedReturn,
        Volatility:       volatility,
        SharpeRatio:      sharpeRatio,
        MaxDrawdown:      maxDrawdown,
        WinRate:          winRate,
        ProfitFactor:     profitFactor,
    }
}

// HistoricalDataPoint represents a historical data point
type HistoricalDataPoint struct {
    Date  time.Time
    Value float64
}

// getHistoricalPositionData gets historical data for a position
func (c *PerformanceCalculator) getHistoricalPositionData(userID, exchangeSegment, tradingSymbol string, startDate, endDate time.Time) ([]HistoricalDataPoint, error) {
    // Implementation
    return []HistoricalDataPoint{}, nil
}

// getHistoricalPortfolioData gets historical data for a portfolio
func (c *PerformanceCalculator) getHistoricalPortfolioData(userID, portfolioID string, startDate, endDate time.Time) ([]HistoricalDataPoint, error) {
    // Implementation
    return []HistoricalDataPoint{}, nil
}

// mean calculates the mean of a slice of float64
func mean(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}
```

### 4. Portfolio Analytics Engine

```go
// portfolio_analytics.go
package portfolio

import (
    "math"
    "time"
)

// PortfolioAnalytics represents analytics for a portfolio
type PortfolioAnalytics struct {
    RiskMetrics      RiskMetrics
    AllocationMetrics AllocationMetrics
    CorrelationMatrix [][]float64
    BetaToMarket     float64
    AlphaToMarket    float64
    ValueAtRisk      float64
}

// RiskMetrics represents risk metrics for a portfolio
type RiskMetrics struct {
    Volatility       float64
    SharpeRatio      float64
    SortinoRatio     float64
    MaxDrawdown      float64
    DownsideDeviation float64
    InformationRatio float64
}

// AllocationMetrics represents allocation metrics for a portfolio
type AllocationMetrics struct {
    SectorAllocation     map[string]float64
    AssetClassAllocation map[string]float64
    GeographicAllocation map[string]float64
    CurrencyAllocation   map[string]float64
}

// PortfolioAnalyticsEngine provides advanced analytics for portfolios
type PortfolioAnalyticsEngine struct {
    positionTracker     *PositionTracker
    portfolioManager    *PortfolioManager
    performanceCalculator *PerformanceCalculator
    db                  *database.Database
    logger              *Logger
}

// NewPortfolioAnalyticsEngine creates a new portfolio analytics engine
func NewPortfolioAnalyticsEngine(
    positionTracker *PositionTracker,
    portfolioManager *PortfolioManager,
    performanceCalculator *PerformanceCalculator,
    db *database.Database,
) *PortfolioAnalyticsEngine {
    return &PortfolioAnalyticsEngine{
        positionTracker:     positionTracker,
        portfolioManager:    portfolioManager,
        performanceCalculator: performanceCalculator,
        db:                  db,
        logger:              NewLogger("portfolio_analytics"),
    }
}

// CalculatePortfolioAnalytics calculates analytics for a portfolio
func (e *PortfolioAnalyticsEngine) CalculatePortfolioAnalytics(userID, portfolioID string) (*PortfolioAnalytics, error) {
    // Get portfolio
    portfolio, err := e.portfolioManager.GetPortfolio(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Calculate risk metrics
    riskMetrics, err := e.calculateRiskMetrics(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Calculate allocation metrics
    allocationMetrics, err := e.calculateAllocationMetrics(portfolio)
    if err != nil {
        return nil, err
    }
    
    // Calculate correlation matrix
    correlationMatrix, err := e.calculateCorrelationMatrix(portfolio)
    if err != nil {
        return nil, err
    }
    
    // Calculate beta to market
    betaToMarket, err := e.calculateBetaToMarket(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Calculate alpha to market
    alphaToMarket, err := e.calculateAlphaToMarket(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Calculate value at risk
    valueAtRisk, err := e.calculateValueAtRisk(portfolio)
    if err != nil {
        return nil, err
    }
    
    return &PortfolioAnalytics{
        RiskMetrics:      *riskMetrics,
        AllocationMetrics: *allocationMetrics,
        CorrelationMatrix: correlationMatrix,
        BetaToMarket:     betaToMarket,
        AlphaToMarket:    alphaToMarket,
        ValueAtRisk:      valueAtRisk,
    }, nil
}

// calculateRiskMetrics calculates risk metrics for a portfolio
func (e *PortfolioAnalyticsEngine) calculateRiskMetrics(userID, portfolioID string) (*RiskMetrics, error) {
    // Get historical data
    startDate := time.Now().AddDate(-1, 0, 0) // 1 year ago
    endDate := time.Now()
    historicalData, err := e.getHistoricalPortfolioData(userID, portfolioID, startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    // Calculate returns
    returns := make([]float64, len(historicalData)-1)
    for i := 1; i < len(historicalData); i++ {
        returns[i-1] = (historicalData[i].Value - historicalData[i-1].Value) / historicalData[i-1].Value
    }
    
    // Calculate volatility
    meanReturn := mean(returns)
    variance := 0.0
    for _, r := range returns {
        variance += math.Pow(r-meanReturn, 2)
    }
    variance /= float64(len(returns))
    volatility := math.Sqrt(variance) * math.Sqrt(252) // Annualized volatility
    
    // Calculate Sharpe ratio
    riskFreeRate := 0.02 // Assume 2% risk-free rate
    sharpeRatio := (meanReturn*252 - riskFreeRate) / volatility
    
    // Calculate Sortino ratio
    negativeReturns := make([]float64, 0)
    for _, r := range returns {
        if r < 0 {
            negativeReturns = append(negativeReturns, r)
        }
    }
    downsideVariance := 0.0
    for _, r := range negativeReturns {
        downsideVariance += math.Pow(r, 2)
    }
    downsideVariance /= float64(len(negativeReturns))
    downsideDeviation := math.Sqrt(downsideVariance) * math.Sqrt(252)
    sortinoRatio := (meanReturn*252 - riskFreeRate) / downsideDeviation
    
    // Calculate maximum drawdown
    maxDrawdown := 0.0
    peak := historicalData[0].Value
    for _, point := range historicalData {
        if point.Value > peak {
            peak = point.Value
        }
        drawdown := (peak - point.Value) / peak
        if drawdown > maxDrawdown {
            maxDrawdown = drawdown
        }
    }
    
    // Calculate information ratio
    // This is a simplified calculation; in practice, you would use a benchmark
    benchmarkReturns := make([]float64, len(returns))
    for i := range benchmarkReturns {
        benchmarkReturns[i] = 0.0001 // Assume a constant daily return for the benchmark
    }
    
    excessReturns := make([]float64, len(returns))
    for i := range returns {
        excessReturns[i] = returns[i] - benchmarkReturns[i]
    }
    
    trackingError := 0.0
    meanExcessReturn := mean(excessReturns)
    for _, r := range excessReturns {
        trackingError += math.Pow(r-meanExcessReturn, 2)
    }
    trackingError = math.Sqrt(trackingError/float64(len(excessReturns))) * math.Sqrt(252)
    
    informationRatio := (meanReturn*252 - mean(benchmarkReturns)*252) / trackingError
    
    return &RiskMetrics{
        Volatility:       volatility,
        SharpeRatio:      sharpeRatio,
        SortinoRatio:     sortinoRatio,
        MaxDrawdown:      maxDrawdown,
        DownsideDeviation: downsideDeviation,
        InformationRatio: informationRatio,
    }, nil
}

// calculateAllocationMetrics calculates allocation metrics for a portfolio
func (e *PortfolioAnalyticsEngine) calculateAllocationMetrics(portfolio *Portfolio) (*AllocationMetrics, error) {
    // Implementation
    return &AllocationMetrics{}, nil
}

// calculateCorrelationMatrix calculates the correlation matrix for a portfolio
func (e *PortfolioAnalyticsEngine) calculateCorrelationMatrix(portfolio *Portfolio) ([][]float64, error) {
    // Implementation
    return [][]float64{}, nil
}

// calculateBetaToMarket calculates the beta to market for a portfolio
func (e *PortfolioAnalyticsEngine) calculateBetaToMarket(userID, portfolioID string) (float64, error) {
    // Implementation
    return 0.0, nil
}

// calculateAlphaToMarket calculates the alpha to market for a portfolio
func (e *PortfolioAnalyticsEngine) calculateAlphaToMarket(userID, portfolioID string) (float64, error) {
    // Implementation
    return 0.0, nil
}

// calculateValueAtRisk calculates the value at risk for a portfolio
func (e *PortfolioAnalyticsEngine) calculateValueAtRisk(portfolio *Portfolio) (float64, error) {
    // Implementation
    return 0.0, nil
}

// getHistoricalPortfolioData gets historical data for a portfolio
func (e *PortfolioAnalyticsEngine) getHistoricalPortfolioData(userID, portfolioID string, startDate, endDate time.Time) ([]HistoricalDataPoint, error) {
    // Implementation
    return []HistoricalDataPoint{}, nil
}
```

### 5. Reporting System

```go
// reporting_system.go
package portfolio

import (
    "time"
    "github.com/trading-platform/backend/internal/database"
)

// Report represents a report
type Report struct {
    ID          string
    UserID      string
    Type        string
    Title       string
    Description string
    Content     map[string]interface{}
    CreatedAt   time.Time
}

// ReportingSystem generates reports for positions and portfolios
type ReportingSystem struct {
    positionTracker     *PositionTracker
    portfolioManager    *PortfolioManager
    performanceCalculator *PerformanceCalculator
    analyticsEngine     *PortfolioAnalyticsEngine
    db                  *database.Database
    logger              *Logger
}

// NewReportingSystem creates a new reporting system
func NewReportingSystem(
    positionTracker *PositionTracker,
    portfolioManager *PortfolioManager,
    performanceCalculator *PerformanceCalculator,
    analyticsEngine *PortfolioAnalyticsEngine,
    db *database.Database,
) *ReportingSystem {
    return &ReportingSystem{
        positionTracker:     positionTracker,
        portfolioManager:    portfolioManager,
        performanceCalculator: performanceCalculator,
        analyticsEngine:     analyticsEngine,
        db:                  db,
        logger:              NewLogger("reporting_system"),
    }
}

// GeneratePositionReport generates a report for a position
func (s *ReportingSystem) GeneratePositionReport(userID, exchangeSegment, tradingSymbol string) (*Report, error) {
    // Get position
    position, err := s.positionTracker.GetPosition(userID, exchangeSegment, tradingSymbol)
    if err != nil {
        return nil, err
    }
    
    // Calculate performance metrics
    startDate := time.Now().AddDate(-1, 0, 0) // 1 year ago
    endDate := time.Now()
    metrics, err := s.performanceCalculator.CalculatePositionPerformance(userID, exchangeSegment, tradingSymbol, startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    // Create report
    report := &Report{
        ID:          generateReportID(),
        UserID:      userID,
        Type:        "position",
        Title:       "Position Report: " + tradingSymbol,
        Description: "Performance report for " + tradingSymbol,
        Content: map[string]interface{}{
            "position":  position,
            "metrics":   metrics,
            "startDate": startDate,
            "endDate":   endDate,
        },
        CreatedAt:   time.Now(),
    }
    
    // Save report
    if err := s.saveReport(report); err != nil {
        s.logger.Error("Failed to save report", "error", err, "user_id", userID, "symbol", tradingSymbol)
        return nil, err
    }
    
    return report, nil
}

// GeneratePortfolioReport generates a report for a portfolio
func (s *ReportingSystem) GeneratePortfolioReport(userID, portfolioID string) (*Report, error) {
    // Get portfolio
    portfolio, err := s.portfolioManager.GetPortfolio(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Calculate performance metrics
    startDate := time.Now().AddDate(-1, 0, 0) // 1 year ago
    endDate := time.Now()
    metrics, err := s.performanceCalculator.CalculatePortfolioPerformance(userID, portfolioID, startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    // Calculate analytics
    analytics, err := s.analyticsEngine.CalculatePortfolioAnalytics(userID, portfolioID)
    if err != nil {
        return nil, err
    }
    
    // Create report
    report := &Report{
        ID:          generateReportID(),
        UserID:      userID,
        Type:        "portfolio",
        Title:       "Portfolio Report: " + portfolio.Name,
        Description: "Performance report for " + portfolio.Name,
        Content: map[string]interface{}{
            "portfolio": portfolio,
            "metrics":   metrics,
            "analytics": analytics,
            "startDate": startDate,
            "endDate":   endDate,
        },
        CreatedAt:   time.Now(),
    }
    
    // Save report
    if err := s.saveReport(report); err != nil {
        s.logger.Error("Failed to save report", "error", err, "user_id", userID, "portfolio_id", portfolioID)
        return nil, err
    }
    
    return report, nil
}

// GetReport gets a report by ID
func (s *ReportingSystem) GetReport(userID, reportID string) (*Report, error) {
    // Implementation
    return &Report{}, nil
}

// GetReports gets all reports for a user
func (s *ReportingSystem) GetReports(userID string) ([]*Report, error) {
    // Implementation
    return []*Report{}, nil
}

// DeleteReport deletes a report
func (s *ReportingSystem) DeleteReport(userID, reportID string) error {
    // Implementation
    return nil
}

// saveReport saves a report to the database
func (s *ReportingSystem) saveReport(report *Report) error {
    // Implementation
    return nil
}

// generateReportID generates a unique ID for a report
func generateReportID() string {
    return "report_" + time.Now().Format("20060102150405")
}
```

## API Integration

The Position and Portfolio Management System will be integrated with the existing API layer:

```go
// portfolio_controller.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/trading-platform/backend/internal/portfolio"
)

// PortfolioController handles portfolio management requests
type PortfolioController struct {
    positionTracker     *portfolio.PositionTracker
    portfolioManager    *portfolio.PortfolioManager
    performanceCalculator *portfolio.PerformanceCalculator
    analyticsEngine     *portfolio.PortfolioAnalyticsEngine
    reportingSystem     *portfolio.ReportingSystem
}

// NewPortfolioController creates a new portfolio controller
func NewPortfolioController(
    positionTracker *portfolio.PositionTracker,
    portfolioManager *portfolio.PortfolioManager,
    performanceCalculator *portfolio.PerformanceCalculator,
    analyticsEngine *portfolio.PortfolioAnalyticsEngine,
    reportingSystem *portfolio.ReportingSystem,
) *PortfolioController {
    return &PortfolioController{
        positionTracker:     positionTracker,
        portfolioManager:    portfolioManager,
        performanceCalculator: performanceCalculator,
        analyticsEngine:     analyticsEngine,
        reportingSystem:     reportingSystem,
    }
}

// RegisterRoutes registers the API routes
func (c *PortfolioController) RegisterRoutes(router *gin.Engine) {
    group := router.Group("/api/portfolio")
    
    // Position routes
    group.GET("/positions", c.GetPositions)
    group.GET("/positions/:exchangeSegment/:tradingSymbol", c.GetPosition)
    group.GET("/positions/:exchangeSegment/:tradingSymbol/performance", c.GetPositionPerformance)
    
    // Portfolio routes
    group.GET("/portfolios", c.GetPortfolios)
    group.POST("/portfolios", c.CreatePortfolio)
    group.GET("/portfolios/:portfolioID", c.GetPortfolio)
    group.PUT("/portfolios/:portfolioID", c.UpdatePortfolio)
    group.DELETE("/portfolios/:portfolioID", c.DeletePortfolio)
    group.POST("/portfolios/:portfolioID/positions", c.AddPositionToPortfolio)
    group.DELETE("/portfolios/:portfolioID/positions/:exchangeSegment/:tradingSymbol", c.RemovePositionFromPortfolio)
    group.GET("/portfolios/:portfolioID/performance", c.GetPortfolioPerformance)
    group.GET("/portfolios/:portfolioID/analytics", c.GetPortfolioAnalytics)
    
    // Report routes
    group.GET("/reports", c.GetReports)
    group.GET("/reports/:reportID", c.GetReport)
    group.POST("/reports/positions/:exchangeSegment/:tradingSymbol", c.GeneratePositionReport)
    group.POST("/reports/portfolios/:portfolioID", c.GeneratePortfolioReport)
    group.DELETE("/reports/:reportID", c.DeleteReport)
}

// GetPositions handles get positions requests
func (c *PortfolioController) GetPositions(ctx *gin.Context) {
    // Implementation
}

// GetPosition handles get position requests
func (c *PortfolioController) GetPosition(ctx *gin.Context) {
    // Implementation
}

// GetPositionPerformance handles get position performance requests
func (c *PortfolioController) GetPositionPerformance(ctx *gin.Context) {
    // Implementation
}

// GetPortfolios handles get portfolios requests
func (c *PortfolioController) GetPortfolios(ctx *gin.Context) {
    // Implementation
}

// CreatePortfolio handles create portfolio requests
func (c *PortfolioController) CreatePortfolio(ctx *gin.Context) {
    // Implementation
}

// GetPortfolio handles get portfolio requests
func (c *PortfolioController) GetPortfolio(ctx *gin.Context) {
    // Implementation
}

// UpdatePortfolio handles update portfolio requests
func (c *PortfolioController) UpdatePortfolio(ctx *gin.Context) {
    // Implementation
}

// DeletePortfolio handles delete portfolio requests
func (c *PortfolioController) DeletePortfolio(ctx *gin.Context) {
    // Implementation
}

// AddPositionToPortfolio handles add position to portfolio requests
func (c *PortfolioController) AddPositionToPortfolio(ctx *gin.Context) {
    // Implementation
}

// RemovePositionFromPortfolio handles remove position from portfolio requests
func (c *PortfolioController) RemovePositionFromPortfolio(ctx *gin.Context) {
    // Implementation
}

// GetPortfolioPerformance handles get portfolio performance requests
func (c *PortfolioController) GetPortfolioPerformance(ctx *gin.Context) {
    // Implementation
}

// GetPortfolioAnalytics handles get portfolio analytics requests
func (c *PortfolioController) GetPortfolioAnalytics(ctx *gin.Context) {
    // Implementation
}

// GetReports handles get reports requests
func (c *PortfolioController) GetReports(ctx *gin.Context) {
    // Implementation
}

// GetReport handles get report requests
func (c *PortfolioController) GetReport(ctx *gin.Context) {
    // Implementation
}

// GeneratePositionReport handles generate position report requests
func (c *PortfolioController) GeneratePositionReport(ctx *gin.Context) {
    // Implementation
}

// GeneratePortfolioReport handles generate portfolio report requests
func (c *PortfolioController) GeneratePortfolioReport(ctx *gin.Context) {
    // Implementation
}

// DeleteReport handles delete report requests
func (c *PortfolioController) DeleteReport(ctx *gin.Context) {
    // Implementation
}
```

## WebSocket Integration

The Position and Portfolio Management System will provide real-time updates via WebSocket:

```go
// portfolio_websocket.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/trading-platform/backend/internal/portfolio"
)

// PortfolioWebSocketHandler handles WebSocket connections for portfolio updates
type PortfolioWebSocketHandler struct {
    positionTracker *portfolio.PositionTracker
    portfolioManager *portfolio.PortfolioManager
    upgrader        websocket.Upgrader
    logger          *Logger
}

// NewPortfolioWebSocketHandler creates a new portfolio WebSocket handler
func NewPortfolioWebSocketHandler(
    positionTracker *portfolio.PositionTracker,
    portfolioManager *portfolio.PortfolioManager,
) *PortfolioWebSocketHandler {
    return &PortfolioWebSocketHandler{
        positionTracker: positionTracker,
        portfolioManager: portfolioManager,
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
        logger: NewLogger("portfolio_websocket"),
    }
}

// RegisterRoutes registers the WebSocket routes
func (h *PortfolioWebSocketHandler) RegisterRoutes(router *gin.Engine) {
    router.GET("/ws/portfolio", h.HandleWebSocket)
}

// HandleWebSocket handles WebSocket connections
func (h *PortfolioWebSocketHandler) HandleWebSocket(ctx *gin.Context) {
    // Implementation
}
```

## Testing Strategy

1. **Unit Tests**: Test individual components in isolation
   - Test position tracking
   - Test portfolio management
   - Test performance calculation
   - Test portfolio analytics

2. **Integration Tests**: Test the integration between components
   - Test position updates and portfolio updates
   - Test performance calculation with historical data
   - Test reporting system

3. **Scenario Tests**: Test specific scenarios
   - Test portfolio creation and management
   - Test position tracking with trades
   - Test performance calculation with different time periods

## Performance Considerations

1. **Caching**: Cache positions and portfolios for fast access
2. **Efficient Data Structures**: Use efficient data structures for position and portfolio tracking
3. **Batch Processing**: Batch position updates for efficient processing
4. **Asynchronous Processing**: Use asynchronous processing for performance calculation and analytics
5. **Database Optimization**: Optimize database queries for position and portfolio data

## Error Handling

1. **Validation Errors**: Validate inputs before processing
2. **Database Errors**: Handle database errors gracefully
3. **Calculation Errors**: Handle calculation errors and edge cases
4. **Logging**: Comprehensive logging for debugging and monitoring
5. **Fallback Mechanisms**: Implement fallback mechanisms for system failures

## Implementation Timeline

1. **Week 1**: Implement Position Tracker and Portfolio Manager
2. **Week 2**: Implement Performance Calculator
3. **Week 3**: Implement Portfolio Analytics Engine
4. **Week 4**: Implement Reporting System
5. **Week 5**: Integrate with API and WebSocket layers

## Conclusion

This Position and Portfolio Management System implementation plan provides a comprehensive approach to position tracking, portfolio management, performance calculation, and portfolio analytics. The implementation follows a modular architecture that allows for flexibility and extensibility, with support for real-time updates and reporting.
