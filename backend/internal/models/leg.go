package models

import (
        "errors"
        "time"
)

// LegType represents the type of leg in a portfolio
type LegType string

const (
        LegTypeOption LegType = "OPTION"
        LegTypeFuture LegType = "FUTURE"
        LegTypeStock  LegType = "STOCK"
)

// Leg represents a single leg in a multi-leg options portfolio
type Leg struct {
        ID                int               `json:"id" bson:"id"`
        PortfolioID       string            `json:"portfolioId" bson:"portfolioId"`
        Symbol            string            `json:"symbol" bson:"symbol"`
        Exchange          string            `json:"exchange" bson:"exchange"`
        Type              LegType           `json:"type" bson:"type"`
        BuySell           string            `json:"buySell" bson:"buySell"`
        OptionType        string            `json:"optionType,omitempty" bson:"optionType,omitempty"`
        StrikePrice       float64           `json:"strikePrice,omitempty" bson:"strikePrice,omitempty"`
        Expiry            time.Time         `json:"expiry,omitempty" bson:"expiry,omitempty"`
        Lots              int               `json:"lots" bson:"lots"`
        LotSize           int               `json:"lotSize" bson:"lotSize"`
        Quantity          int               `json:"quantity" bson:"quantity"`
        EntryPrice        float64           `json:"entryPrice" bson:"entryPrice"`
        CurrentPrice      float64           `json:"currentPrice" bson:"currentPrice"`
        ExitPrice         float64           `json:"exitPrice,omitempty" bson:"exitPrice,omitempty"`
        
        // Strike Selection
        StrikeSelectionMode StrikeSelectionMode `json:"strikeSelectionMode" bson:"strikeSelectionMode"`
        StrikeSelectionValue float64           `json:"strikeSelectionValue,omitempty" bson:"strikeSelectionValue,omitempty"`
        StrikeSelectionType string            `json:"strikeSelectionType,omitempty" bson:"strikeSelectionType,omitempty"`
        
        // Entry Parameters
        EntryOrderType     OrderType         `json:"entryOrderType" bson:"entryOrderType"`
        EntryLimitPrice    float64           `json:"entryLimitPrice,omitempty" bson:"entryLimitPrice,omitempty"`
        EntryTriggerPrice  float64           `json:"entryTriggerPrice,omitempty" bson:"entryTriggerPrice,omitempty"`
        EntryPriceBuffer   float64           `json:"entryPriceBuffer" bson:"entryPriceBuffer"`
        MaxEntryRetries    int               `json:"maxEntryRetries" bson:"maxEntryRetries"`
        EntryRetryInterval int               `json:"entryRetryInterval" bson:"entryRetryInterval"`
        
        // Exit Parameters
        ExitOrderType      OrderType         `json:"exitOrderType" bson:"exitOrderType"`
        ExitLimitPrice     float64           `json:"exitLimitPrice,omitempty" bson:"exitLimitPrice,omitempty"`
        ExitTriggerPrice   float64           `json:"exitTriggerPrice,omitempty" bson:"exitTriggerPrice,omitempty"`
        ExitPriceBuffer    float64           `json:"exitPriceBuffer" bson:"exitPriceBuffer"`
        MaxExitRetries     int               `json:"maxExitRetries" bson:"maxExitRetries"`
        ExitRetryInterval  int               `json:"exitRetryInterval" bson:"exitRetryInterval"`
        
        // Leg-specific Conditions
        EntryCondition     string            `json:"entryCondition,omitempty" bson:"entryCondition,omitempty"`
        ExitCondition      string            `json:"exitCondition,omitempty" bson:"exitCondition,omitempty"`
        DependentLegID     int               `json:"dependentLegId,omitempty" bson:"dependentLegId,omitempty"`
        ExecutionPriority  int               `json:"executionPriority" bson:"executionPriority"`
        
        // Leg-specific Target/Stop Loss
        IndividualTarget   float64           `json:"individualTarget,omitempty" bson:"individualTarget,omitempty"`
        IndividualStopLoss float64           `json:"individualStopLoss,omitempty" bson:"individualStopLoss,omitempty"`
        TrailTarget        bool              `json:"trailTarget" bson:"trailTarget"`
        TrailStopLoss      bool              `json:"trailStopLoss" bson:"trailStopLoss"`
        TrailValue         float64           `json:"trailValue,omitempty" bson:"trailValue,omitempty"`
        
        // Greeks
        Delta              float64           `json:"delta" bson:"delta"`
        Gamma              float64           `json:"gamma" bson:"gamma"`
        Theta              float64           `json:"theta" bson:"theta"`
        Vega               float64           `json:"vega" bson:"vega"`
        
        // Execution Status
        Status             string            `json:"status" bson:"status"`
        EntryOrderID       string            `json:"entryOrderId,omitempty" bson:"entryOrderId,omitempty"`
        ExitOrderID        string            `json:"exitOrderId,omitempty" bson:"exitOrderId,omitempty"`
        EntryTime          time.Time         `json:"entryTime,omitempty" bson:"entryTime,omitempty"`
        ExitTime           time.Time         `json:"exitTime,omitempty" bson:"exitTime,omitempty"`
        EntrySlippage      float64           `json:"entrySlippage" bson:"entrySlippage"`
        ExitSlippage       float64           `json:"exitSlippage" bson:"exitSlippage"`
        
        // Performance Metrics
        UnrealizedPnL      float64           `json:"unrealizedPnL" bson:"unrealizedPnL"`
        RealizedPnL        float64           `json:"realizedPnL" bson:"realizedPnL"`
        TotalPnL           float64           `json:"totalPnL" bson:"totalPnL"`
        PnLPercentage      float64           `json:"pnLPercentage" bson:"pnLPercentage"`
        
        // Additional Information
        Notes              string            `json:"notes,omitempty" bson:"notes,omitempty"`
        Tags               []string          `json:"tags,omitempty" bson:"tags,omitempty"`
        CreatedAt          time.Time         `json:"createdAt" bson:"createdAt"`
        UpdatedAt          time.Time         `json:"updatedAt" bson:"updatedAt"`
}

