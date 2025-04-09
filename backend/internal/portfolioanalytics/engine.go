package portfolioanalytics

import (
        "context"
        "errors"
        "fmt"
        "math"
        "sync"
        "time"
)

// PortfolioAnalyticsEngine is the main engine for portfolio analytics
type PortfolioAnalyticsEngine struct {
        portfolios       map[string]*Portfolio
        positions        map[string][]*Position
        performanceCache map[string]*PerformanceMetrics
        riskCache        map[string]*RiskMetrics
        mutex            sync.RWMutex
        dataProvider     DataProvider
        calculationQueue chan *AnalyticsTask
        workers          int
        isRunning        bool
        stopChan         chan struct{}
}

// Portfolio represents a collection of positions
type Portfolio struct {
        ID               string
        Name             string
        Description      string
        Positions        []*Position
        Tags             []string
        CreatedAt        time.Time
        UpdatedAt        time.Time
        StrategyID       string
        UserID           string
        PerformanceCache *PerformanceMetrics
        RiskCache        *RiskMetrics
}

// Position represents a trading position
type Position struct {
        ID              string
        Symbol          string
        Quantity        int
        EntryPrice      float64
        CurrentPrice    float64
        EntryTime       time.Time
        ExitTime        *time.Time
        ExitPrice       *float64
        TransactionType string // "BUY" or "SELL"
        ProductType     string // "MIS", "NRML", "CNC"
        Exchange        string
        ExpiryDate      *time.Time
        StrikePrice     *float64
        OptionType      *string // "CE" or "PE"
        PortfolioID     string
        StrategyID      string
        Tags            []string
        Greeks          *Greeks
}

// Greeks represents option Greeks
type Greeks struct {
        Delta     float64
        Gamma     float64
        Theta     float64
        Vega      float64
        Rho       float64
        UpdatedAt time.Time
}

// PerformanceMetrics represents performance metrics for a portfolio
type PerformanceMetrics struct {
        TotalPnL            float64
        RealizedPnL         float64
        UnrealizedPnL       float64
        PnLPercentage       float64
        CAGR                float64
        Volatility          float64
        SharpeRatio         float64
        SortinoRatio        float64
        MaxDrawdown         float64
        WinRate             float64
        ProfitFactor        float64
        AverageWin          float64
        AverageLoss         float64
        ExpectancyRatio     float64
        ReturnOnCapital     float64
        DailyPnL            map[string]float64
        CumulativePnL       map[string]float64
        RollingPerformance  map[string]float64
        PerformanceBySymbol map[string]float64
        UpdatedAt           time.Time
}

// RiskMetrics represents risk metrics for a portfolio
type RiskMetrics struct {
        ValueAtRisk         float64
        ConditionalVaR      float64
        BetaToMarket        float64
        PortfolioVolatility float64
        CorrelationMatrix   map[string]map[string]float64
        StressTestResults   map[string]float64
        SectorExposure      map[string]float64
        AssetClassExposure  map[string]float64
        ConcentrationRisk   float64
        LiquidityRisk       float64
        OptionExposure      map[string]float64
        DeltaExposure       float64
        GammaExposure       float64
        ThetaExposure       float64
        VegaExposure        float64
        RhoExposure         float64
        UpdatedAt           time.Time
}

// AnalyticsTask represents a task for the analytics engine
type AnalyticsTask struct {
        TaskType    string
        PortfolioID string
        Callback    func(interface{}, error)
}

// DataProvider interface for getting market data
type DataProvider interface {
        GetCurrentPrice(ctx context.Context, symbol string, exchange string) (float64, error)
        GetHistoricalPrices(ctx context.Context, symbol string, exchange string, startDate time.Time, endDate time.Time, interval string) (map[time.Time]float64, error)
        GetOptionChain(ctx context.Context, symbol string, exchange string, expiryDate time.Time) ([]*OptionData, error)
        GetGreeks(ctx context.Context, symbol string, exchange string, strikePrice float64, expiryDate time.Time, optionType string) (*Greeks, error)
        GetMarketIndices(ctx context.Context) (map[string]float64, error)
        GetVolatilityIndex(ctx context.Context, symbol string) (float64, error)
}

// OptionData represents data for an option
type OptionData struct {
        Symbol      string
        StrikePrice float64
        ExpiryDate  time.Time
        OptionType  string
        LastPrice   float64
        BidPrice    float64
        AskPrice    float64
        Volume      int
        OpenInterest int
        ImpliedVolatility float64
        Delta       float64
        Gamma       float64
        Theta       float64
        Vega        float64
        Rho         float64
}

