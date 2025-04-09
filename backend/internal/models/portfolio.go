package models

import (
        "errors"
        "regexp"
        "time"
)

// PortfolioStatus represents the current status of a portfolio
type PortfolioStatus string

const (
        PortfolioStatusPending   PortfolioStatus = "PENDING"
        PortfolioStatusActive    PortfolioStatus = "ACTIVE"
        PortfolioStatusCompleted PortfolioStatus = "COMPLETED"
        PortfolioStatusFailed    PortfolioStatus = "FAILED"
)

// StrikeSelectionMode represents the mode for selecting option strikes
type StrikeSelectionMode string

const (
        StrikeSelectionModeNormal   StrikeSelectionMode = "NORMAL"
        StrikeSelectionModeRelative StrikeSelectionMode = "RELATIVE"
        StrikeSelectionModeBoth     StrikeSelectionMode = "BOTH"
)

// UnderlyingReference represents the reference for underlying price
type UnderlyingReference string

const (
        UnderlyingReferenceFuture UnderlyingReference = "FUTURE"
        UnderlyingReferenceSpot   UnderlyingReference = "SPOT"
)

// PriceType represents the type of price to use for calculations
type PriceType string

const (
        PriceTypeLTP      PriceType = "LTP"
        PriceTypeBidAsk   PriceType = "BID_ASK"
        PriceTypeBidAskAvg PriceType = "BID_ASK_AVG"
)

