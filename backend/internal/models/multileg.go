package models

import (
	"errors"
	"time"
)

// LegType defines the type of leg in a multileg strategy
type LegType string

const (
	LegTypeBuy       LegType = "BUY"
	LegTypeSell      LegType = "SELL"
	LegTypeBuyToOpen LegType = "BUY_TO_OPEN"
	LegTypeSellToOpen LegType = "SELL_TO_OPEN"
	LegTypeBuyToClose LegType = "BUY_TO_CLOSE"
	LegTypeSellToClose LegType = "SELL_TO_CLOSE"
)

// LegStatus defines the status of a leg
type LegStatus string

const (
	LegStatusPending  LegStatus = "PENDING"
	LegStatusActive   LegStatus = "ACTIVE"
	LegStatusExecuted LegStatus = "EXECUTED"
	LegStatusCanceled LegStatus = "CANCELED"
	LegStatusFailed   LegStatus = "FAILED"
)

// ExecutionType defines how a leg should be executed
type ExecutionType string

const (
	ExecutionTypeMarket     ExecutionType = "MARKET"
	ExecutionTypeLimit      ExecutionType = "LIMIT"
	ExecutionTypeStop       ExecutionType = "STOP"
	ExecutionTypeStopLimit  ExecutionType = "STOP_LIMIT"
	ExecutionTypeTrailing   ExecutionType = "TRAILING"
)

// HedgeType defines the type of hedging to apply
type HedgeType string

const (
	HedgeTypeNone     HedgeType = "NONE"
	HedgeTypeDelta    HedgeType = "DELTA"
	HedgeTypeGamma    HedgeType = "GAMMA"
	HedgeTypeVega     HedgeType = "VEGA"
	HedgeTypeTheta    HedgeType = "THETA"
	HedgeTypeDynamic  HedgeType = "DYNAMIC"
)

// MultilegStrategy represents a complex trading strategy with multiple legs
type MultilegStrategy struct {
	ID              string           `json:"id" bson:"_id,omitempty"`
	Name            string           `json:"name" bson:"name"`
	Description     string           `json:"description" bson:"description"`
	UserID          string           `json:"userId" bson:"userId"`
	PortfolioID     string           `json:"portfolioId" bson:"portfolioId"`
	Legs            []Leg            `json:"legs" bson:"legs"`
	ExecutionParams ExecutionParams  `json:"executionParams" bson:"executionParams"`
	RiskParams      RiskParameters   `json:"riskParams" bson:"riskParams"`
	HedgeParams     HedgeParameters  `json:"hedgeParams" bson:"hedgeParams"`
	Status          string           `json:"status" bson:"status"`
	Tags            []string         `json:"tags" bson:"tags"`
	CreatedAt       time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt" bson:"updatedAt"`
	LastExecutedAt  time.Time        `json:"lastExecutedAt,omitempty" bson:"lastExecutedAt,omitempty"`
}

