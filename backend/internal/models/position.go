package models

import (
	"errors"
	"time"
)

// PositionStatus represents the status of a position
type PositionStatus string

const (
	PositionStatusOpen    PositionStatus = "OPEN"
	PositionStatusClosed  PositionStatus = "CLOSED"
	PositionStatusPartial PositionStatus = "PARTIAL"
)

// PositionDirection represents the direction of a position
type PositionDirection string

const (
	PositionDirectionLong  PositionDirection = "LONG"
	PositionDirectionShort PositionDirection = "SHORT"
)

// Position represents a trading position
type Position struct {
	ID             string            `json:"id" bson:"_id,omitempty"`
	UserID         string            `json:"userId" bson:"userId"`
	OrderID        string            `json:"orderId" bson:"orderId"`
	Symbol         string            `json:"symbol" bson:"symbol"`
	Exchange       string            `json:"exchange" bson:"exchange"`
	Direction      PositionDirection `json:"direction" bson:"direction"`
	EntryPrice     float64           `json:"entryPrice" bson:"entryPrice"`
	ExitPrice      float64           `json:"exitPrice" bson:"exitPrice"`
	Quantity       int               `json:"quantity" bson:"quantity"`
	ExitQuantity   int               `json:"exitQuantity" bson:"exitQuantity"`
	Status         PositionStatus    `json:"status" bson:"status"`
	ProductType    ProductType       `json:"productType" bson:"productType"`
	InstrumentType InstrumentType    `json:"instrumentType" bson:"instrumentType"`
	OptionType     OptionType        `json:"optionType,omitempty" bson:"optionType,omitempty"`
	StrikePrice    float64           `json:"strikePrice,omitempty" bson:"strikePrice,omitempty"`
	Expiry         time.Time         `json:"expiry,omitempty" bson:"expiry,omitempty"`
	UnrealizedPnL  float64           `json:"unrealizedPnL" bson:"unrealizedPnL"`
	RealizedPnL    float64           `json:"realizedPnL" bson:"realizedPnL"`
	Greeks         Greeks            `json:"greeks" bson:"greeks"`
	PortfolioID    string            `json:"portfolioId,omitempty" bson:"portfolioId,omitempty"`
	StrategyID     string            `json:"strategyId,omitempty" bson:"strategyId,omitempty"`
	LegID          string            `json:"legId,omitempty" bson:"legId,omitempty"`
	Tags           []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	CreatedAt      time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt" bson:"updatedAt"`
}

// Greeks represents the option Greeks
type Greeks struct {
	Delta float64 `json:"delta" bson:"delta"`
	Gamma float64 `json:"gamma" bson:"gamma"`
	Theta float64 `json:"theta" bson:"theta"`
	Vega  float64 `json:"vega" bson:"vega"`
}

// AggregatedPosition represents an aggregated view of positions
type AggregatedPosition struct {
	Key           string    `json:"key"`
	GroupBy       string    `json:"groupBy"`
	TotalQuantity int       `json:"totalQuantity"`
	NetQuantity   int       `json:"netQuantity"`
	TotalValue    float64   `json:"totalValue"`
	NetValue      float64   `json:"netValue"`
	PnL           float64   `json:"pnl"`
	Greeks        Greeks    `json:"greeks"`
	PositionCount int       `json:"positionCount"`
}

// PositionFilter represents filter criteria for positions
type PositionFilter struct {
	UserID         string
	Symbol         string
	Status         PositionStatus
	Direction      PositionDirection
	ProductType    ProductType
	InstrumentType InstrumentType
	PortfolioID    string
	StrategyID     string
	OrderID        string
	FromDate       time.Time
	ToDate         time.Time
	Tags           []string
}

// Validate validates the position
func (p *Position) Validate() error {
	if p.UserID == "" {
		return errors.New("user ID is required")
	}
	if p.Symbol == "" {
		return errors.New("symbol is required")
	}
	if p.Exchange == "" {
		return errors.New("exchange is required")
	}
	if p.Direction != PositionDirectionLong && p.Direction != PositionDirectionShort {
		return errors.New("invalid direction")
	}
	if p.EntryPrice <= 0 {
		return errors.New("entry price must be greater than zero")
	}
	if p.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if p.ExitQuantity < 0 {
		return errors.New("exit quantity cannot be negative")
	}
	if p.ExitQuantity > p.Quantity {
		return errors.New("exit quantity cannot exceed position quantity")
	}
	if p.Status != PositionStatusOpen && p.Status != PositionStatusClosed && p.Status != PositionStatusPartial {
		return errors.New("invalid status")
	}
	if p.ProductType == "" {
		return errors.New("product type is required")
	}
	if p.InstrumentType == "" {
		return errors.New("instrument type is required")
	}

	// Validate option-specific fields
	if p.InstrumentType == InstrumentTypeOption {
		if p.OptionType != OptionTypeCall && p.OptionType != OptionTypePut {
			return errors.New("invalid option type")
		}
		if p.StrikePrice <= 0 {
			return errors.New("strike price must be greater than zero")
		}
		if p.Expiry.IsZero() {
			return errors.New("expiry is required for options")
		}
	}

	return nil
}

// CalculateTotalPnL calculates the total P&L (realized + unrealized)
func (p *Position) CalculateTotalPnL() float64 {
	return p.RealizedPnL + p.UnrealizedPnL
}

// CalculatePnLPercentage calculates the P&L as a percentage of the initial investment
func (p *Position) CalculatePnLPercentage() float64 {
	initialInvestment := p.EntryPrice * float64(p.Quantity)
	if initialInvestment == 0 {
		return 0
	}
	return (p.CalculateTotalPnL() / initialInvestment) * 100
}

// CalculateDaysTillExpiry calculates the number of days until expiry
func (p *Position) CalculateDaysTillExpiry() int {
	if p.InstrumentType != InstrumentTypeOption && p.InstrumentType != InstrumentTypeFuture {
		return 0
	}
	if p.Expiry.IsZero() {
		return 0
	}
	
	now := time.Now()
	if now.After(p.Expiry) {
		return 0
	}
	
	return int(p.Expiry.Sub(now).Hours() / 24)
}

// IsExpired checks if the position is expired
func (p *Position) IsExpired() bool {
	if p.InstrumentType != InstrumentTypeOption && p.InstrumentType != InstrumentTypeFuture {
		return false
	}
	if p.Expiry.IsZero() {
		return false
	}
	
	return time.Now().After(p.Expiry)
}

// RemainingQuantity returns the remaining quantity in the position
func (p *Position) RemainingQuantity() int {
	return p.Quantity - p.ExitQuantity
}

// IsFullyClosed checks if the position is fully closed
func (p *Position) IsFullyClosed() bool {
	return p.Status == PositionStatusClosed || p.ExitQuantity >= p.Quantity
}

// IsPartiallyExited checks if the position is partially exited
func (p *Position) IsPartiallyExited() bool {
	return p.ExitQuantity > 0 && p.ExitQuantity < p.Quantity
}

// UpdateStatus updates the position status based on exit quantity
func (p *Position) UpdateStatus() {
	if p.ExitQuantity >= p.Quantity {
		p.Status = PositionStatusClosed
	} else if p.ExitQuantity > 0 {
		p.Status = PositionStatusPartial
	} else {
		p.Status = PositionStatusOpen
	}
}
