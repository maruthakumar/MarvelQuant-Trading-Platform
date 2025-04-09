// Package common provides models shared across all broker implementations
package common

// Credentials represents authentication credentials for a broker
type Credentials struct {
	APIKey     string
	SecretKey  string
	UserID     string
	Password   string
	TwoFactorCode string
}

// Session represents an authenticated session
type Session struct {
	Token        string
	UserID       string
	ExpiresAt    int64
	RefreshToken string
}

// Order represents a trading order
type Order struct {
	ExchangeSegment      string
	ExchangeInstrumentID string
	ProductType          string
	OrderType            string
	OrderSide            string
	TimeInForce          string
	DisclosedQuantity    int
	OrderQuantity        int
	LimitPrice           float64
	StopPrice            float64
	OrderUniqueIdentifier string
	ClientID             string
	
	// XTS Client specific fields
	APIOrderSource       string
	
	// Zerodha specific fields
	Variety              string
	TradingSymbol        string
}

// ModifyOrder represents an order modification request
type ModifyOrder struct {
	OrderID              string
	ExchangeSegment      string
	ExchangeInstrumentID string
	OrderType            string
	OrderQuantity        int
	LimitPrice           float64
	StopPrice            float64
	ClientID             string
}

// OrderResponse represents a response to an order operation
type OrderResponse struct {
	OrderID         string
	ExchangeOrderID string
	Status          string
	StatusMessage   string
	RejectionReason string
}

// OrderBook represents a collection of orders
type OrderBook struct {
	Orders []OrderDetails
}

// OrderDetails represents detailed information about an order
type OrderDetails struct {
	OrderID              string
	ExchangeOrderID      string
	ExchangeSegment      string
	ExchangeInstrumentID string
	OrderSide            string
	OrderType            string
	ProductType          string
	TimeInForce          string
	OrderQuantity        int
	FilledQuantity       int
	RemainingQuantity    int
	LimitPrice           float64
	StopPrice            float64
	OrderStatus          string
	OrderTimestamp       int64
	LastUpdateTimestamp  int64
	CancelTimestamp      int64
	ClientID             string
}

// Position represents a trading position
type Position struct {
	ExchangeSegment      string
	ExchangeInstrumentID string
	ProductType          string
	Quantity             int
	BuyQuantity          int
	SellQuantity         int
	NetQuantity          int
	AveragePrice         float64
	LastPrice            float64
	RealizedProfit       float64
	UnrealizedProfit     float64
	ClientID             string
}

// Holding represents a security holding
type Holding struct {
	ExchangeSegment      string
	ExchangeInstrumentID string
	TradingSymbol        string
	ISIN                 string
	Quantity             int
	AveragePrice         float64
	LastPrice            float64
	RealizedProfit       float64
	UnrealizedProfit     float64
	ClientID             string
}

// Quote represents market data for a security
type Quote struct {
	ExchangeSegment      string
	ExchangeInstrumentID string
	TradingSymbol        string
	LastPrice            float64
	Open                 float64
	High                 float64
	Low                  float64
	Close                float64
	Volume               int64
	BidPrice             float64
	BidSize              int
	AskPrice             float64
	AskSize              int
	Timestamp            int64
}