// Portfolio represents a multi-leg options portfolio in the system
type Portfolio struct {
        ID                 string            `json:"id" bson:"_id,omitempty"`
        UserID             string            `json:"userId" bson:"userId"`
        Name               string            `json:"name" bson:"name"`
        StrategyID         string            `json:"strategyId" bson:"strategyId"`
        Status             PortfolioStatus   `json:"status" bson:"status"`
        
        // Default Portfolio Settings
        Exchange           string            `json:"exchange" bson:"exchange"`
        Symbol             string            `json:"symbol" bson:"symbol"`
        Expiry             time.Time         `json:"expiry" bson:"expiry"`
        DefaultLots        int               `json:"defaultLots" bson:"defaultLots"`
        PredefinedStrategy string            `json:"predefinedStrategy" bson:"predefinedStrategy"`
        StrikeSelection    StrikeSelectionMode `json:"strikeSelection" bson:"strikeSelection"`
        UnderlyingRef      UnderlyingReference `json:"underlyingRef" bson:"underlyingRef"`
        PriceType          PriceType         `json:"priceType" bson:"priceType"`
        StrikeStep         float64           `json:"strikeStep" bson:"strikeStep"`
        IsPositional       bool              `json:"isPositional" bson:"isPositional"`
        BuyTradesFirst     bool              `json:"buyTradesFirst" bson:"buyTradesFirst"`
        AllowFarStrikes    bool              `json:"allowFarStrikes" bson:"allowFarStrikes"`
        ImpliedSynthetic   bool              `json:"impliedSynthetic" bson:"impliedSynthetic"`
        ValueAllLots       float64           `json:"valueAllLots" bson:"valueAllLots"`
        ValuePerLot        float64           `json:"valuePerLot" bson:"valuePerLot"`
        
        // Execution Parameters
        ProductType        ProductType       `json:"productType" bson:"productType"`
        FailureAction      FailureAction     `json:"failureAction" bson:"failureAction"`
        LegExecutionMode   LegExecutionMode  `json:"legExecutionMode" bson:"legExecutionMode"`
        QuantityByExposure float64           `json:"quantityByExposure,omitempty" bson:"quantityByExposure,omitempty"`
        MaxLots            int               `json:"maxLots" bson:"maxLots"`
        PremiumGap         float64           `json:"premiumGap,omitempty" bson:"premiumGap,omitempty"`
        RunOnDays          []string          `json:"runOnDays" bson:"runOnDays"`
        StartTime          string            `json:"startTime" bson:"startTime"`
        EndTime            string            `json:"endTime" bson:"endTime"`
        SquareOffTime      string            `json:"squareOffTime" bson:"squareOffTime"`
        ExecutionMode      ExecutionMode     `json:"executionMode" bson:"executionMode"`
        EntryOrderType     OrderType         `json:"entryOrderType" bson:"entryOrderType"`
        EstimatedMargin    float64           `json:"estimatedMargin" bson:"estimatedMargin"`
        
        // Range Breakout Settings
        RangeBreakoutEnabled bool              `json:"rangeBreakoutEnabled" bson:"rangeBreakoutEnabled"`
        RangeStartTime      string            `json:"rangeStartTime,omitempty" bson:"rangeStartTime,omitempty"`
        RangeEndTime        string            `json:"rangeEndTime,omitempty" bson:"rangeEndTime,omitempty"`
        HighBuffer          float64           `json:"highBuffer,omitempty" bson:"highBuffer,omitempty"`
        LowBuffer           float64           `json:"lowBuffer,omitempty" bson:"lowBuffer,omitempty"`
        OppositeSideSL      bool              `json:"oppositeSideSL,omitempty" bson:"oppositeSideSL,omitempty"`
        RangeBuffer         float64           `json:"rangeBuffer,omitempty" bson:"rangeBuffer,omitempty"`
        BreakoutNotReqLegs  []int             `json:"breakoutNotReqLegs,omitempty" bson:"breakoutNotReqLegs,omitempty"`
        
        // Extra Conditions
        GapUpMinimum       float64           `json:"gapUpMinimum,omitempty" bson:"gapUpMinimum,omitempty"`
        GapUpMaximum       float64           `json:"gapUpMaximum,omitempty" bson:"gapUpMaximum,omitempty"`
        GapDownMinimum     float64           `json:"gapDownMinimum,omitempty" bson:"gapDownMinimum,omitempty"`
        GapDownMaximum     float64           `json:"gapDownMaximum,omitempty" bson:"gapDownMaximum,omitempty"`
        PreviousDayReference string           `json:"previousDayReference,omitempty" bson:"previousDayReference,omitempty"`
        CombinedWaitAndTrade float64          `json:"combinedWaitAndTrade,omitempty" bson:"combinedWaitAndTrade,omitempty"`
        LegEntryConditions map[int]string    `json:"legEntryConditions,omitempty" bson:"legEntryConditions,omitempty"`
        
        // Other Settings
        KeepAllUsersInSync bool              `json:"keepAllUsersInSync" bson:"keepAllUsersInSync"`
        TrailWaitAndTrade  bool              `json:"trailWaitAndTrade" bson:"trailWaitAndTrade"`
        ExecuteDelay       int               `json:"executeDelay" bson:"executeDelay"`
        ReExecuteDelay     int               `json:"reExecuteDelay" bson:"reExecuteDelay"`
        StraddleWidthMultiplier float64      `json:"straddleWidthMultiplier,omitempty" bson:"straddleWidthMultiplier,omitempty"`
        OnActionTrigger    string            `json:"onActionTrigger" bson:"onActionTrigger"`
        
        // Monitoring Settings
        PositionalTimes    map[string]string `json:"positionalTimes,omitempty" bson:"positionalTimes,omitempty"`
        LegMonitoringType  MonitoringType    `json:"legMonitoringType" bson:"legMonitoringType"`
        CombinedMonitoringType MonitoringType `json:"combinedMonitoringType" bson:"combinedMonitoringType"`
        MonitoringInterval int               `json:"monitoringInterval,omitempty" bson:"monitoringInterval,omitempty"`
        
        // Dynamic Hedge Settings
        MinHedgeDistance   int               `json:"minHedgeDistance,omitempty" bson:"minHedgeDistance,omitempty"`
        MaxHedgeDistance   int               `json:"maxHedgeDistance,omitempty" bson:"maxHedgeDistance,omitempty"`
        MaxHedgePremium    float64           `json:"maxHedgePremium,omitempty" bson:"maxHedgePremium,omitempty"`
        MinHedgeOI         int               `json:"minHedgeOI,omitempty" bson:"minHedgeOI,omitempty"`
        UnsatisfiedHedgeAction string         `json:"unsatisfiedHedgeAction,omitempty" bson:"unsatisfiedHedgeAction,omitempty"`
        DeltaTarget        float64           `json:"deltaTarget,omitempty" bson:"deltaTarget,omitempty"`
        
        // Target Settings
        TargetType         TargetType        `json:"targetType" bson:"targetType"`
        TargetValue        float64           `json:"targetValue" bson:"targetValue"`
        OnTargetAction     string            `json:"onTargetAction" bson:"onTargetAction"`
        ProfitLockThreshold float64          `json:"profitLockThreshold,omitempty" bson:"profitLockThreshold,omitempty"`
        MinimumProfitLock  float64           `json:"minimumProfitLock,omitempty" bson:"minimumProfitLock,omitempty"`
        ProfitTrailAmount  float64           `json:"profitTrailAmount,omitempty" bson:"profitTrailAmount,omitempty"`
        ProfitTrailValue   float64           `json:"profitTrailValue,omitempty" bson:"profitTrailValue,omitempty"`
        
        // Stop Loss Settings
        StopLossType       StopLossType      `json:"stopLossType" bson:"stopLossType"`
        StopLossValue      float64           `json:"stopLossValue" bson:"stopLossValue"`
        StopLossWaitSeconds int              `json:"stopLossWaitSeconds" bson:"stopLossWaitSeconds"`
        OnStopLossAction   string            `json:"onStopLossAction" bson:"onStopLossAction"`
        StopLossTrailAmount float64          `json:"stopLossTrailAmount,omitempty" bson:"stopLossTrailAmount,omitempty"`
        StopLossTrailValue float64           `json:"stopLossTrailValue,omitempty" bson:"stopLossTrailValue,omitempty"`
        
        // Exit Settings
        ExitMode           ExitMode          `json:"exitMode" bson:"exitMode"`
        ExitOrderType      OrderType         `json:"exitOrderType" bson:"exitOrderType"`
        ExitPriceBuffer    float64           `json:"exitPriceBuffer" bson:"exitPriceBuffer"`
        MaxExitRetries     int               `json:"maxExitRetries" bson:"maxExitRetries"`
        ExitRetryInterval  int               `json:"exitRetryInterval" bson:"exitRetryInterval"`
        ConvertToMarket    bool              `json:"convertToMarket" bson:"convertToMarket"`
        EnablePartialExits bool              `json:"enablePartialExits" bson:"enablePartialExits"`
        MinExitPercentage  float64           `json:"minExitPercentage,omitempty" bson:"minExitPercentage,omitempty"`
        ExitSequence       string            `json:"exitSequence,omitempty" bson:"exitSequence,omitempty"`
        
        // Broker Settings
        BrokerSelection    string            `json:"brokerSelection" bson:"brokerSelection"`
        BrokerSpecificSettings map[string]interface{} `json:"brokerSpecificSettings,omitempty" bson:"brokerSpecificSettings,omitempty"`
        
        // Performance Metrics
        EntryValue         float64           `json:"entryValue" bson:"entryValue"`
        CurrentValue       float64           `json:"currentValue" bson:"currentValue"`
        MaxValue           float64           `json:"maxValue" bson:"maxValue"`
        MinValue           float64           `json:"minValue" bson:"minValue"`
        UnrealizedPnL      float64           `json:"unrealizedPnL" bson:"unrealizedPnL"`
        RealizedPnL        float64           `json:"realizedPnL" bson:"realizedPnL"`
        TotalPnL           float64           `json:"totalPnL" bson:"totalPnL"`
        PnLPercentage      float64           `json:"pnLPercentage" bson:"pnLPercentage"`
        Delta              float64           `json:"delta" bson:"delta"`
        Gamma              float64           `json:"gamma" bson:"gamma"`
        Theta              float64           `json:"theta" bson:"theta"`
        Vega               float64           `json:"vega" bson:"vega"`
        
        // Execution Tracking
        ExecutionStartTime time.Time         `json:"executionStartTime,omitempty" bson:"executionStartTime,omitempty"`
        ExecutionEndTime   time.Time         `json:"executionEndTime,omitempty" bson:"executionEndTime,omitempty"`
        LastMonitorTime    time.Time         `json:"lastMonitorTime,omitempty" bson:"lastMonitorTime,omitempty"`
        ExecutionLogs      []string          `json:"executionLogs,omitempty" bson:"executionLogs,omitempty"`
        
        Legs               []Leg             `json:"legs" bson:"legs"`
        CreatedAt          time.Time         `json:"createdAt" bson:"createdAt"`
        UpdatedAt          time.Time         `json:"updatedAt" bson:"updatedAt"`
}

