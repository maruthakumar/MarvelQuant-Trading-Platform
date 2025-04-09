package models

import (
	"time"
)

// SimulationAccount represents a simulation account for backtesting and paper trading
type SimulationAccount struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"userId" db:"user_id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	InitialBalance  float64   `json:"initialBalance" db:"initial_balance"`
	CurrentBalance  float64   `json:"currentBalance" db:"current_balance"`
	Currency        string    `json:"currency" db:"currency"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
	IsActive        bool      `json:"isActive" db:"is_active"`
	SimulationType  string    `json:"simulationType" db:"simulation_type"` // "PAPER" or "BACKTEST"
	RiskSettings    *RiskSettings `json:"riskSettings" db:"risk_settings"`
	MarketSettings  *MarketSettings `json:"marketSettings" db:"market_settings"`
}

// RiskSettings represents risk management settings for a simulation account
type RiskSettings struct {
	MaxPositionSize       float64 `json:"maxPositionSize" db:"max_position_size"`
	MaxDrawdown           float64 `json:"maxDrawdown" db:"max_drawdown"`
	MaxDailyLoss          float64 `json:"maxDailyLoss" db:"max_daily_loss"`
	MaxOpenPositions      int     `json:"maxOpenPositions" db:"max_open_positions"`
	MaxLeverage           float64 `json:"maxLeverage" db:"max_leverage"`
	StopLossRequired      bool    `json:"stopLossRequired" db:"stop_loss_required"`
	TakeProfitRecommended bool    `json:"takeProfitRecommended" db:"take_profit_recommended"`
}

// MarketSettings represents market simulation settings for a simulation account
type MarketSettings struct {
	SlippageModel      string  `json:"slippageModel" db:"slippage_model"` // "FIXED", "PERCENTAGE", "VARIABLE"
	SlippageValue      float64 `json:"slippageValue" db:"slippage_value"`
	LatencyModel       string  `json:"latencyModel" db:"latency_model"` // "FIXED", "VARIABLE", "REALISTIC"
	LatencyValue       int     `json:"latencyValue" db:"latency_value"` // in milliseconds
	PriceFeedSource    string  `json:"priceFeedSource" db:"price_feed_source"`
	CommissionModel    string  `json:"commissionModel" db:"commission_model"` // "FIXED", "PERCENTAGE", "TIERED"
	CommissionValue    float64 `json:"commissionValue" db:"commission_value"`
	SpreadModel        string  `json:"spreadModel" db:"spread_model"` // "FIXED", "VARIABLE", "REALISTIC"
	SpreadValue        float64 `json:"spreadValue" db:"spread_value"`
	AllowShortSelling  bool    `json:"allowShortSelling" db:"allow_short_selling"`
	AllowFractionalLots bool   `json:"allowFractionalLots" db:"allow_fractional_lots"`
}

// SimulationTransaction represents a transaction in a simulation account
type SimulationTransaction struct {
	ID                 string    `json:"id" db:"id"`
	SimulationAccountID string   `json:"simulationAccountId" db:"simulation_account_id"`
	Type               string    `json:"type" db:"type"` // "DEPOSIT", "WITHDRAWAL", "FEE", "INTEREST", "DIVIDEND", "P&L"
	Amount             float64   `json:"amount" db:"amount"`
	Balance            float64   `json:"balance" db:"balance"` // Balance after transaction
	Description        string    `json:"description" db:"description"`
	ReferenceID        string    `json:"referenceId" db:"reference_id"` // ID of related entity (order, position)
	ReferenceType      string    `json:"referenceType" db:"reference_type"` // Type of related entity
	Timestamp          time.Time `json:"timestamp" db:"timestamp"`
}

// SimulationOrder represents an order in the simulation system
type SimulationOrder struct {
	Order              // Embed the base Order struct
	SimulationAccountID string    `json:"simulationAccountId" db:"simulation_account_id"`
	SimulatedFillPrice float64    `json:"simulatedFillPrice" db:"simulated_fill_price"`
	SimulatedFillTime  time.Time  `json:"simulatedFillTime" db:"simulated_fill_time"`
	SlippageAmount     float64    `json:"slippageAmount" db:"slippage_amount"`
	LatencyMs          int        `json:"latencyMs" db:"latency_ms"`
	CommissionAmount   float64    `json:"commissionAmount" db:"commission_amount"`
	IsBacktestOrder    bool       `json:"isBacktestOrder" db:"is_backtest_order"`
	BacktestDate       *time.Time `json:"backtestDate" db:"backtest_date"`
}

// SimulationPosition represents a position in the simulation system
type SimulationPosition struct {
	Position              // Embed the base Position struct
	SimulationAccountID string    `json:"simulationAccountId" db:"simulation_account_id"`
	SimulatedEntryPrice float64   `json:"simulatedEntryPrice" db:"simulated_entry_price"`
	SimulatedMarketPrice float64  `json:"simulatedMarketPrice" db:"simulated_market_price"`
	TotalCommission     float64   `json:"totalCommission" db:"total_commission"`
	TotalSlippage       float64   `json:"totalSlippage" db:"total_slippage"`
	IsBacktestPosition  bool      `json:"isBacktestPosition" db:"is_backtest_position"`
	BacktestDate        *time.Time `json:"backtestDate" db:"backtest_date"`
}

// BacktestSession represents a backtesting session
type BacktestSession struct {
	ID                 string    `json:"id" db:"id"`
	SimulationAccountID string   `json:"simulationAccountId" db:"simulation_account_id"`
	Name               string    `json:"name" db:"name"`
	Description        string    `json:"description" db:"description"`
	StartDate          time.Time `json:"startDate" db:"start_date"`
	EndDate            time.Time `json:"endDate" db:"end_date"`
	Symbols            []string  `json:"symbols" db:"symbols"`
	Timeframe          string    `json:"timeframe" db:"timeframe"` // "1m", "5m", "15m", "1h", "1d", etc.
	InitialBalance     float64   `json:"initialBalance" db:"initial_balance"`
	FinalBalance       float64   `json:"finalBalance" db:"final_balance"`
	TotalTrades        int       `json:"totalTrades" db:"total_trades"`
	WinningTrades      int       `json:"winningTrades" db:"winning_trades"`
	LosingTrades       int       `json:"losingTrades" db:"losing_trades"`
	ProfitFactor       float64   `json:"profitFactor" db:"profit_factor"`
	SharpeRatio        float64   `json:"sharpeRatio" db:"sharpe_ratio"`
	MaxDrawdown        float64   `json:"maxDrawdown" db:"max_drawdown"`
	AnnualizedReturn   float64   `json:"annualizedReturn" db:"annualized_return"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	CompletedAt        *time.Time `json:"completedAt" db:"completed_at"`
	Status             string    `json:"status" db:"status"` // "PENDING", "RUNNING", "COMPLETED", "FAILED"
	StrategyID         string    `json:"strategyId" db:"strategy_id"`
	Parameters         map[string]interface{} `json:"parameters" db:"parameters"`
}

