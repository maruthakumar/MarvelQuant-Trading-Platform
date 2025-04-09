package models

import (
        "errors"
        "time"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
        OrderStatusPending   OrderStatus = "PENDING"
        OrderStatusExecuted  OrderStatus = "EXECUTED"
        OrderStatusCancelled OrderStatus = "CANCELLED"
        OrderStatusRejected  OrderStatus = "REJECTED"
        OrderStatusPartial   OrderStatus = "PARTIAL"
)

// OrderType represents the type of order
type OrderType string

const (
        OrderTypeMarket  OrderType = "MARKET"
        OrderTypeLimit   OrderType = "LIMIT"
        OrderTypeSLLimit OrderType = "SL_LIMIT"
)

// OrderDirection represents the direction of the order (buy/sell)
type OrderDirection string

const (
        OrderDirectionBuy  OrderDirection = "BUY"
        OrderDirectionSell OrderDirection = "SELL"
)

// ProductType represents the product type for the order
type ProductType string

const (
        ProductTypeMIS  ProductType = "MIS"
        ProductTypeNRML ProductType = "NRML"
        ProductTypeCNC  ProductType = "CNC"
)

// InstrumentType represents the type of instrument
type InstrumentType string

const (
        InstrumentTypeOption InstrumentType = "OPTION"
        InstrumentTypeFuture InstrumentType = "FUTURE"
        InstrumentTypeStock  InstrumentType = "STOCK"
)

// OptionType represents the type of option
type OptionType string

const (
        OptionTypeCall OptionType = "CE"
        OptionTypePut  OptionType = "PE"
)

// Order represents an order in the trading system
type Order struct {
        ID              string          `json:"id" bson:"_id,omitempty"`
        UserID          string          `json:"userId" bson:"userId"`
        Symbol          string          `json:"symbol" bson:"symbol"`
        Exchange        string          `json:"exchange" bson:"exchange"`
        OrderType       OrderType       `json:"orderType" bson:"orderType"`
        Direction       OrderDirection  `json:"direction" bson:"direction"`
        Quantity        int             `json:"quantity" bson:"quantity"`
        FilledQuantity  int             `json:"filledQuantity" bson:"filledQuantity"`
        Price           float64         `json:"price" bson:"price"`
        TriggerPrice    float64         `json:"triggerPrice,omitempty" bson:"triggerPrice,omitempty"`
        Status          OrderStatus     `json:"status" bson:"status"`
        ProductType     ProductType     `json:"productType" bson:"productType"`
        InstrumentType  InstrumentType  `json:"instrumentType" bson:"instrumentType"`
        OptionType      OptionType      `json:"optionType,omitempty" bson:"optionType,omitempty"`
        StrikePrice     float64         `json:"strikePrice,omitempty" bson:"strikePrice,omitempty"`
        Expiry          time.Time       `json:"expiry,omitempty" bson:"expiry,omitempty"`
        PortfolioID     string          `json:"portfolioId,omitempty" bson:"portfolioId,omitempty"`
        StrategyID      string          `json:"strategyId,omitempty" bson:"strategyId,omitempty"`
        LegID           int             `json:"legId,omitempty" bson:"legId,omitempty"`
        ParentOrderID   string          `json:"parentOrderId,omitempty" bson:"parentOrderId,omitempty"`
        BrokerOrderID   string          `json:"brokerOrderId,omitempty" bson:"brokerOrderId,omitempty"`
        AveragePrice    float64         `json:"averagePrice" bson:"averagePrice"`
        Slippage        float64         `json:"slippage" bson:"slippage"`
        ExecutionTime   time.Time       `json:"executionTime,omitempty" bson:"executionTime,omitempty"`
        CreatedAt       time.Time       `json:"createdAt" bson:"createdAt"`
        UpdatedAt       time.Time       `json:"updatedAt" bson:"updatedAt"`
        Tags            []string        `json:"tags,omitempty" bson:"tags,omitempty"`
        Notes           string          `json:"notes,omitempty" bson:"notes,omitempty"`
        ErrorMessage    string          `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`
}

