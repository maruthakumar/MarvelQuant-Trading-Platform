package models

import "time"

// Session represents an authenticated session with XTS
type Session struct {
	Token           string
	UserID          string
	IsInvestorClient bool
	ExpiresAt       time.Time
}

// Order represents an order to be placed on XTS
type Order struct {
	ExchangeSegment       string  `json:"exchangeSegment"`
	ExchangeInstrumentID  string  `json:"exchangeInstrumentID"`
	ProductType           string  `json:"productType"`
	OrderType             string  `json:"orderType"`
	OrderSide             string  `json:"orderSide"`
	TimeInForce           string  `json:"timeInForce"`
	DisclosedQuantity     int     `json:"disclosedQuantity"`
	OrderQuantity         int     `json:"orderQuantity"`
	LimitPrice            float64 `json:"limitPrice"`
	StopPrice             float64 `json:"stopPrice"`
	OrderUniqueIdentifier string  `json:"orderUniqueIdentifier"`
	APIOrderSource        string  `json:"apiOrderSource,omitempty"`
	ClientID              string  `json:"clientID,omitempty"`
}

// OrderResponse represents the response from placing an order
type OrderResponse struct {
	OrderID         string `json:"orderID"`
	ExchangeOrderID string `json:"exchangeOrderID"`
	Status          string `json:"status"`
	ErrorCode       string `json:"errorCode,omitempty"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
}

// BracketOrder represents a bracket order to be placed on XTS
type BracketOrder struct {
	ExchangeSegment       string  `json:"exchangeSegment"`
	ExchangeInstrumentID  string  `json:"exchangeInstrumentID"`
	OrderType             string  `json:"orderType"`
	OrderSide             string  `json:"orderSide"`
	DisclosedQuantity     int     `json:"disclosedQuantity"`
	OrderQuantity         int     `json:"orderQuantity"`
	LimitPrice            float64 `json:"limitPrice"`
	SquareOff             float64 `json:"squarOff"`
	StopLossPrice         float64 `json:"stopLossPrice"`
	TrailingStoploss      float64 `json:"trailingStoploss"`
	IsProOrder            bool    `json:"isProOrder"`
	APIOrderSource        string  `json:"apiOrderSource"`
	OrderUniqueIdentifier string  `json:"orderUniqueIdentifier"`
}

// CoverOrder represents a cover order to be placed on XTS
type CoverOrder struct {
	ExchangeSegment       string  `json:"exchangeSegment"`
	ExchangeInstrumentID  string  `json:"exchangeInstrumentID"`
	OrderSide             string  `json:"orderSide"`
	OrderType             string  `json:"orderType"`
	OrderQuantity         int     `json:"orderQuantity"`
	DisclosedQuantity     int     `json:"disclosedQuantity"`
	LimitPrice            float64 `json:"limitPrice"`
	StopPrice             float64 `json:"stopPrice"`
	OrderUniqueIdentifier string  `json:"orderUniqueIdentifier"`
	APIOrderSource        string  `json:"apiOrderSource"`
}

// ModifyOrder represents an order modification request
type ModifyOrder struct {
	AppOrderID              int     `json:"appOrderID"`
	ModifiedProductType     string  `json:"modifiedProductType"`
	ModifiedOrderType       string  `json:"modifiedOrderType"`
	ModifiedOrderQuantity   int     `json:"modifiedOrderQuantity"`
	ModifiedDisclosedQuantity int   `json:"modifiedDisclosedQuantity"`
	ModifiedLimitPrice      float64 `json:"modifiedLimitPrice"`
	ModifiedStopPrice       float64 `json:"modifiedStopPrice"`
	ModifiedTimeInForce     string  `json:"modifiedTimeInForce"`
	OrderUniqueIdentifier   string  `json:"orderUniqueIdentifier"`
	ClientID                string  `json:"clientID,omitempty"`
}

// Quote represents a market quote
type Quote struct {
	ExchangeSegment      string       `json:"exchangeSegment"`
	ExchangeInstrumentID string       `json:"exchangeInstrumentID"`
	Timestamp            time.Time    `json:"timestamp"`
	LastTradedPrice      float64      `json:"lastTradedPrice"`
	LastTradedQuantity   int          `json:"lastTradedQuantity"`
	TotalBuyQuantity     int          `json:"totalBuyQuantity"`
	TotalSellQuantity    int          `json:"totalSellQuantity"`
	BestBids             []PriceLevel `json:"bestBids"`
	BestAsks             []PriceLevel `json:"bestAsks"`
	Open                 float64      `json:"open"`
	High                 float64      `json:"high"`
	Low                  float64      `json:"low"`
	Close                float64      `json:"close"`
	TotalTradedVolume    int          `json:"totalTradedVolume"`
}

// PriceLevel represents a price level in the order book
type PriceLevel struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// Position represents a trading position
type Position struct {
	ExchangeSegment      string  `json:"exchangeSegment"`
	ExchangeInstrumentID string  `json:"exchangeInstrumentID"`
	ProductType          string  `json:"productType"`
	Quantity             int     `json:"quantity"`
	AveragePrice         float64 `json:"averagePrice"`
	CurrentPrice         float64 `json:"currentPrice"`
	PnL                  float64 `json:"pnl"`
	PnLPercentage        float64 `json:"pnlPercentage"`
	PositionID           string  `json:"positionID,omitempty"`
}

// Instrument represents a trading instrument
type Instrument struct {
	ExchangeSegment      string  `json:"exchangeSegment"`
	ExchangeInstrumentID string  `json:"exchangeInstrumentID"`
	InstrumentType       string  `json:"instrumentType"`
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	Series               string  `json:"series"`
	NameWithSeries       string  `json:"nameWithSeries"`
	InstrumentID         int     `json:"instrumentID"`
	TickSize             float64 `json:"tickSize"`
	LotSize              int     `json:"lotSize"`
	ExpiryDate           string  `json:"expiryDate,omitempty"`
	StrikePrice          float64 `json:"strikePrice,omitempty"`
	OptionType           string  `json:"optionType,omitempty"`
}

// LoginRequest represents a login request to XTS
type LoginRequest struct {
	AppKey     string `json:"appKey"`
	SecretKey  string `json:"secretKey"`
	Source     string `json:"source"`
}

// LoginResponse represents the response from a login request
type LoginResponse struct {
	Type        string `json:"type"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Result      struct {
		Token            string `json:"token"`
		UserID           string `json:"userID"`
		IsInvestorClient bool   `json:"isInvestorClient"`
		ExchangeSegments []string `json:"exchangeSegments"`
	} `json:"result"`
}