// NewPortfolioAnalyticsEngine creates a new portfolio analytics engine
func NewPortfolioAnalyticsEngine(dataProvider DataProvider, workers int) *PortfolioAnalyticsEngine {
        return &PortfolioAnalyticsEngine{
                portfolios:       make(map[string]*Portfolio),
                positions:        make(map[string][]*Position),
                performanceCache: make(map[string]*PerformanceMetrics),
                riskCache:        make(map[string]*RiskMetrics),
                dataProvider:     dataProvider,
                calculationQueue: make(chan *AnalyticsTask, 1000),
                workers:          workers,
                stopChan:         make(chan struct{}),
        }
}

// Start starts the portfolio analytics engine
func (e *PortfolioAnalyticsEngine) Start() error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        if e.isRunning {
                return errors.New("analytics engine is already running")
        }

        e.isRunning = true
        e.stopChan = make(chan struct{})

        // Start worker goroutines
        for i := 0; i < e.workers; i++ {
                go e.worker()
        }

        return nil
}

// Stop stops the portfolio analytics engine
func (e *PortfolioAnalyticsEngine) Stop() {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        if !e.isRunning {
                return
        }

        close(e.stopChan)
        e.isRunning = false
}

// worker processes tasks from the calculation queue
func (e *PortfolioAnalyticsEngine) worker() {
        for {
                select {
                case <-e.stopChan:
                        return
                case task := <-e.calculationQueue:
                        var result interface{}
                        var err error

                        switch task.TaskType {
                        case "performance":
                                result, err = e.calculatePerformanceMetrics(task.PortfolioID)
                        case "risk":
                                result, err = e.calculateRiskMetrics(task.PortfolioID)
                        case "update_prices":
                                err = e.updatePositionPrices(task.PortfolioID)
                                result = nil
                        case "update_greeks":
                                err = e.updatePositionGreeks(task.PortfolioID)
                                result = nil
                        }

                        if task.Callback != nil {
                                task.Callback(result, err)
                        }
                }
        }
}

// AddPortfolio adds a portfolio to the engine
func (e *PortfolioAnalyticsEngine) AddPortfolio(portfolio *Portfolio) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        if portfolio.ID == "" {
                return errors.New("portfolio ID cannot be empty")
        }

        if _, exists := e.portfolios[portfolio.ID]; exists {
                return fmt.Errorf("portfolio with ID %s already exists", portfolio.ID)
        }

        e.portfolios[portfolio.ID] = portfolio
        e.positions[portfolio.ID] = portfolio.Positions

        return nil
}

// GetPortfolio gets a portfolio by ID
func (e *PortfolioAnalyticsEngine) GetPortfolio(portfolioID string) (*Portfolio, error) {
        e.mutex.RLock()
        defer e.mutex.RUnlock()

        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return nil, fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        return portfolio, nil
}

// UpdatePortfolio updates a portfolio
func (e *PortfolioAnalyticsEngine) UpdatePortfolio(portfolio *Portfolio) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        if portfolio.ID == "" {
                return errors.New("portfolio ID cannot be empty")
        }

        if _, exists := e.portfolios[portfolio.ID]; !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolio.ID)
        }

        e.portfolios[portfolio.ID] = portfolio
        e.positions[portfolio.ID] = portfolio.Positions

        return nil
}

// DeletePortfolio deletes a portfolio
func (e *PortfolioAnalyticsEngine) DeletePortfolio(portfolioID string) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        if _, exists := e.portfolios[portfolioID]; !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        delete(e.portfolios, portfolioID)
        delete(e.positions, portfolioID)
        delete(e.performanceCache, portfolioID)
        delete(e.riskCache, portfolioID)

        return nil
}

// AddPosition adds a position to a portfolio
func (e *PortfolioAnalyticsEngine) AddPosition(portfolioID string, position *Position) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        position.PortfolioID = portfolioID
        portfolio.Positions = append(portfolio.Positions, position)
        e.positions[portfolioID] = portfolio.Positions

        // Invalidate cache
        delete(e.performanceCache, portfolioID)
        delete(e.riskCache, portfolioID)

        return nil
}

// UpdatePosition updates a position
func (e *PortfolioAnalyticsEngine) UpdatePosition(position *Position) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        portfolioID := position.PortfolioID
        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        found := false
        for i, p := range portfolio.Positions {
                if p.ID == position.ID {
                        portfolio.Positions[i] = position
                        found = true
                        break
                }
        }

        if !found {
                return fmt.Errorf("position with ID %s not found in portfolio %s", position.ID, portfolioID)
        }

        e.positions[portfolioID] = portfolio.Positions

        // Invalidate cache
        delete(e.performanceCache, portfolioID)
        delete(e.riskCache, portfolioID)

        return nil
}