// PortfolioFilter represents filters for querying portfolios
type PortfolioFilter struct {
        UserID       string          `json:"userId,omitempty"`
        Name         string          `json:"name,omitempty"`
        StrategyID   string          `json:"strategyId,omitempty"`
        Status       PortfolioStatus `json:"status,omitempty"`
        Symbol       string          `json:"symbol,omitempty"`
        Exchange     string          `json:"exchange,omitempty"`
        FromDate     time.Time       `json:"fromDate,omitempty"`
        ToDate       time.Time       `json:"toDate,omitempty"`
}

// Validate validates the portfolio data
func (p *Portfolio) Validate() error {
        // Check required fields
        if p.UserID == "" {
                return errors.New("user ID is required")
        }
        if p.Name == "" {
                return errors.New("portfolio name is required")
        }
        if p.Symbol == "" {
                return errors.New("symbol is required")
        }
        if p.Exchange == "" {
                return errors.New("exchange is required")
        }
        if p.Expiry.IsZero() {
                return errors.New("expiry date is required")
        }
        if p.DefaultLots <= 0 {
                return errors.New("default lots must be greater than zero")
        }

        // Validate portfolio status
        switch p.Status {
        case PortfolioStatusPending, PortfolioStatusActive, PortfolioStatusCompleted, PortfolioStatusFailed:
                // Valid statuses
        default:
                return errors.New("invalid portfolio status")
        }

        // Validate strike selection mode
        switch p.StrikeSelection {
        case StrikeSelectionModeNormal, StrikeSelectionModeRelative, StrikeSelectionModeBoth:
                // Valid strike selection modes
        default:
                return errors.New("invalid strike selection mode")
        }

        // Validate underlying reference
        switch p.UnderlyingRef {
        case UnderlyingReferenceFuture, UnderlyingReferenceSpot:
                // Valid underlying references
        default:
                return errors.New("invalid underlying reference")
        }

        // Validate price type
        switch p.PriceType {
        case PriceTypeLTP, PriceTypeBidAsk, PriceTypeBidAskAvg:
                // Valid price types
        default:
                return errors.New("invalid price type")
        }

        // Validate strike step
        if p.StrikeStep <= 0 {
                return errors.New("strike step must be greater than zero")
        }

        // Validate product type
        switch p.ProductType {
        case ProductTypeMIS, ProductTypeNRML, ProductTypeCNC:
                // Valid product types
        default:
                return errors.New("invalid product type")
        }

        // Validate failure action
        switch p.FailureAction {
        case FailureActionKeepPlacedLegs, FailureActionExitPlacedLegs:
                // Valid failure actions
        default:
                return errors.New("invalid failure action")
        }

        // Validate leg execution mode
        switch p.LegExecutionMode {
        case LegExecutionModeParallel, LegExecutionModeSequential:
                // Valid leg execution modes
        default:
                return errors.New("invalid leg execution mode")
        }

        // Validate max lots
        if p.MaxLots <= 0 {
                return errors.New("max lots must be greater than zero")
        }

        // Validate run on days
        validDays := map[string]bool{
                "MONDAY": true, "TUESDAY": true, "WEDNESDAY": true, 
                "THURSDAY": true, "FRIDAY": true, "SATURDAY": true, "SUNDAY": true,
        }
        for _, day := range p.RunOnDays {
                if !validDays[day] {
                        return errors.New("invalid day in run on days: " + day)
                }
        }

        // Validate time formats (HH:MM:SS)
        timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)
        if !timeRegex.MatchString(p.StartTime) {
                return errors.New("invalid start time format (use HH:MM:SS)")
        }
        if !timeRegex.MatchString(p.EndTime) {
                return errors.New("invalid end time format (use HH:MM:SS)")
        }
        if !timeRegex.MatchString(p.SquareOffTime) {
                return errors.New("invalid square off time format (use HH:MM:SS)")
        }

        // Validate execution mode
        switch p.ExecutionMode {
        case ExecutionModeTime, ExecutionModeSignal, ExecutionModeCombinedPremium, 
             ExecutionModeManual, ExecutionModeUnderlyingLevel:
                // Valid execution modes
        default:
                return errors.New("invalid execution mode")
        }

        // Validate entry order type
        switch p.EntryOrderType {
        case OrderTypeMarket, OrderTypeLimit, OrderTypeSLLimit:
                // Valid order types
        default:
                return errors.New("invalid entry order type")
        }

        // Validate range breakout settings if enabled
        if p.RangeBreakoutEnabled {
                if !timeRegex.MatchString(p.RangeStartTime) {
                        return errors.New("invalid range start time format (use HH:MM:SS)")
                }
                if !timeRegex.MatchString(p.RangeEndTime) {
                        return errors.New("invalid range end time format (use HH:MM:SS)")
                }
        }

        // Validate monitoring types
        switch p.LegMonitoringType {
        case MonitoringTypeRealtime, MonitoringTypeMinuteClose, MonitoringTypeInterval:
                // Valid monitoring types
        default:
                return errors.New("invalid leg monitoring type")
        }
        switch p.CombinedMonitoringType {
        case MonitoringTypeRealtime, MonitoringTypeMinuteClose, MonitoringTypeInterval:
                // Valid monitoring types
        default:
                return errors.New("invalid combined monitoring type")
        }

        // Validate monitoring interval if interval monitoring is used
        if (p.LegMonitoringType == MonitoringTypeInterval || 
            p.CombinedMonitoringType == MonitoringTypeInterval) && p.MonitoringInterval <= 0 {
                return errors.New("monitoring interval must be greater than zero when interval monitoring is used")
        }

        // Validate target type
        switch p.TargetType {
        case TargetTypeCombinedProfit, TargetTypeCombinedPremium, TargetTypeUnderlying:
                // Valid target types
        default:
                return errors.New("invalid target type")
        }

        // Validate target value
        if p.TargetValue <= 0 {
                return errors.New("target value must be greater than zero")
        }

        // Validate stop loss type
        switch p.StopLossType {
        case StopLossTypeCombinedLoss, StopLossTypeCombinedPremium, 
             StopLossTypeLossAndUnderlyingRange, StopLossTypeDeltaTheta:
                // Valid stop loss types
        default:
                return errors.New("invalid stop loss type")
        }

        // Validate stop loss value
        if p.StopLossValue <= 0 {
                return errors.New("stop loss value must be greater than zero")
        }

        // Validate exit mode
        switch p.ExitMode {
        case ExitModeNormal, ExitModeLegByLeg, ExitModeReverseEntrySequence:
                // Valid exit modes
        default:
                return errors.New("invalid exit mode")
        }

        // Validate exit order type
        switch p.ExitOrderType {
        case OrderTypeMarket, OrderTypeLimit, OrderTypeSLLimit:
                // Valid order types
        default:
                return errors.New("invalid exit order type")
        }

        // Validate exit parameters
        if p.ExitPriceBuffer < 0 {
                return errors.New("exit price buffer cannot be negative")
        }
        if p.MaxExitRetries < 0 {
                return errors.New("max exit retries cannot be negative")
        }
        if p.ExitRetryInterval < 0 {
                return errors.New("exit retry interval cannot be negative")
        }

        // Validate partial exit settings
        if p.EnablePartialExits && (p.MinExitPercentage <= 0 || p.MinExitPercentage > 100) {
                return errors.New("min exit percentage must be between 0 and 100 when partial exits are enabled")
        }

        // Validate legs
        if len(p.Legs) == 0 {
                return errors.New("portfolio must have at least one leg")
        }
        for i, leg := range p.Legs {
                if err := leg.Validate(); err != nil {
                        return errors.New("invalid leg at index " + string(i) + ": " + err.Error())
                }
        }

        return nil
}