// Validate validates the leg data
func (l *Leg) Validate() error {
        // Check required fields
        if l.PortfolioID == "" {
                return errors.New("portfolio ID is required")
        }
        if l.Symbol == "" {
                return errors.New("symbol is required")
        }
        if l.Exchange == "" {
                return errors.New("exchange is required")
        }
        if l.Lots <= 0 {
                return errors.New("lots must be greater than zero")
        }
        if l.LotSize <= 0 {
                return errors.New("lot size must be greater than zero")
        }

        // Validate leg type
        switch l.Type {
        case LegTypeOption, LegTypeFuture, LegTypeStock:
                // Valid leg types
        default:
                return errors.New("invalid leg type")
        }

        // Validate buy/sell direction
        if l.BuySell != string(OrderDirectionBuy) && l.BuySell != string(OrderDirectionSell) {
                return errors.New("invalid buy/sell direction")
        }

        // Validate option-specific fields
        if l.Type == LegTypeOption {
                if l.StrikePrice <= 0 {
                        return errors.New("strike price must be greater than zero for options")
                }
                if l.Expiry.IsZero() {
                        return errors.New("expiry date is required for options")
                }
                if l.OptionType != string(OptionTypeCall) && l.OptionType != string(OptionTypePut) {
                        return errors.New("invalid option type")
                }
        }

        // Validate strike selection mode
        switch l.StrikeSelectionMode {
        case StrikeSelectionModeNormal, StrikeSelectionModeRelative, StrikeSelectionModeBoth:
                // Valid strike selection modes
        default:
                return errors.New("invalid strike selection mode")
        }

        // Validate entry order type
        switch l.EntryOrderType {
        case OrderTypeMarket, OrderTypeLimit, OrderTypeSLLimit:
                // Valid order types
        default:
                return errors.New("invalid entry order type")
        }

        // Validate exit order type
        switch l.ExitOrderType {
        case OrderTypeMarket, OrderTypeLimit, OrderTypeSLLimit:
                // Valid order types
        default:
                return errors.New("invalid exit order type")
        }

        // Validate limit prices for limit orders
        if l.EntryOrderType == OrderTypeLimit && l.EntryLimitPrice <= 0 {
                return errors.New("entry limit price must be greater than zero for limit orders")
        }
        if l.ExitOrderType == OrderTypeLimit && l.ExitLimitPrice <= 0 && l.Status != "PENDING" {
                return errors.New("exit limit price must be greater than zero for limit orders")
        }

        // Validate trigger prices for stop-loss limit orders
        if l.EntryOrderType == OrderTypeSLLimit && l.EntryTriggerPrice <= 0 {
                return errors.New("entry trigger price must be greater than zero for stop-loss limit orders")
        }
        if l.ExitOrderType == OrderTypeSLLimit && l.ExitTriggerPrice <= 0 && l.Status != "PENDING" {
                return errors.New("exit trigger price must be greater than zero for stop-loss limit orders")
        }

        // Validate buffer values
        if l.EntryPriceBuffer < 0 {
                return errors.New("entry price buffer cannot be negative")
        }
        if l.ExitPriceBuffer < 0 {
                return errors.New("exit price buffer cannot be negative")
        }

        // Validate retry parameters
        if l.MaxEntryRetries < 0 {
                return errors.New("max entry retries cannot be negative")
        }
        if l.EntryRetryInterval < 0 {
                return errors.New("entry retry interval cannot be negative")
        }
        if l.MaxExitRetries < 0 {
                return errors.New("max exit retries cannot be negative")
        }
        if l.ExitRetryInterval < 0 {
                return errors.New("exit retry interval cannot be negative")
        }

        // Validate execution priority
        if l.ExecutionPriority < 0 {
                return errors.New("execution priority cannot be negative")
        }

        // Validate target and stop loss values if set
        if l.IndividualTarget <= 0 && l.IndividualTarget != 0 {
                return errors.New("individual target must be greater than zero if set")
        }
        if l.IndividualStopLoss <= 0 && l.IndividualStopLoss != 0 {
                return errors.New("individual stop loss must be greater than zero if set")
        }

        // Validate trail value if trailing is enabled
        if (l.TrailTarget || l.TrailStopLoss) && l.TrailValue <= 0 {
                return errors.New("trail value must be greater than zero when trailing is enabled")
        }

        // Validate status
        validStatuses := map[string]bool{
                "PENDING": true, "ACTIVE": true, "COMPLETED": true, "FAILED": true, "CANCELLED": true,
        }
        if !validStatuses[l.Status] {
                return errors.New("invalid leg status")
        }

        return nil
}