// OrderFilter represents filters for querying orders
type OrderFilter struct {
        UserID         string          `json:"userId,omitempty"`
        Symbol         string          `json:"symbol,omitempty"`
        Status         OrderStatus     `json:"status,omitempty"`
        Direction      OrderDirection  `json:"direction,omitempty"`
        ProductType    ProductType     `json:"productType,omitempty"`
        InstrumentType InstrumentType  `json:"instrumentType,omitempty"`
        PortfolioID    string          `json:"portfolioId,omitempty"`
        StrategyID     string          `json:"strategyId,omitempty"`
        FromDate       time.Time       `json:"fromDate,omitempty"`
        ToDate         time.Time       `json:"toDate,omitempty"`
        Tags           []string        `json:"tags,omitempty"`
}

// Validate validates the order data
func (o *Order) Validate() error {
        // Check required fields
        if o.UserID == "" {
                return errors.New("user ID is required")
        }
        if o.Symbol == "" {
                return errors.New("symbol is required")
        }
        if o.Exchange == "" {
                return errors.New("exchange is required")
        }
        if o.Quantity <= 0 {
                return errors.New("quantity must be greater than zero")
        }

        // Validate order type
        switch o.OrderType {
        case OrderTypeMarket, OrderTypeLimit, OrderTypeSLLimit:
                // Valid order types
        default:
                return errors.New("invalid order type")
        }

        // Validate direction
        switch o.Direction {
        case OrderDirectionBuy, OrderDirectionSell:
                // Valid directions
        default:
                return errors.New("invalid order direction")
        }

        // Validate product type
        switch o.ProductType {
        case ProductTypeMIS, ProductTypeNRML, ProductTypeCNC:
                // Valid product types
        default:
                return errors.New("invalid product type")
        }

        // Validate instrument type
        switch o.InstrumentType {
        case InstrumentTypeOption, InstrumentTypeFuture, InstrumentTypeStock:
                // Valid instrument types
        default:
                return errors.New("invalid instrument type")
        }

        // Validate option-specific fields
        if o.InstrumentType == InstrumentTypeOption {
                if o.StrikePrice <= 0 {
                        return errors.New("strike price must be greater than zero for options")
                }
                if o.Expiry.IsZero() {
                        return errors.New("expiry date is required for options")
                }
                switch o.OptionType {
                case OptionTypeCall, OptionTypePut:
                        // Valid option types
                default:
                        return errors.New("invalid option type")
                }
        }

        // Validate price for limit orders
        if o.OrderType == OrderTypeLimit || o.OrderType == OrderTypeSLLimit {
                if o.Price <= 0 {
                        return errors.New("price must be greater than zero for limit orders")
                }
        }

        // Validate trigger price for stop-loss limit orders
        if o.OrderType == OrderTypeSLLimit {
                if o.TriggerPrice <= 0 {
                        return errors.New("trigger price must be greater than zero for stop-loss limit orders")
                }
        }

        // Validate filled quantity
        if o.FilledQuantity < 0 || o.FilledQuantity > o.Quantity {
                return errors.New("filled quantity must be between 0 and total quantity")
        }

        // Validate status
        switch o.Status {
        case OrderStatusPending, OrderStatusExecuted, OrderStatusCancelled, OrderStatusRejected, OrderStatusPartial:
                // Valid statuses
        default:
                return errors.New("invalid order status")
        }

        // Validate status and filled quantity consistency
        if o.Status == OrderStatusExecuted && o.FilledQuantity != o.Quantity {
                return errors.New("executed orders must have filled quantity equal to total quantity")
        }
        if o.Status == OrderStatusPartial && (o.FilledQuantity <= 0 || o.FilledQuantity >= o.Quantity) {
                return errors.New("partial orders must have filled quantity between 0 and total quantity")
        }

        return nil
}

// CalculateSlippage calculates the slippage for the order
func (o *Order) CalculateSlippage() float64 {
        if o.Status != OrderStatusExecuted && o.Status != OrderStatusPartial {
                return 0
        }

        if o.OrderType == OrderTypeMarket {
                return 0 // No slippage calculation for market orders
        }

        if o.Direction == OrderDirectionBuy {
                return o.AveragePrice - o.Price
        }

        return o.Price - o.AveragePrice
}

// IsComplete checks if the order is completely filled
func (o *Order) IsComplete() bool {
        return o.FilledQuantity == o.Quantity
}

// RemainingQuantity returns the remaining quantity to be filled
func (o *Order) RemainingQuantity() int {
        return o.Quantity - o.FilledQuantity
}