// CalculatePnL calculates the profit and loss for the portfolio
func (p *Portfolio) CalculatePnL() {
        p.UnrealizedPnL = 0
        p.RealizedPnL = 0
        p.Delta = 0
        p.Gamma = 0
        p.Theta = 0
        p.Vega = 0

        // Calculate P&L and Greeks from all legs
        for _, leg := range p.Legs {
                p.UnrealizedPnL += leg.UnrealizedPnL
                p.RealizedPnL += leg.RealizedPnL
                p.Delta += leg.Delta
                p.Gamma += leg.Gamma
                p.Theta += leg.Theta
                p.Vega += leg.Vega
        }

        p.TotalPnL = p.UnrealizedPnL + p.RealizedPnL

        // Calculate percentage P&L
        if p.EntryValue > 0 {
                p.PnLPercentage = (p.TotalPnL / p.EntryValue) * 100
        }
}

// UpdateCurrentValues updates the current values of all legs and recalculates P&L
func (p *Portfolio) UpdateCurrentValues() {
        p.CurrentValue = 0
        for i := range p.Legs {
                p.CurrentValue += p.Legs[i].CurrentPrice * float64(p.Legs[i].Quantity)
        }
        p.CalculatePnL()

        // Update max and min values
        if p.CurrentValue > p.MaxValue {
                p.MaxValue = p.CurrentValue
        }
        if p.MinValue == 0 || p.CurrentValue < p.MinValue {
                p.MinValue = p.CurrentValue
        }
}

