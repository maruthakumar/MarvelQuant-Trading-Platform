package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// XTSBroker implements the Broker interface for XTS
type XTSBroker struct {
	config      BrokerConfig
	httpClient  *http.Client
	accessToken string
	userID      string
	isConnected bool
}

// NewXTSBroker creates a new XTS broker
func NewXTSBroker() *XTSBroker {
	return &XTSBroker{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Initialize initializes the broker with configuration
func (b *XTSBroker) Initialize(config BrokerConfig) error {
	b.config = config
	
	// Login to XTS API
	loginURL := fmt.Sprintf("%s/user/login", b.config.Endpoint)
	
	payload := map[string]string{
		"appKey":    b.config.APIKey,
		"secretKey": b.config.APISecret,
		"source":    "WEBAPI",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}
	
	var loginResponse struct {
		Result struct {
			UserID      string `json:"userID"`
			AccessToken string `json:"token"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal login response: %w", err)
	}
	
	if loginResponse.Status != "Success" {
		return fmt.Errorf("login failed: %s", loginResponse.Message)
	}
	
	b.accessToken = loginResponse.Result.AccessToken
	b.userID = loginResponse.Result.UserID
	b.isConnected = true
	
	return nil
}

// PlaceOrder places an order
func (b *XTSBroker) PlaceOrder(ctx context.Context, request OrderRequest) (*OrderResponse, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	orderURL := fmt.Sprintf("%s/orders", b.config.Endpoint)
	
	// Map our generic order request to XTS-specific format
	xtsOrderType := "MARKET"
	switch request.OrderType {
	case OrderTypeMarket:
		xtsOrderType = "MARKET"
	case OrderTypeLimit:
		xtsOrderType = "LIMIT"
	case OrderTypeStopLoss:
		xtsOrderType = "STOPLOSS"
	case OrderTypeStopLimit:
		xtsOrderType = "STOPLIMIT"
	}
	
	xtsTransactionType := "BUY"
	if request.TransactionType == TransactionTypeSell {
		xtsTransactionType = "SELL"
	}
	
	xtsProductType := "NRML"
	if request.ProductType == ProductTypeMIS {
		xtsProductType = "MIS"
	}
	
	payload := map[string]interface{}{
		"exchangeSegment":   request.Exchange,
		"exchangeInstrumentID": request.Symbol,
		"orderType":         xtsOrderType,
		"orderSide":         xtsTransactionType,
		"productType":       xtsProductType,
		"quantity":          request.Quantity,
	}
	
	if request.OrderType == OrderTypeLimit || request.OrderType == OrderTypeStopLimit {
		payload["price"] = request.Price
	}
	
	if request.OrderType == OrderTypeStopLoss || request.OrderType == OrderTypeStopLimit {
		payload["triggerPrice"] = request.TriggerPrice
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", orderURL, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create order request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send order request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read order response: %w", err)
	}
	
	var orderResponse struct {
		Result struct {
			OrderID string `json:"AppOrderID"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}
	
	if orderResponse.Status != "Success" {
		return &OrderResponse{
			Success:      false,
			ErrorMessage: orderResponse.Message,
		}, nil
	}
	
	return &OrderResponse{
		Success: true,
		OrderID: orderResponse.Result.OrderID,
	}, nil
}

// ModifyOrder modifies an existing order
func (b *XTSBroker) ModifyOrder(ctx context.Context, orderID string, request OrderRequest) (*OrderResponse, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	modifyURL := fmt.Sprintf("%s/orders/%s", b.config.Endpoint, orderID)
	
	// Map our generic order request to XTS-specific format
	payload := map[string]interface{}{
		"appOrderID": orderID,
		"quantity":   request.Quantity,
	}
	
	if request.OrderType == OrderTypeLimit || request.OrderType == OrderTypeStopLimit {
		payload["price"] = request.Price
	}
	
	if request.OrderType == OrderTypeStopLoss || request.OrderType == OrderTypeStopLimit {
		payload["triggerPrice"] = request.TriggerPrice
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modify order payload: %w", err)
	}
	
	req, err := http.NewRequest("PUT", modifyURL, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create modify order request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send modify order request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read modify order response: %w", err)
	}
	
	var modifyResponse struct {
		Result struct {
			OrderID string `json:"AppOrderID"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &modifyResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal modify order response: %w", err)
	}
	
	if modifyResponse.Status != "Success" {
		return &OrderResponse{
			Success:      false,
			ErrorMessage: modifyResponse.Message,
		}, nil
	}
	
	return &OrderResponse{
		Success: true,
		OrderID: modifyResponse.Result.OrderID,
	}, nil
}

// CancelOrder cancels an order
func (b *XTSBroker) CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	cancelURL := fmt.Sprintf("%s/orders/%s", b.config.Endpoint, orderID)
	
	req, err := http.NewRequest("DELETE", cancelURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cancel order request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send cancel order request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read cancel order response: %w", err)
	}
	
	var cancelResponse struct {
		Result struct {
			OrderID string `json:"AppOrderID"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &cancelResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cancel order response: %w", err)
	}
	
	if cancelResponse.Status != "Success" {
		return &OrderResponse{
			Success:      false,
			ErrorMessage: cancelResponse.Message,
		}, nil
	}
	
	return &OrderResponse{
		Success: true,
		OrderID: cancelResponse.Result.OrderID,
	}, nil
}

// GetOrder gets an order by ID
func (b *XTSBroker) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	orderURL := fmt.Sprintf("%s/orders/%s", b.config.Endpoint, orderID)
	
	req, err := http.NewRequest("GET", orderURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get order request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get order request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read get order response: %w", err)
	}
	
	var orderResponse struct {
		Result struct {
			OrderID         string  `json:"AppOrderID"`
			ExchangeOrderID string  `json:"OrderID"`
			Symbol          string  `json:"ExchangeInstrumentID"`
			Exchange        string  `json:"ExchangeSegment"`
			OrderSide       string  `json:"OrderSide"`
			OrderType       string  `json:"OrderType"`
			ProductType     string  `json:"ProductType"`
			Quantity        int     `json:"Quantity"`
			Price           float64 `json:"Price"`
			TriggerPrice    float64 `json:"TriggerPrice"`
			Status          string  `json:"OrderStatus"`
			StatusMessage   string  `json:"StatusMessage"`
			OrderTime       string  `json:"OrderGeneratedDateTime"`
			LastUpdateTime  string  `json:"LastUpdateDateTime"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal get order response: %w", err)
	}
	
	if orderResponse.Status != "Success" {
		return nil, fmt.Errorf("get order failed: %s", orderResponse.Message)
	}
	
	// Map XTS-specific order status to our generic status
	var orderStatus OrderStatus
	switch orderResponse.Result.Status {
	case "Pending":
		orderStatus = OrderStatusPending
	case "Open":
		orderStatus = OrderStatusOpen
	case "Completed":
		orderStatus = OrderStatusCompleted
	case "Cancelled":
		orderStatus = OrderStatusCancelled
	case "Rejected":
		orderStatus = OrderStatusRejected
	default:
		orderStatus = OrderStatusPending
	}
	
	// Map XTS-specific order type to our generic type
	var orderType OrderType
	switch orderResponse.Result.OrderType {
	case "MARKET":
		orderType = OrderTypeMarket
	case "LIMIT":
		orderType = OrderTypeLimit
	case "STOPLOSS":
		orderType = OrderTypeStopLoss
	case "STOPLIMIT":
		orderType = OrderTypeStopLimit
	default:
		orderType = OrderTypeMarket
	}
	
	// Map XTS-specific transaction type to our generic type
	var transactionType TransactionType
	if orderResponse.Result.OrderSide == "BUY" {
		transactionType = TransactionTypeBuy
	} else {
		transactionType = TransactionTypeSell
	}
	
	// Map XTS-specific product type to our generic type
	var productType ProductType
	if orderResponse.Result.ProductType == "MIS" {
		productType = ProductTypeMIS
	} else {
		productType = ProductTypeNRML
	}
	
	// Parse timestamps
	orderTime, _ := time.Parse(time.RFC3339, orderResponse.Result.OrderTime)
	updateTime, _ := time.Parse(time.RFC3339, orderResponse.Result.LastUpdateTime)
	
	return &Order{
		ID:              orderResponse.Result.OrderID,
		BrokerOrderID:   orderResponse.Result.ExchangeOrderID,
		Symbol:          orderResponse.Result.Symbol,
		Exchange:        orderResponse.Result.Exchange,
		OrderType:       orderType,
		TransactionType: transactionType,
		ProductType:     productType,
		Quantity:        orderResponse.Result.Quantity,
		Price:           orderResponse.Result.Price,
		TriggerPrice:    orderResponse.Result.TriggerPrice,
		Status:          orderStatus,
		Message:         orderResponse.Result.StatusMessage,
		OrderTimestamp:  orderTime,
		UpdatedAt:       updateTime,
	}, nil
}

// GetOrders gets all orders
func (b *XTSBroker) GetOrders(ctx context.Context) ([]Order, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	ordersURL := fmt.Sprintf("%s/orders", b.config.Endpoint)
	
	req, err := http.NewRequest("GET", ordersURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get orders request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get orders request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read get orders response: %w", err)
	}
	
	var ordersResponse struct {
		Result []struct {
			OrderID         string  `json:"AppOrderID"`
			ExchangeOrderID string  `json:"OrderID"`
			Symbol          string  `json:"ExchangeInstrumentID"`
			Exchange        string  `json:"ExchangeSegment"`
			OrderSide       string  `json:"OrderSide"`
			OrderType       string  `json:"OrderType"`
			ProductType     string  `json:"ProductType"`
			Quantity        int     `json:"Quantity"`
			Price           float64 `json:"Price"`
			TriggerPrice    float64 `json:"TriggerPrice"`
			Status          string  `json:"OrderStatus"`
			StatusMessage   string  `json:"StatusMessage"`
			OrderTime       string  `json:"OrderGeneratedDateTime"`
			LastUpdateTime  string  `json:"LastUpdateDateTime"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &ordersResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal get orders response: %w", err)
	}
	
	if ordersResponse.Status != "Success" {
		return nil, fmt.Errorf("get orders failed: %s", ordersResponse.Message)
	}
	
	orders := make([]Order, 0, len(ordersResponse.Result))
	
	for _, xtsOrder := range ordersResponse.Result {
		// Map XTS-specific order status to our generic status
		var orderStatus OrderStatus
		switch xtsOrder.Status {
		case "Pending":
			orderStatus = OrderStatusPending
		case "Open":
			orderStatus = OrderStatusOpen
		case "Completed":
			orderStatus = OrderStatusCompleted
		case "Cancelled":
			orderStatus = OrderStatusCancelled
		case "Rejected":
			orderStatus = OrderStatusRejected
		default:
			orderStatus = OrderStatusPending
		}
		
		// Map XTS-specific order type to our generic type
		var orderType OrderType
		switch xtsOrder.OrderType {
		case "MARKET":
			orderType = OrderTypeMarket
		case "LIMIT":
			orderType = OrderTypeLimit
		case "STOPLOSS":
			orderType = OrderTypeStopLoss
		case "STOPLIMIT":
			orderType = OrderTypeStopLimit
		default:
			orderType = OrderTypeMarket
		}
		
		// Map XTS-specific transaction type to our generic type
		var transactionType TransactionType
		if xtsOrder.OrderSide == "BUY" {
			transactionType = TransactionTypeBuy
		} else {
			transactionType = TransactionTypeSell
		}
		
		// Map XTS-specific product type to our generic type
		var productType ProductType
		if xtsOrder.ProductType == "MIS" {
			productType = ProductTypeMIS
		} else {
			productType = ProductTypeNRML
		}
		
		// Parse timestamps
		orderTime, _ := time.Parse(time.RFC3339, xtsOrder.OrderTime)
		updateTime, _ := time.Parse(time.RFC3339, xtsOrder.LastUpdateTime)
		
		orders = append(orders, Order{
			ID:              xtsOrder.OrderID,
			BrokerOrderID:   xtsOrder.ExchangeOrderID,
			Symbol:          xtsOrder.Symbol,
			Exchange:        xtsOrder.Exchange,
			OrderType:       orderType,
			TransactionType: transactionType,
			ProductType:     productType,
			Quantity:        xtsOrder.Quantity,
			Price:           xtsOrder.Price,
			TriggerPrice:    xtsOrder.TriggerPrice,
			Status:          orderStatus,
			Message:         xtsOrder.StatusMessage,
			OrderTimestamp:  orderTime,
			UpdatedAt:       updateTime,
		})
	}
	
	return orders, nil
}

// GetPositions gets all positions
func (b *XTSBroker) GetPositions(ctx context.Context) ([]Position, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	positionsURL := fmt.Sprintf("%s/portfolio/positions", b.config.Endpoint)
	
	req, err := http.NewRequest("GET", positionsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get positions request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get positions request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read get positions response: %w", err)
	}
	
	var positionsResponse struct {
		Result []struct {
			Symbol          string  `json:"ExchangeInstrumentID"`
			Exchange        string  `json:"ExchangeSegment"`
			ProductType     string  `json:"ProductType"`
			Quantity        int     `json:"Quantity"`
			AveragePrice    float64 `json:"AveragePrice"`
			LastPrice       float64 `json:"LastPrice"`
			RealizedPnL     float64 `json:"RealizedPnL"`
			UnrealizedPnL   float64 `json:"UnrealizedPnL"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &positionsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal get positions response: %w", err)
	}
	
	if positionsResponse.Status != "Success" {
		return nil, fmt.Errorf("get positions failed: %s", positionsResponse.Message)
	}
	
	positions := make([]Position, 0, len(positionsResponse.Result))
	
	for _, xtsPosition := range positionsResponse.Result {
		// Map XTS-specific product type to our generic type
		var productType ProductType
		if xtsPosition.ProductType == "MIS" {
			productType = ProductTypeMIS
		} else {
			productType = ProductTypeNRML
		}
		
		positions = append(positions, Position{
			Symbol:          xtsPosition.Symbol,
			Exchange:        xtsPosition.Exchange,
			ProductType:     productType,
			Quantity:        xtsPosition.Quantity,
			AveragePrice:    xtsPosition.AveragePrice,
			LastPrice:       xtsPosition.LastPrice,
			PnL:             xtsPosition.RealizedPnL + xtsPosition.UnrealizedPnL,
			RealizedPnL:     xtsPosition.RealizedPnL,
			UnrealizedPnL:   xtsPosition.UnrealizedPnL,
			Timestamp:       time.Now(),
		})
	}
	
	return positions, nil
}

// GetQuote gets a quote for a symbol
func (b *XTSBroker) GetQuote(ctx context.Context, symbol, exchange string) (*Quote, error) {
	if !b.isConnected {
		return nil, fmt.Errorf("broker not connected")
	}
	
	quoteURL := fmt.Sprintf("%s/marketdata/quotes?exchangeSegment=%s&exchangeInstrumentID=%s", 
		b.config.Endpoint, url.QueryEscape(exchange), url.QueryEscape(symbol))
	
	req, err := http.NewRequest("GET", quoteURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get quote request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get quote request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read get quote response: %w", err)
	}
	
	var quoteResponse struct {
		Result struct {
			LastPrice     float64 `json:"LastPrice"`
			BidPrice      float64 `json:"BidPrice"`
			AskPrice      float64 `json:"AskPrice"`
			Volume        int     `json:"Volume"`
			OpenInterest  int     `json:"OpenInterest"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	
	err = json.Unmarshal(body, &quoteResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal get quote response: %w", err)
	}
	
	if quoteResponse.Status != "Success" {
		return nil, fmt.Errorf("get quote failed: %s", quoteResponse.Message)
	}
	
	return &Quote{
		Symbol:        symbol,
		Exchange:      exchange,
		LastPrice:     quoteResponse.Result.LastPrice,
		BidPrice:      quoteResponse.Result.BidPrice,
		AskPrice:      quoteResponse.Result.AskPrice,
		Volume:        quoteResponse.Result.Volume,
		OpenInterest:  quoteResponse.Result.OpenInterest,
		Timestamp:     time.Now(),
	}, nil
}

// SubscribeQuotes subscribes to quotes for symbols
func (b *XTSBroker) SubscribeQuotes(ctx context.Context, symbols []string, exchange string) error {
	// This would typically use WebSocket connection to subscribe to market data
	// For this implementation, we'll return not implemented
	return ErrNotImplemented
}

// UnsubscribeQuotes unsubscribes from quotes for symbols
func (b *XTSBroker) UnsubscribeQuotes(ctx context.Context, symbols []string, exchange string) error {
	// This would typically use WebSocket connection to unsubscribe from market data
	// For this implementation, we'll return not implemented
	return ErrNotImplemented
}

// Close closes the broker connection
func (b *XTSBroker) Close() error {
	if !b.isConnected {
		return nil
	}
	
	logoutURL := fmt.Sprintf("%s/user/logout", b.config.Endpoint)
	
	req, err := http.NewRequest("DELETE", logoutURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", b.accessToken)
	
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send logout request: %w", err)
	}
	defer resp.Body.Close()
	
	b.isConnected = false
	return nil
}