// BacktestResult represents a single result point in a backtest
type BacktestResult struct {
	ID               string    `json:"id" db:"id"`
	BacktestSessionID string   `json:"backtestSessionId" db:"backtest_session_id"`
	Timestamp        time.Time `json:"timestamp" db:"timestamp"`
	EquityCurve      float64   `json:"equityCurve" db:"equity_curve"`
	DrawdownCurve    float64   `json:"drawdownCurve" db:"drawdown_curve"`
	OpenPositions    int       `json:"openPositions" db:"open_positions"`
	CumulativePnL    float64   `json:"cumulativePnL" db:"cumulative_pnl"`
	DailyPnL         float64   `json:"dailyPnL" db:"daily_pnl"`
	MarketValue      float64   `json:"marketValue" db:"market_value"`
	CashBalance      float64   `json:"cashBalance" db:"cash_balance"`
}

// MarketDataSnapshot represents a snapshot of market data for simulation
type MarketDataSnapshot struct {
	ID          string    `json:"id" db:"id"`
	Symbol      string    `json:"symbol" db:"symbol"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Open        float64   `json:"open" db:"open"`
	High        float64   `json:"high" db:"high"`
	Low         float64   `json:"low" db:"low"`
	Close       float64   `json:"close" db:"close"`
	Volume      int64     `json:"volume" db:"volume"`
	Bid         float64   `json:"bid" db:"bid"`
	Ask         float64   `json:"ask" db:"ask"`
	BidSize     int       `json:"bidSize" db:"bid_size"`
	AskSize     int       `json:"askSize" db:"ask_size"`
	Timeframe   string    `json:"timeframe" db:"timeframe"` // "1m", "5m", "15m", "1h", "1d", etc.
	Source      string    `json:"source" db:"source"`
	IsSimulated bool      `json:"isSimulated" db:"is_simulated"`
}