// DeletePosition deletes a position from a portfolio
func (e *PortfolioAnalyticsEngine) DeletePosition(portfolioID string, positionID string) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        found := false
        newPositions := make([]*Position, 0, len(portfolio.Positions))
        for _, p := range portfolio.Positions {
                if p.ID == positionID {
                        found = true
                } else {
                        newPositions = append(newPositions, p)
                }
        }

        if !found {
                return fmt.Errorf("position with ID %s not found in portfolio %s", positionID, portfolioID)
        }

        portfolio.Positions = newPositions
        e.positions[portfolioID] = newPositions

        // Invalidate cache
        delete(e.performanceCache, portfolioID)
        delete(e.riskCache, portfolioID)

        return nil
}

// CalculatePerformanceMetrics calculates performance metrics for a portfolio
func (e *PortfolioAnalyticsEngine) CalculatePerformanceMetrics(portfolioID string) (*PerformanceMetrics, error) {
        e.mutex.RLock()
        defer e.mutex.RUnlock()

        // Check cache first
        if metrics, exists := e.performanceCache[portfolioID]; exists && metrics.UpdatedAt.After(time.Now().Add(-1*time.Hour)) {
                return metrics, nil
        }

        return e.calculatePerformanceMetrics(portfolioID)
}

// calculatePerformanceMetrics is the internal implementation of CalculatePerformanceMetrics
func (e *PortfolioAnalyticsEngine) calculatePerformanceMetrics(portfolioID string) (*PerformanceMetrics, error) {
        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return nil, fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        positions := e.positions[portfolioID]
        if len(positions) == 0 {
                return &PerformanceMetrics{
                        UpdatedAt: time.Now(),
                }, nil
        }

        // Calculate performance metrics
        // This is a simplified implementation
        var totalPnL, realizedPnL, unrealizedPnL float64
        var totalInvestment float64
        var winCount, lossCount int
        var totalWin, totalLoss float64

        for _, position := range positions {
                investment := float64(position.Quantity) * position.EntryPrice
                totalInvestment += investment

                if position.ExitTime != nil && position.ExitPrice != nil {
                        // Closed position
                        pnl := float64(position.Quantity) * (*position.ExitPrice - position.EntryPrice)
                        if position.TransactionType == "SELL" {
                                pnl = -pnl
                        }
                        realizedPnL += pnl

                        if pnl > 0 {
                                winCount++
                                totalWin += pnl
                        } else {
                                lossCount++
                                totalLoss += math.Abs(pnl)
                        }
                } else {
                        // Open position
                        pnl := float64(position.Quantity) * (position.CurrentPrice - position.EntryPrice)
                        if position.TransactionType == "SELL" {
                                pnl = -pnl
                        }
                        unrealizedPnL += pnl
                }
        }

        totalPnL = realizedPnL + unrealizedPnL
        pnlPercentage := 0.0
        if totalInvestment > 0 {
                pnlPercentage = totalPnL / totalInvestment * 100
        }

        winRate := 0.0
        if winCount+lossCount > 0 {
                winRate = float64(winCount) / float64(winCount+lossCount) * 100
        }

        averageWin := 0.0
        if winCount > 0 {
                averageWin = totalWin / float64(winCount)
        }

        averageLoss := 0.0
        if lossCount > 0 {
                averageLoss = totalLoss / float64(lossCount)
        }

        profitFactor := 0.0
        if totalLoss > 0 {
                profitFactor = totalWin / totalLoss
        }

        // Create performance metrics
        metrics := &PerformanceMetrics{
                TotalPnL:        totalPnL,
                RealizedPnL:     realizedPnL,
                UnrealizedPnL:   unrealizedPnL,
                PnLPercentage:   pnlPercentage,
                WinRate:         winRate,
                ProfitFactor:    profitFactor,
                AverageWin:      averageWin,
                AverageLoss:     averageLoss,
                ReturnOnCapital: pnlPercentage,
                UpdatedAt:       time.Now(),
        }

        // Cache the metrics
        e.performanceCache[portfolioID] = metrics
        portfolio.PerformanceCache = metrics

        return metrics, nil
}

// CalculateRiskMetrics calculates risk metrics for a portfolio
func (e *PortfolioAnalyticsEngine) CalculateRiskMetrics(portfolioID string) (*RiskMetrics, error) {
        e.mutex.RLock()
        defer e.mutex.RUnlock()

        // Check cache first
        if metrics, exists := e.riskCache[portfolioID]; exists && metrics.UpdatedAt.After(time.Now().Add(-1*time.Hour)) {
                return metrics, nil
        }

        return e.calculateRiskMetrics(portfolioID)
}