// Leg represents a single leg in a multileg strategy
type Leg struct {
	ID              string        `json:"id" bson:"_id,omitempty"`
	StrategyID      string        `json:"strategyId" bson:"strategyId"`
	Symbol          string        `json:"symbol" bson:"symbol"`
	Type            LegType       `json:"type" bson:"type"`
	Quantity        int           `json:"quantity" bson:"quantity"`
	ExecutionType   ExecutionType `json:"executionType" bson:"executionType"`
	Price           float64       `json:"price,omitempty" bson:"price,omitempty"`
	StopPrice       float64       `json:"stopPrice,omitempty" bson:"stopPrice,omitempty"`
	TrailingAmount  float64       `json:"trailingAmount,omitempty" bson:"trailingAmount,omitempty"`
	TrailingPercent float64       `json:"trailingPercent,omitempty" bson:"trailingPercent,omitempty"`
	Status          LegStatus     `json:"status" bson:"status"`
	OrderID         string        `json:"orderId,omitempty" bson:"orderId,omitempty"`
	ExecutionTime   time.Time     `json:"executionTime,omitempty" bson:"executionTime,omitempty"`
	ExecutedPrice   float64       `json:"executedPrice,omitempty" bson:"executedPrice,omitempty"`
	Sequence        int           `json:"sequence" bson:"sequence"`
	DependsOn       []string      `json:"dependsOn,omitempty" bson:"dependsOn,omitempty"`
	CreatedAt       time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// ExecutionParams defines parameters for strategy execution
type ExecutionParams struct {
	Sequential      bool      `json:"sequential" bson:"sequential"`
	SimultaneousLegs bool     `json:"simultaneousLegs" bson:"simultaneousLegs"`
	TimeWindow      int       `json:"timeWindow" bson:"timeWindow"` // in seconds
	MaxSlippage     float64   `json:"maxSlippage" bson:"maxSlippage"` // in percentage
	EntryConditions []Condition `json:"entryConditions,omitempty" bson:"entryConditions,omitempty"`
	ExitConditions  []Condition `json:"exitConditions,omitempty" bson:"exitConditions,omitempty"`
	RangeBreakout   RangeBreakout `json:"rangeBreakout,omitempty" bson:"rangeBreakout,omitempty"`
}

// RangeBreakout defines parameters for range breakout execution
type RangeBreakout struct {
	Enabled      bool    `json:"enabled" bson:"enabled"`
	UpperBound   float64 `json:"upperBound" bson:"upperBound"`
	LowerBound   float64 `json:"lowerBound" bson:"lowerBound"`
	Confirmation int     `json:"confirmation" bson:"confirmation"` // number of ticks to confirm breakout
	Symbol       string  `json:"symbol" bson:"symbol"`
}

// HedgeParameters defines parameters for strategy hedging
type HedgeParameters struct {
	Type           HedgeType `json:"type" bson:"type"`
	Instrument     string    `json:"instrument,omitempty" bson:"instrument,omitempty"`
	Ratio          float64   `json:"ratio,omitempty" bson:"ratio,omitempty"`
	RebalanceFreq  int       `json:"rebalanceFreq,omitempty" bson:"rebalanceFreq,omitempty"` // in minutes
	DynamicThreshold float64 `json:"dynamicThreshold,omitempty" bson:"dynamicThreshold,omitempty"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
}

// Validate validates the multileg strategy
func (s *MultilegStrategy) Validate() error {
	if s.Name == "" {
		return errors.New("strategy name is required")
	}
	if s.UserID == "" {
		return errors.New("user ID is required")
	}
	if s.PortfolioID == "" {
		return errors.New("portfolio ID is required")
	}
	if len(s.Legs) == 0 {
		return errors.New("at least one leg is required")
	}
	
	// Validate each leg
	for _, leg := range s.Legs {
		if err := leg.Validate(); err != nil {
			return err
		}
	}
	
	// Validate execution parameters
	if s.ExecutionParams.MaxSlippage < 0 {
		return errors.New("max slippage cannot be negative")
	}
	
	// Validate risk parameters
	if s.RiskParams.MaxLoss <= 0 {
		return errors.New("max loss must be greater than zero")
	}
	
	return nil
}

// Validate validates the leg
func (l *Leg) Validate() error {
	if l.Symbol == "" {
		return errors.New("symbol is required")
	}
	if l.Type == "" {
		return errors.New("leg type is required")
	}
	if l.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if l.ExecutionType == "" {
		return errors.New("execution type is required")
	}
	
	// Validate price based on execution type
	switch l.ExecutionType {
	case ExecutionTypeLimit, ExecutionTypeStopLimit:
		if l.Price <= 0 {
			return errors.New("price must be greater than zero for limit orders")
		}
	case ExecutionTypeStop, ExecutionTypeStopLimit:
		if l.StopPrice <= 0 {
			return errors.New("stop price must be greater than zero for stop orders")
		}
	case ExecutionTypeTrailing:
		if l.TrailingAmount <= 0 && l.TrailingPercent <= 0 {
			return errors.New("trailing amount or percent must be greater than zero for trailing orders")
		}
	}
	
	return nil
}
