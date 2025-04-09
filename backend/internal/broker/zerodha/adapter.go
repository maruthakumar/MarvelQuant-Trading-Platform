// Package zerodha provides the Zerodha implementation of the broker interface
package zerodha

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/trading-platform/backend/internal/broker/common"
	kiteconnect "github.com/zerodha/gokiteconnect"
)

// ZerodhaAdapter implements the BrokerClient interface for Zerodha
type ZerodhaAdapter struct {
	apiKey      string
	apiSecret   string
	redirectURL string
	baseURL     string
	debug       bool
	client      *kiteconnect.Client
	httpClient  *http.Client
	accessToken string
	userID      string
}

// NewZerodhaAdapter creates a new Zerodha adapter
func NewZerodhaAdapter(config *common.ZerodhaConfig) (*ZerodhaAdapter, error) {
	if config == nil {
		return nil, errors.New("Zerodha configuration is required")
	}

	if config.APIKey == "" || config.APISecret == "" {
		return nil, errors.New("API key and API secret are required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.kite.trade" // Default URL
	}

	// Create Kite Connect client
	kc := kiteconnect.New(config.APIKey)

	// Set debug if needed
	// kc.SetDebug(true)

	// Set custom base URL if provided
	if baseURL != "https://api.kite.trade" {
		kc.SetBaseURI(baseURL)
	}

	return &ZerodhaAdapter{
		apiKey:      config.APIKey,
		apiSecret:   config.APISecret,
		redirectURL: config.RedirectURL,
		baseURL:     baseURL,
		debug:       false,
		client:      kc,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// SetDebug enables or disables debug mode
func (z *ZerodhaAdapter) SetDebug(debug bool) {
	z.debug = debug
	z.client.SetDebug(debug)
}

// Login authenticates with the Zerodha API
// Note: Zerodha uses a two-step authentication process
// 1. Get a request token from the Kite Connect login URL
// 2. Exchange the request token for an access token
func (z *ZerodhaAdapter) Login(credentials *common.Credentials) (*common.Session, error) {
	if credentials == nil || credentials.TwoFactorCode == "" {
		// If no request token is provided, return the login URL
		loginURL := z.client.GetLoginURL()
		return nil, fmt.Errorf("please visit %s to get the request token", loginURL)
	}

	// Exchange request token for access token
	data, err := z.client.GenerateSession(credentials.TwoFactorCode, z.apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session: %w", err)
	}

	// Set the access token for subsequent requests
	z.client.SetAccessToken(data.AccessToken)
	z.accessToken = data.AccessToken
	z.userID = data.UserID

	// Calculate expiry time (Zerodha tokens typically expire at 6 AM the next day)
	// For simplicity, we'll set it to expire in 24 hours
	expiryTime := time.Now().Add(24 * time.Hour).Unix()

	return &common.Session{
		Token:     data.AccessToken,
		UserID:    data.UserID,
		ExpiresAt: expiryTime,
	}, nil
}

// Logout invalidates the current session
func (z *ZerodhaAdapter) Logout() error {
	if z.accessToken == "" {
		return errors.New("not logged in")
	}

	err := z.client.InvalidateSession()
	if err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}

	// Clear the access token and user ID
	z.accessToken = ""
	z.userID = ""

	return nil
}

// PlaceOrder places a new order with the Zerodha API
func (z *ZerodhaAdapter) PlaceOrder(order *common.Order) (*common.OrderResponse, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	if order == nil {
		return nil, errors.New("order is required")
	}

	// Map common order to Zerodha order
	orderParams := kiteconnect.OrderParams{
		Exchange:        mapExchangeSegment(order.ExchangeSegment),
		TradingSymbol:   order.TradingSymbol,
		TransactionType: mapOrderSide(order.OrderSide),
		Quantity:        order.OrderQuantity,
		Product:         mapProductType(order.ProductType),
		OrderType:       mapOrderType(order.OrderType),
		Validity:        mapTimeInForce(order.TimeInForce),
		Price:           order.LimitPrice,
		TriggerPrice:    order.StopPrice,
		Tag:             order.OrderUniqueIdentifier, // Use tag for order identifier
	}

	// Set variety based on order type
	// For simplicity, we'll use regular variety for all orders
	// In a real implementation, this would be more sophisticated
	variety := kiteconnect.VarietyRegular

	// Place the order
	orderResponse, err := z.client.PlaceOrder(variety, orderParams)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	// Convert the response to the common OrderResponse model
	response := &common.OrderResponse{
		OrderID: orderResponse.OrderID,
		Status:  "PLACED", // Zerodha doesn't return status in place order response
	}

	return response, nil
}

// ModifyOrder modifies an existing order with the Zerodha API
func (z *ZerodhaAdapter) ModifyOrder(order *common.ModifyOrder) (*common.OrderResponse, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	if order == nil {
		return nil, errors.New("order is required")
	}

	if order.OrderID == "" {
		return nil, errors.New("order ID is required")
	}

	// Map common order to Zerodha order
	orderParams := kiteconnect.OrderParams{
		Quantity:     order.OrderQuantity,
		OrderType:    mapOrderType(order.OrderType),
		Price:        order.LimitPrice,
		TriggerPrice: order.StopPrice,
	}

	// Set variety based on order type
	// For simplicity, we'll use regular variety for all orders
	variety := kiteconnect.VarietyRegular

	// Modify the order
	orderResponse, err := z.client.ModifyOrder(variety, order.OrderID, orderParams)
	if err != nil {
		return nil, fmt.Errorf("failed to modify order: %w", err)
	}

	// Convert the response to the common OrderResponse model
	response := &common.OrderResponse{
		OrderID: orderResponse.OrderID,
		Status:  "MODIFIED", // Zerodha doesn't return status in modify order response
	}

	return response, nil
}

// CancelOrder cancels an existing order with the Zerodha API
func (z *ZerodhaAdapter) CancelOrder(orderID string, clientID string) (*common.OrderResponse, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	if orderID == "" {
		return nil, errors.New("order ID is required")
	}

	// Set variety based on order type
	// For simplicity, we'll use regular variety for all orders
	variety := kiteconnect.VarietyRegular

	// Cancel the order
	orderResponse, err := z.client.CancelOrder(variety, orderID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	// Convert the response to the common OrderResponse model
	response := &common.OrderResponse{
		OrderID: orderResponse.OrderID,
		Status:  "CANCELLED", // Zerodha doesn't return status in cancel order response
	}

	return response, nil
}

// GetOrderBook retrieves the order book from the Zerodha API
func (z *ZerodhaAdapter) GetOrderBook(clientID string) (*common.OrderBook, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	// Get the order book
	orders, err := z.client.GetOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	// Convert the response to the common OrderBook model
	orderBook := &common.OrderBook{
		Orders: make([]common.OrderDetails, len(orders)),
	}

	for i, order := range orders {
		orderBook.Orders[i] = common.OrderDetails{
			OrderID:              order.OrderID,
			ExchangeOrderID:      order.ExchangeOrderID,
			ExchangeSegment:      mapExchange(order.Exchange),
			ExchangeInstrumentID: "", // Zerodha doesn't provide instrument ID in order book
			OrderSide:            mapTransactionType(order.TransactionType),
			OrderType:            mapZerodhaOrderType(order.OrderType),
			ProductType:          mapZerodhaProductType(order.Product),
			TimeInForce:          mapZerodhaValidity(order.Validity),
			OrderQuantity:        order.Quantity,
			FilledQuantity:       order.FilledQuantity,
			RemainingQuantity:    order.PendingQuantity,
			LimitPrice:           order.Price,
			StopPrice:            order.TriggerPrice,
			OrderStatus:          mapZerodhaStatus(order.Status),
			OrderTimestamp:       order.OrderTimestamp.Unix() * 1000, // Convert to milliseconds
			LastUpdateTimestamp:  order.ExchangeTimestamp.Unix() * 1000,
			CancelTimestamp:      0, // Zerodha doesn't provide cancel timestamp
			ClientID:             z.userID,
		}
	}

	return orderBook, nil
}

// GetPositions retrieves the positions from the Zerodha API
func (z *ZerodhaAdapter) GetPositions(clientID string) ([]common.Position, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	// Get the positions
	positionsResponse, err := z.client.GetPositions()
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// Combine day and net positions
	allPositions := append(positionsResponse.Day, positionsResponse.Net...)

	// Convert the response to the common Position model
	positions := make([]common.Position, len(allPositions))

	for i, position := range allPositions {
		positions[i] = common.Position{
			ExchangeSegment:      mapExchange(position.Exchange),
			ExchangeInstrumentID: "", // Zerodha doesn't provide instrument ID in positions
			ProductType:          mapZerodhaProductType(position.Product),
			Quantity:             position.Quantity,
			BuyQuantity:          position.BuyQuantity,
			SellQuantity:         position.SellQuantity,
			NetQuantity:          position.Quantity,
			AveragePrice:         position.AveragePrice,
			LastPrice:            position.LastPrice,
			RealizedProfit:       position.Pnl,
			UnrealizedProfit:     0, // Zerodha doesn't provide unrealized profit
			ClientID:             z.userID,
		}
	}

	return positions, nil
}

// GetHoldings retrieves the holdings from the Zerodha API
func (z *ZerodhaAdapter) GetHoldings(clientID string) ([]common.Holding, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	// Get the holdings
	holdingsResponse, err := z.client.GetHoldings()
	if err != nil {
		return nil, fmt.Errorf("failed to get holdings: %w", err)
	}

	// Convert the response to the common Holding model
	holdings := make([]common.Holding, len(holdingsResponse))

	for i, holding := range holdingsResponse {
		holdings[i] = common.Holding{
			ExchangeSegment:      mapExchange(holding.Exchange),
			ExchangeInstrumentID: "", // Zerodha doesn't provide instrument ID in holdings
			TradingSymbol:        holding.TradingSymbol,
			ISIN:                 holding.ISIN,
			Quantity:             holding.Quantity,
			AveragePrice:         holding.AveragePrice,
			LastPrice:            holding.LastPrice,
			RealizedProfit:       0, // Zerodha doesn't provide realized profit
			UnrealizedProfit:     holding.PnL,
			ClientID:             z.userID,
		}
	}

	return holdings, nil
}

// GetQuote retrieves quotes for the specified symbols from the Zerodha API
func (z *ZerodhaAdapter) GetQuote(symbols []string) (map[string]common.Quote, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	if len(symbols) == 0 {
		return nil, errors.New("at least one symbol is required")
	}

	// Get the quotes
	quotesResponse, err := z.client.GetQuote(symbols...)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}

	// Convert the response to the common Quote model
	quotes := make(map[string]common.Quote, len(quotesResponse))

	for symbol, quote := range quotesResponse {
		quotes[symbol] = common.Quote{
			ExchangeSegment:      mapExchange(quote.InstrumentToken),
			ExchangeInstrumentID: fmt.Sprintf("%d", quote.InstrumentToken),
			TradingSymbol:        symbol,
			LastPrice:            quote.LastPrice,
			Open:                 quote.OHLC.Open,
			High:                 quote.OHLC.High,
			Low:                  quote.OHLC.Low,
			Close:                quote.OHLC.Close,
			Volume:               int64(quote.Volume),
			BidPrice:             quote.Depth.Buy[0].Price,
			BidSize:              quote.Depth.Buy[0].Quantity,
			AskPrice:             quote.Depth.Sell[0].Price,
			AskSize:              quote.Depth.Sell[0].Quantity,
			Timestamp:            quote.LastTradeTime.Unix() * 1000, // Convert to milliseconds
		}
	}

	return quotes, nil
}

// SubscribeToQuotes subscribes to real-time quotes for the specified symbols
func (z *ZerodhaAdapter) SubscribeToQuotes(symbols []string) (chan common.Quote, error) {
	if z.accessToken == "" {
		return nil, errors.New("not logged in")
	}

	if len(symbols) == 0 {
		return nil, errors.New("at least one symbol is required")
	}

	// Create a channel for quotes
	quoteChan := make(chan common.Quote, 100)

	// In a real implementation, this would establish a WebSocket connection
	// using the Kite Ticker API and stream quotes to the channel.
	// For now, we'll return an error indicating that this is not yet implemented.

	// TODO: Implement WebSocket connection for real-time quotes using Kite Ticker

	return quoteChan, errors.New("real-time quotes not yet implemented")
}

// UnsubscribeFromQuotes unsubscribes from real-time quotes for the specified symbols
func (z *ZerodhaAdapter) UnsubscribeFromQuotes(symbols []string) error {
	if z.accessToken == "" {
		return errors.New("not logged in")
	}

	if len(symbols) == 0 {
		return errors.New("at least one symbol is required")
	}

	// In a real implementation, this would close the WebSocket connection
	// or unsubscribe from the specified symbols. For now, we'll return an error
	// indicating that this is not yet implemented.

	// TODO: Implement WebSocket unsubscribe for real-time quotes

	return errors.New("real-time quotes not yet implemented")
}

// Helper functions to map between common and Zerodha-specific values

// mapExchangeSegment maps common exchange segment to Zerodha exchange
func mapExchangeSegment(segment string) string {
	switch segment {
	case "NSECM":
		return "NSE"
	case "BSECM":
		return "BSE"
	case "NSEFO":
		return "NFO"
	case "BSEFO":
		return "BFO"
	case "NSECD":
		return "CDS"
	case "MCXFO":
		return "MCX"
	default:
		return segment
	}
}

// mapExchange maps Zerodha exchange to common exchange segment
func mapExchange(exchange string) string {
	switch exchange {
	case "NSE":
		return "NSECM"
	case "BSE":
		return "BSECM"
	case "NFO":
		return "NSEFO"
	case "BFO":
		return "BSEFO"
	case "CDS":
		return "NSECD"
	case "MCX":
		return "MCXFO"
	default:
		return exchange
	}
}

// mapOrderSide maps common order side to Zerodha transaction type
func mapOrderSide(side string) string {
	switch side {
	case "BUY":
		return kiteconnect.TransactionTypeBuy
	case "SELL":
		return kiteconnect.TransactionTypeSell
	default:
		return side
	}
}

// mapTransactionType maps Zerodha transaction type to common order side
func mapTransactionType(transactionType string) string {
	switch transactionType {
	case kiteconnect.TransactionTypeBuy:
		return "BUY"
	case kiteconnect.TransactionTypeSell:
		return "SELL"
	default:
		return transactionType
	}
}

// mapProductType maps common product type to Zerodha product
func mapProductType(productType string) string {
	switch productType {
	case "MIS":
		return kiteconnect.ProductMIS
	case "NRML":
		return kiteconnect.ProductNRML
	case "CNC":
		return kiteconnect.ProductCNC
	default:
		return productType
	}
}

// mapZerodhaProductType maps Zerodha product to common product type
func mapZerodhaProductType(product string) string {
	switch product {
	case kiteconnect.ProductMIS:
		return "MIS"
	case kiteconnect.ProductNRML:
		return "NRML"
	case kiteconnect.ProductCNC:
		return "CNC"
	default:
		return product
	}
}

// mapOrderType maps common order type to Zerodha order type
func mapOrderType(orderType string) string {
	switch orderType {
	case "MARKET":
		return kiteconnect.OrderTypeMarket
	case "LIMIT":
		return kiteconnect.OrderTypeLimit
	case "SL":
		return kiteconnect.OrderTypeSL
	case "SL-M":
		return kiteconnect.OrderTypeSLM
	default:
		return orderType
	}
}

// mapZerodhaOrderType maps Zerodha order type to common order type
func mapZerodhaOrderType(orderType string) string {
	switch orderType {
	case kiteconnect.OrderTypeMarket:
		return "MARKET"
	case kiteconnect.OrderTypeLimit:
		return "LIMIT"
	case kiteconnect.OrderTypeSL:
		return "SL"
	case kiteconnect.OrderTypeSLM:
		return "SL-M"
	default:
		return orderType
	}
}

// mapTimeInForce maps common time in force to Zerodha validity
func mapTimeInForce(timeInForce string) string {
	switch timeInForce {
	case "DAY":
		return kiteconnect.ValidityDay
	case "IOC":
		return kiteconnect.ValidityIOC
	default:
		return timeInForce
	}
}

// mapZerodhaValidity maps Zerodha validity to common time in force
func mapZerodhaValidity(validity string) string {
	switch validity {
	case kiteconnect.ValidityDay:
		return "DAY"
	case kiteconnect.ValidityIOC:
		return "IOC"
	default:
		return validity
	}
}

// mapZerodhaStatus maps Zerodha status to common order status
func mapZerodhaStatus(status string) string {
	switch status {
	case "COMPLETE":
		return "FILLED"
	case "REJECTED":
		return "REJECTED"
	case "CANCELLED":
		return "CANCELLED"
	case "PENDING":
		return "PENDING"
	case "OPEN":
		return "OPEN"
	default:
		return status
	}
}