// calculateRiskMetrics is the internal implementation of CalculateRiskMetrics
func (e *PortfolioAnalyticsEngine) calculateRiskMetrics(portfolioID string) (*RiskMetrics, error) {
        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return nil, fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        positions := e.positions[portfolioID]
        if len(positions) == 0 {
                return &RiskMetrics{
                        UpdatedAt: time.Now(),
                }, nil
        }

        // Calculate risk metrics
        // This is a simplified implementation
        var deltaExposure, gammaExposure, thetaExposure, vegaExposure, rhoExposure float64
        sectorExposure := make(map[string]float64)
        assetClassExposure := make(map[string]float64)
        optionExposure := make(map[string]float64)

        for _, position := range positions {
                if position.ExitTime != nil {
                        // Skip closed positions
                        continue
                }

                value := float64(position.Quantity) * position.CurrentPrice
                
                // Asset class exposure
                assetClass := "Equity" // Default
                if position.OptionType != nil {
                        assetClass = "Options"
                        optionType := *position.OptionType
                        optionExposure[optionType] += value

                        // Greeks exposure
                        if position.Greeks != nil {
                                deltaExposure += position.Greeks.Delta * value
                                gammaExposure += position.Greeks.Gamma * value
                                thetaExposure += position.Greeks.Theta * value
                                vegaExposure += position.Greeks.Vega * value
                                rhoExposure += position.Greeks.Rho * value
                        }
                }
                assetClassExposure[assetClass] += value

                // Sector exposure (simplified)
                sector := "Unknown" // In a real implementation, this would be determined from a sector database
                sectorExposure[sector] += value
        }

        // Create risk metrics
        metrics := &RiskMetrics{
                SectorExposure:     sectorExposure,
                AssetClassExposure: assetClassExposure,
                OptionExposure:     optionExposure,
                DeltaExposure:      deltaExposure,
                GammaExposure:      gammaExposure,
                ThetaExposure:      thetaExposure,
                VegaExposure:       vegaExposure,
                RhoExposure:        rhoExposure,
                UpdatedAt:          time.Now(),
        }

        // Cache the metrics
        e.riskCache[portfolioID] = metrics
        portfolio.RiskCache = metrics

        return metrics, nil
}

// updatePositionPrices updates current prices for all positions in a portfolio
func (e *PortfolioAnalyticsEngine) updatePositionPrices(portfolioID string) error {
        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        positions := e.positions[portfolioID]
        ctx := context.Background()

        for i, position := range positions {
                if position.ExitTime != nil {
                        // Skip closed positions
                        continue
                }

                price, err := e.dataProvider.GetCurrentPrice(ctx, position.Symbol, position.Exchange)
                if err != nil {
                        return fmt.Errorf("failed to get current price for %s: %w", position.Symbol, err)
                }

                positions[i].CurrentPrice = price
        }

        portfolio.Positions = positions
        e.positions[portfolioID] = positions

        // Invalidate cache
        delete(e.performanceCache, portfolioID)
        delete(e.riskCache, portfolioID)

        return nil
}

// updatePositionGreeks updates Greeks for option positions in a portfolio
func (e *PortfolioAnalyticsEngine) updatePositionGreeks(portfolioID string) error {
        portfolio, exists := e.portfolios[portfolioID]
        if !exists {
                return fmt.Errorf("portfolio with ID %s not found", portfolioID)
        }

        positions := e.positions[portfolioID]
        ctx := context.Background()

        for i, position := range positions {
                if position.ExitTime != nil {
                        // Skip closed positions
                        continue
                }

                if position.OptionType == nil || position.StrikePrice == nil || position.ExpiryDate == nil {
                        // Skip non-option positions
                        continue
                }

                greeks, err := e.dataProvider.GetGreeks(ctx, position.Symbol, position.Exchange, *position.StrikePrice, *position.ExpiryDate, *position.OptionType)
                if err != nil {
                        return fmt.Errorf("failed to get Greeks for %s: %w", position.Symbol, err)
                }

                positions[i].Greeks = greeks
        }

        portfolio.Positions = positions
        e.positions[portfolioID] = positions

        // Invalidate cache
        delete(e.riskCache, portfolioID)

        return nil
}

// QueueTask queues a task for asynchronous processing
func (e *PortfolioAnalyticsEngine) QueueTask(taskType string, portfolioID string, callback func(interface{}, error)) error {
        e.mutex.RLock()
        defer e.mutex.RUnlock()

        if !e.isRunning {
                return errors.New("analytics engine is not running")
        }

        task := &AnalyticsTask{
                TaskType:    taskType,
                PortfolioID: portfolioID,
                Callback:    callback,
        }

        select {
        case e.calculationQueue <- task:
                return nil
        default:
                return errors.New("calculation queue is full")
        }
}
