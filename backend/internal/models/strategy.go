package models

import (
	"errors"
	"time"
)

// StrategyType defines the type of strategy
type StrategyType string

const (
	StrategyTypeManual    StrategyType = "MANUAL"
	StrategyTypeAutomated StrategyType = "AUTOMATED"
	StrategyTypeAlgo      StrategyType = "ALGO"
)

// StrategyStatus defines the status of a strategy
type StrategyStatus string

const (
	StrategyStatusDraft   StrategyStatus = "DRAFT"
	StrategyStatusActive  StrategyStatus = "ACTIVE"
	StrategyStatusPaused  StrategyStatus = "PAUSED"
	StrategyStatusStopped StrategyStatus = "STOPPED"
	StrategyStatusFailed  StrategyStatus = "FAILED"
	StrategyStatusArchived StrategyStatus = "ARCHIVED"
)

// ScheduleFrequency defines the frequency of strategy execution
type ScheduleFrequency string

const (
	ScheduleFrequencyOnce      ScheduleFrequency = "ONCE"
	ScheduleFrequencyDaily     ScheduleFrequency = "DAILY"
	ScheduleFrequencyWeekly    ScheduleFrequency = "WEEKLY"
	ScheduleFrequencyMonthly   ScheduleFrequency = "MONTHLY"
	ScheduleFrequencyCustom    ScheduleFrequency = "CUSTOM"
)

// Strategy represents a trading strategy
type Strategy struct {
	ID              string         `json:"id" bson:"_id,omitempty"`
	Name            string         `json:"name" bson:"name"`
	Description     string         `json:"description" bson:"description"`
	UserID          string         `json:"userId" bson:"userId"`
	Type            StrategyType   `json:"type" bson:"type"`
	Status          StrategyStatus `json:"status" bson:"status"`
	EntryConditions []Condition    `json:"entryConditions" bson:"entryConditions"`
	ExitConditions  []Condition    `json:"exitConditions" bson:"exitConditions"`
	RiskParameters  RiskParameters `json:"riskParameters" bson:"riskParameters"`
	Instruments     []string       `json:"instruments" bson:"instruments"`
	Tags            []string       `json:"tags" bson:"tags"`
	CreatedAt       time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt" bson:"updatedAt"`
	LastExecutedAt  time.Time      `json:"lastExecutedAt,omitempty" bson:"lastExecutedAt,omitempty"`
}

// Condition represents a trading condition
type Condition struct {
	Type      string      `json:"type" bson:"type"`
	Parameter string      `json:"parameter" bson:"parameter"`
	Operator  string      `json:"operator" bson:"operator"`
	Value     interface{} `json:"value" bson:"value"`
}

// RiskParameters represents risk management parameters
type RiskParameters struct {
	MaxPositionSize     float64 `json:"maxPositionSize" bson:"maxPositionSize"`
	MaxLoss             float64 `json:"maxLoss" bson:"maxLoss"`
	MaxDailyLoss        float64 `json:"maxDailyLoss" bson:"maxDailyLoss"`
	TrailingStopPercent float64 `json:"trailingStopPercent" bson:"trailingStopPercent"`
	TakeProfitPercent   float64 `json:"takeProfitPercent" bson:"takeProfitPercent"`
}

// StrategySchedule represents a schedule for strategy execution
type StrategySchedule struct {
	StrategyID string            `json:"strategyId" bson:"strategyId"`
	Frequency  ScheduleFrequency `json:"frequency" bson:"frequency"`
	StartTime  time.Time         `json:"startTime" bson:"startTime"`
	EndTime    time.Time         `json:"endTime,omitempty" bson:"endTime,omitempty"`
	DaysOfWeek []int             `json:"daysOfWeek,omitempty" bson:"daysOfWeek,omitempty"` // 0 = Sunday, 1 = Monday, etc.
	Enabled    bool              `json:"enabled" bson:"enabled"`
	CreatedAt  time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt" bson:"updatedAt"`
}

// StrategyPerformance represents performance metrics for a strategy
type StrategyPerformance struct {
	StrategyID    string    `json:"strategyId" bson:"strategyId"`
	TotalPnL      float64   `json:"totalPnL" bson:"totalPnL"`
	WinCount      int       `json:"winCount" bson:"winCount"`
	LossCount     int       `json:"lossCount" bson:"lossCount"`
	TotalTrades   int       `json:"totalTrades" bson:"totalTrades"`
	WinRate       float64   `json:"winRate" bson:"winRate"`
	MaxDrawdown   float64   `json:"maxDrawdown" bson:"maxDrawdown"`
	OrderCount    int       `json:"orderCount" bson:"orderCount"`
	PositionCount int       `json:"positionCount" bson:"positionCount"`
	StartDate     time.Time `json:"startDate" bson:"startDate"`
	EndDate       time.Time `json:"endDate" bson:"endDate"`
}

// Validate validates the strategy
func (s *Strategy) Validate() error {
	if s.Name == "" {
		return errors.New("strategy name is required")
	}
	if s.UserID == "" {
		return errors.New("user ID is required")
	}
	if s.Type == "" {
		return errors.New("strategy type is required")
	}
	if len(s.Instruments) == 0 {
		return errors.New("at least one instrument is required")
	}
	
	// Validate risk parameters
	if s.RiskParameters.MaxPositionSize <= 0 {
		return errors.New("max position size must be greater than zero")
	}
	
	return nil
}

// Validate validates the strategy schedule
func (s *StrategySchedule) Validate() error {
	if s.StrategyID == "" {
		return errors.New("strategy ID is required")
	}
	if s.Frequency == "" {
		return errors.New("schedule frequency is required")
	}
	if s.StartTime.IsZero() {
		return errors.New("start time is required")
	}
	
	// Validate days of week for weekly frequency
	if s.Frequency == ScheduleFrequencyWeekly && len(s.DaysOfWeek) == 0 {
		return errors.New("days of week are required for weekly frequency")
	}
	
	// Validate that days of week are valid (0-6)
	for _, day := range s.DaysOfWeek {
		if day < 0 || day > 6 {
			return errors.New("days of week must be between 0 and 6")
		}
	}
	
	return nil
}