// IsActive checks if the portfolio is active and should run today
func (p *Portfolio) IsActive() bool {
        if p.Status != PortfolioStatusActive {
                return false
        }

        // Check if portfolio should run today
        today := time.Now().Weekday().String()
        for _, day := range p.RunOnDays {
                if day == today {
                        return true
                }
        }

        return false
}

// ShouldExecuteNow checks if the portfolio should execute at the current time
func (p *Portfolio) ShouldExecuteNow() bool {
        if !p.IsActive() {
                return false
        }

        if p.ExecutionMode != ExecutionModeTime {
                return false // Only time-based portfolios are checked
        }

        now := time.Now()
        currentTime := now.Format("15:04:05")

        // Check if current time is between start and end time
        return currentTime >= p.StartTime && currentTime <= p.EndTime
}

// ShouldSquareOff checks if the portfolio should square off positions at the current time
func (p *Portfolio) ShouldSquareOff() bool {
        if !p.IsActive() {
                return false
        }

        now := time.Now()
        currentTime := now.Format("15:04:05")

        // Check if current time is at or after square off time
        return currentTime >= p.SquareOffTime
}

// AddExecutionLog adds a log entry to the execution logs
func (p *Portfolio) AddExecutionLog(log string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        p.ExecutionLogs = append(p.ExecutionLogs, timestamp+": "+log)
}