// OrderBook represents the order book
type OrderBook struct {
	Orders []OrderDetails `json:"orders"`
}

// OrderDetails represents the details of an order
type OrderDetails struct {
	AppOrderID           int     `json:"appOrderID"`
	OrderSide            string  `json:"orderSide"`
	OrderType            string  `json:"orderType"`
	ProductType          string  `json:"productType"`
	TimeInForce          string  `json:"timeInForce"`
	OrderQuantity        int     `json:"orderQuantity"`
	OrderPrice           float64 `json:"orderPrice"`
	OrderStopPrice       float64 `json:"orderStopPrice"`
	OrderStatus          string  `json:"orderStatus"`
	OrderAverageTradedPrice float64 `json:"orderAverageTradedPrice"`
	OrderDisclosedQuantity int   `json:"orderDisclosedQuantity"`
	OrderTradedQuantity  int     `json:"orderTradedQuantity"`
	ExchangeOrderID      string  `json:"exchangeOrderID"`
	ExchangeSegment      string  `json:"exchangeSegment"`
	ExchangeInstrumentID string  `json:"exchangeInstrumentID"`
	OrderGeneratedDateTime string `json:"orderGeneratedDateTime"`
	LastUpdateDateTime   string  `json:"lastUpdateDateTime"`
	CancelRejectReason   string  `json:"cancelRejectReason"`
	OrderUniqueIdentifier string `json:"orderUniqueIdentifier"`
}

// Constants for XTS API
const (
	// Product types
	ProductMIS  = "MIS"
	ProductNRML = "NRML"

	// Order types
	OrderTypeMarket     = "MARKET"
	OrderTypeLimit      = "LIMIT"
	OrderTypeStopMarket = "STOPMARKET"
	OrderTypeStopLimit  = "STOPLIMIT"

	// Transaction types
	TransactionTypeBuy  = "BUY"
	TransactionTypeSell = "SELL"

	// Squareoff modes
	SquareoffDaywise = "DayWise"
	SquareoffNetwise = "Netwise"

	// Squareoff position quantity types
	SquareoffQuantityExact      = "ExactQty"
	SquareoffQuantityPercentage = "Percentage"

	// Validity
	ValidityDay = "DAY"

	// Exchange segments
	ExchangeNSECM = "NSECM"
	ExchangeNSEFO = "NSEFO"
	ExchangeNSECD = "NSECD"
	ExchangeMCXFO = "MCXFO"
	ExchangeBSECM = "BSECM"
	ExchangeBSEFO = "BSEFO"
)