// CalculatePnL calculates the profit and loss for the leg
func (l *Leg) CalculatePnL() {
        multiplier := float64(l.Quantity)
        
        if l.BuySell == string(OrderDirectionBuy) {
                l.UnrealizedPnL = (l.CurrentPrice - l.EntryPrice) * multiplier
        } else {
                l.UnrealizedPnL = (l.EntryPrice - l.CurrentPrice) * multiplier
        }
        
        l.TotalPnL = l.UnrealizedPnL + l.RealizedPnL
        
        // Calculate percentage P&L
        if l.EntryPrice > 0 {
                investment := l.EntryPrice * multiplier
                l.PnLPercentage = (l.TotalPnL / investment) * 100
        }
}

// UpdateCurrentPrice updates the current price and recalculates P&L
func (l *Leg) UpdateCurrentPrice(price float64) {
        l.CurrentPrice = price
        l.CalculatePnL()
}

// CalculateQuantity calculates the total quantity based on lots and lot size
func (l *Leg) CalculateQuantity() {
        l.Quantity = l.Lots * l.LotSize
}

// CalculateEntrySlippage calculates the slippage for entry order
func (l *Leg) CalculateEntrySlippage() {
        if l.EntryOrderType == OrderTypeMarket {
                return // No slippage calculation for market orders
        }
        
        if l.BuySell == string(OrderDirectionBuy) {
                l.EntrySlippage = l.EntryPrice - l.EntryLimitPrice
        } else {
                l.EntrySlippage = l.EntryLimitPrice - l.EntryPrice
        }
}

// CalculateExitSlippage calculates the slippage for exit order
func (l *Leg) CalculateExitSlippage() {
        if l.ExitOrderType == OrderTypeMarket || l.ExitLimitPrice == 0 {
                return // No slippage calculation for market orders or if exit limit price is not set
        }
        
        if l.BuySell == string(OrderDirectionBuy) {
                l.ExitSlippage = l.ExitLimitPrice - l.ExitPrice
        } else {
                l.ExitSlippage = l.ExitPrice - l.ExitLimitPrice
        }
}

// GetSymbolWithExpiry returns the full symbol with expiry and strike information
func (l *Leg) GetSymbolWithExpiry() string {
        if l.Type != LegTypeOption {
                return l.Symbol
        }
        
        expiryStr := l.Expiry.Format("02JAN06")
        return l.Symbol + expiryStr + l.OptionType + l.StrikePrice
}

// Clone creates a copy of the leg
func (l *Leg) Clone() *Leg {
        clone := *l
        clone.ID = 0 // Will be reassigned
        clone.EntryOrderID = ""
        clone.ExitOrderID = ""
        clone.EntryTime = time.Time{}
        clone.ExitTime = time.Time{}
        clone.Status = "PENDING"
        clone.EntryPrice = 0
        clone.ExitPrice = 0
        clone.EntrySlippage = 0
        clone.ExitSlippage = 0
        clone.UnrealizedPnL = 0
        clone.RealizedPnL = 0
        clone.TotalPnL = 0
        clone.PnLPercentage = 0
        clone.CreatedAt = time.Now()
        clone.UpdatedAt = time.Now()
        
        return &clone
}
