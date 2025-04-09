package orderexecution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// XTSBrokerAdapter implements the BrokerAdapter interface for XTS
type XTSBrokerAdapter struct {
	baseURL       string
	apiKey        string
	apiSecret     string
	clientCode    string
	authToken     string
	tokenExpiry   time.Time
	mutex         sync.RWMutex
	httpClient    *http.Client
	isInteractive bool // Interactive or MarketData API
}

// XTSAuthResponse represents the authentication response from XTS
type XTSAuthResponse struct {
	Type             string `json:"type"`
	Result           bool   `json:"result"`
	Description      string `json:"description"`
	InfoID           int    `json:"infoID"`
	InfoMsg          string `json:"infoMsg"`
	Token            string `json:"token"`
	UserID           string `json:"userID"`
	Enabled          bool   `json:"enabled"`
	AccountID        string `json:"accountID"`
	AccountName      string `json:"accountName"`
	EmailID          string `json:"emailID"`
	MobileNumber     string `json:"mobileNumber"`
	ExchangeSegments []struct {
		ExchSegID   int    `json:"exchSegId"`
		ExchSegName string `json:"exchSegName"`
	} `json:"exchangeSegments"`
	OrderTypes []struct {
		OrderTypeID   int    `json:"orderTypeId"`
		OrderTypeName string `json:"orderTypeName"`
	} `json:"orderTypes"`
	ProductTypes []struct {
		ProductTypeID   int    `json:"productTypeId"`
		ProductTypeName string `json:"productTypeName"`
	} `json:"productTypes"`
}

// XTSOrderResponse represents the order response from XTS
type XTSOrderResponse struct {
	Type        string `json:"type"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
	InfoID      int    `json:"infoID"`
	InfoMsg     string `json:"infoMsg"`
	OrderID     string `json:"OrderUniqueIdentifier"`
	AppOrderID  string `json:"AppOrderID"`
}

// XTSOrderStatusResponse represents the order status response from XTS
type XTSOrderStatusResponse struct {
	Type        string `json:"type"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
	InfoID      int    `json:"infoID"`
	InfoMsg     string `json:"infoMsg"`
	OrderStatus struct {
		LoginID                 string  `json:"LoginId"`
		ClientID                string  `json:"ClientId"`
		AppOrderID              string  `json:"AppOrderID"`
		OrderReferenceID        string  `json:"OrderReferenceID"`
		GeneratedBy             string  `json:"GeneratedBy"`
		ExchangeOrderID         string  `json:"ExchangeOrderID"`
		OrderCategoryType       string  `json:"OrderCategoryType"`
		ExchangeSegment         string  `json:"ExchangeSegment"`
		ExchangeInstrumentID    string  `json:"ExchangeInstrumentID"`
		OrderSide               string  `json:"OrderSide"`
		OrderType               string  `json:"OrderType"`
		ProductType             string  `json:"ProductType"`
		TimeInForce             string  `json:"TimeInForce"`
		OrderPrice              float64 `json:"OrderPrice"`
		OrderQuantity           int     `json:"OrderQuantity"`
		OrderStopPrice          float64 `json:"OrderStopPrice"`
		OrderStatus             string  `json:"OrderStatus"`
		OrderAverageTradedPrice float64 `json:"OrderAverageTradedPrice"`
		LeavesQuantity          int     `json:"LeavesQuantity"`
		CumulativeQuantity      int     `json:"CumulativeQuantity"`
		OrderDisclosedQuantity  int     `json:"OrderDisclosedQuantity"`
		OrderGeneratedDateTime  string  `json:"OrderGeneratedDateTime"`
		ExchangeTransactTime    string  `json:"ExchangeTransactTime"`
		LastUpdateDateTime      string  `json:"LastUpdateDateTime"`
		OrderExpiryDate         string  `json:"OrderExpiryDate"`
		CancelRejectReason      string  `json:"CancelRejectReason"`
		OrderUniqueIdentifier   string  `json:"OrderUniqueIdentifier"`
	} `json:"OrderStatus"`
}

// XTSOrdersResponse represents the orders response from XTS
type XTSOrdersResponse struct {
	Type        string `json:"type"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
	InfoID      int    `json:"infoID"`
	InfoMsg     string `json:"infoMsg"`
	OrderBook   []struct {
		LoginID                 string  `json:"LoginId"`
		ClientID                string  `json:"ClientId"`
		AppOrderID              string  `json:"AppOrderID"`
		OrderReferenceID        string  `json:"OrderReferenceID"`
		GeneratedBy             string  `json:"GeneratedBy"`
		ExchangeOrderID         string  `json:"ExchangeOrderID"`
		OrderCategoryType       string  `json:"OrderCategoryType"`
		ExchangeSegment         string  `json:"ExchangeSegment"`
		ExchangeInstrumentID    string  `json:"ExchangeInstrumentID"`
		OrderSide               string  `json:"OrderSide"`
		OrderType               string  `json:"OrderType"`
		ProductType             string  `json:"ProductType"`
		TimeInForce             string  `json:"TimeInForce"`
		OrderPrice              float64 `json:"OrderPrice"`
		OrderQuantity           int     `json:"OrderQuantity"`
		OrderStopPrice          float64 `json:"OrderStopPrice"`
		OrderStatus             string  `json:"OrderStatus"`
		OrderAverageTradedPrice float64 `json:"OrderAverageTradedPrice"`
		LeavesQuantity          int     `json:"LeavesQuantity"`
		CumulativeQuantity      int     `json:"CumulativeQuantity"`
		OrderDisclosedQuantity  int     `json:"OrderDisclosedQuantity"`
		OrderGeneratedDateTime  string  `json:"OrderGeneratedDateTime"`
		ExchangeTransactTime    string  `json:"ExchangeTransactTime"`
		LastUpdateDateTime      string  `json:"LastUpdateDateTime"`
		OrderExpiryDate         string  `json:"OrderExpiryDate"`
		CancelRejectReason      string  `json:"CancelRejectReason"`
		OrderUniqueIdentifier   string  `json:"OrderUniqueIdentifier"`
	} `json:"OrderBook"`
}

// NewXTSBrokerAdapter creates a new XTS broker adapter
func NewXTSBrokerAdapter(baseURL, apiKey, apiSecret, clientCode string, isInteractive bool) *XTSBrokerAdapter {
	return &XTSBrokerAdapter{
		baseURL:       baseURL,
		apiKey:        apiKey,
		apiSecret:     apiSecret,
		clientCode:    clientCode,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		isInteractive: isInteractive,
	}
}

// authenticate authenticates with the XTS API
func (a *XTSBrokerAdapter) authenticate(ctx context.Context) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if token is still valid
	if a.authToken != "" && time.Now().Before(a.tokenExpiry) {
		return nil
	}

	// Prepare request
	url := a.baseURL
	if a.isInteractive {
		url += "/interactive/user/session"
	} else {
		url += "/marketdata/auth/login"
	}

	payload := map[string]string{
		"appKey":    a.apiKey,
		"secretKey": a.apiSecret,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse response
	var authResp XTSAuthResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return err
	}

	if !authResp.Result {
		return fmt.Errorf("authentication failed: %s", authResp.Description)
	}

	// Store token
	a.authToken = authResp.Token
	a.tokenExpiry = time.Now().Add(8 * time.Hour) // Tokens are valid for 8 hours

	return nil
}

// mapOrderType maps our order type to XTS order type
func (a *XTSBrokerAdapter) mapOrderType(orderType OrderType) string {
	switch orderType {
	case Market:
		return "MARKET"
	case Limit:
		return "LIMIT"
	case StopLoss:
		return "STOPLOSS_LIMIT"
	case StopLossMarket:
		return "STOPLOSS_MARKET"
	default:
		return "LIMIT"
	}
}

// mapTransactionType maps our transaction type to XTS transaction type
func (a *XTSBrokerAdapter) mapTransactionType(transactionType TransactionType) string {
	switch transactionType {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return "BUY"
	}
}

// mapProductType maps our product type to XTS product type
func (a *XTSBrokerAdapter) mapProductType(productType ProductType) string {
	switch productType {
	case Intraday:
		return "MIS"
	case Normal:
		return "NRML"
	case CashAndCarry:
		return "CNC"
	default:
		return "NRML"
	}
}

// mapValidityType maps our validity type to XTS validity type
func (a *XTSBrokerAdapter) mapValidityType(validityType ValidityType) string {
	switch validityType {
	case Day:
		return "DAY"
	case IOC:
		return "IOC"
	case GTC:
		return "GTC"
	default:
		return "DAY"
	}
}

// mapXTSOrderStatus maps XTS order status to our order status
func (a *XTSBrokerAdapter) mapXTSOrderStatus(xtsStatus string) OrderStatus {
	switch xtsStatus {
	case "Filled":
		return Executed
	case "PartiallyFilled":
		return PartiallyExecuted
	case "Cancelled":
		return Cancelled
	case "Rejected":
		return Rejected
	case "New":
		return Open
	case "PendingNew":
		return Pending
	default:
		return Open
	}
}

// PlaceOrder places an order with XTS
func (a *XTSBrokerAdapter) PlaceOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	// Authenticate if needed
	err := a.authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare request
	url := a.baseURL + "/interactive/orders"

	// Map order parameters
	orderSide := a.mapTransactionType(request.TransactionType)
	orderType := a.mapOrderType(request.OrderType)
	productType := a.mapProductType(request.Product)
	timeInForce := a.mapValidityType(request.Validity)

	// Create payload
	payload := map[string]interface{}{
		"exchangeSegment":      request.Exchange,
		"exchangeInstrumentID": request.Symbol,
		"orderType":            orderType,
		"orderSide":            orderSide,
		"orderQuantity":        request.Quantity,
		"productType":          productType,
		"timeInForce":          timeInForce,
		"disclosedQuantity":    request.Quantity,
		"orderPrice":           request.Price,
		"clientID":             a.clientCode,
	}

	// Add trigger price for stop orders
	if request.OrderType == StopLoss || request.OrderType == StopLossMarket {
		payload["stopPrice"] = request.TriggerPrice
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", a.authToken)

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var orderResp XTSOrderResponse
	err = json.Unmarshal(body, &orderResp)
	if err != nil {
		return nil, err
	}

	if !orderResp.Result {
		return nil, fmt.Errorf("order placement failed: %s", orderResp.Description)
	}

	// Create order response
	order := &Order{
		ID:              orderResp.OrderID,
		Symbol:          request.Symbol,
		Quantity:        request.Quantity,
		Price:           request.Price,
		OrderType:       request.OrderType,
		TransactionType: request.TransactionType,
		Status:          Pending,
		FilledQuantity:  0,
		AveragePrice:    0,
		PlacedAt:        time.Now(),
		UpdatedAt:       time.Now(),
		Validity:        request.Validity,
		TriggerPrice:    request.TriggerPrice,
		Exchange:        request.Exchange,
		Product:         request.Product,
		BrokerOrderID:   orderResp.OrderID,
		StrategyID:      request.StrategyID,
		Tags:            request.Tags,
	}

	return &OrderResponse{
		Order:  order,
		Status: true,
	}, nil
}

// ModifyOrder modifies an existing order with XTS
func (a *XTSBrokerAdapter) ModifyOrder(ctx context.Context, orderID string, request *OrderRequest) (*OrderResponse, error) {
	// Authenticate if needed
	err := a.authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare request
	url := a.baseURL + "/interactive/orders"

	// Map order parameters
	orderType := a.mapOrderType(request.OrderType)
	timeInForce := a.mapValidityType(request.Validity)

	// Create payload
	payload := map[string]interface{}{
		"appOrderID":        orderID,
		"modifiedOrderType": orderType,
		"modifiedOrderQuantity": request.Quantity,
		"modifiedDisclosedQuantity": request.Quantity,
		"modifiedOrderPrice": request.Price,
		"modifiedTimeInForce": timeInForce,
		"clientID": a.clientCode,
	}

	// Add trigger price for stop orders
	if request.OrderType == StopLoss || request.OrderType == StopLossMarket {
		payload["modifiedStopPrice"] = request.TriggerPrice
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", a.authToken)

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var orderResp XTSOrderResponse
	err = json.Unmarshal(body, &orderResp)
	if err != nil {
		return nil, err
	}

	if !orderResp.Result {
		return nil, fmt.Errorf("order modification failed: %s", orderResp.Description)
	}

	// Get updated order status
	return a.GetOrderStatus(ctx, orderID)
}

// CancelOrder cancels an existing order with XTS
func (a *XTSBrokerAdapter) CancelOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	// Authenticate if needed
	err := a.authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare request
	url := a.baseURL + "/interactive/orders"

	// Create payload
	payload := map[string]interface{}{
		"appOrderID": orderID,
		"clientID":   a.clientCode,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", a.authToken)

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var orderResp XTSOrderResponse
	err = json.Unmarshal(body, &orderResp)
	if err != nil {
		return nil, err
	}

	if !orderResp.Result {
		return nil, fmt.Errorf("order cancellation failed: %s", orderResp.Description)
	}

	// Get updated order status
	return a.GetOrderStatus(ctx, orderID)
}

// GetOrderStatus gets the status of an order from XTS
func (a *XTSBrokerAdapter) GetOrderStatus(ctx context.Context, orderID string) (*Order, error) {
	// Authenticate if needed
	err := a.authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare request
	url := fmt.Sprintf("%s/interactive/orders?appOrderID=%s&clientID=%s", a.baseURL, orderID, a.clientCode)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", a.authToken)

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var statusResp XTSOrderStatusResponse
	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, err
	}

	if !statusResp.Result {
		return nil, fmt.Errorf("get order status failed: %s", statusResp.Description)
	}

	// Map XTS order status to our order status
	status := a.mapXTSOrderStatus(statusResp.OrderStatus.OrderStatus)

	// Create order
	order := &Order{
		ID:              statusResp.OrderStatus.OrderUniqueIdentifier,
		Symbol:          statusResp.OrderStatus.ExchangeInstrumentID,
		Quantity:        statusResp.OrderStatus.OrderQuantity,
		Price:           statusResp.OrderStatus.OrderPrice,
		OrderType:       a.mapXTSOrderType(statusResp.OrderStatus.OrderType),
		TransactionType: a.mapXTSTransactionType(statusResp.OrderStatus.OrderSide),
		Status:          status,
		FilledQuantity:  statusResp.OrderStatus.CumulativeQuantity,
		AveragePrice:    statusResp.OrderStatus.OrderAverageTradedPrice,
		PlacedAt:        a.parseXTSDateTime(statusResp.OrderStatus.OrderGeneratedDateTime),
		UpdatedAt:       a.parseXTSDateTime(statusResp.OrderStatus.LastUpdateDateTime),
		Validity:        a.mapXTSValidityType(statusResp.OrderStatus.TimeInForce),
		TriggerPrice:    statusResp.OrderStatus.OrderStopPrice,
		Exchange:        statusResp.OrderStatus.ExchangeSegment,
		Product:         a.mapXTSProductType(statusResp.OrderStatus.ProductType),
		BrokerOrderID:   statusResp.OrderStatus.OrderUniqueIdentifier,
		Message:         statusResp.OrderStatus.CancelRejectReason,
	}

	return order, nil
}

// GetOrders gets all orders from XTS
func (a *XTSBrokerAdapter) GetOrders(ctx context.Context) ([]*Order, error) {
	// Authenticate if needed
	err := a.authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare request
	url := fmt.Sprintf("%s/interactive/orders?clientID=%s", a.baseURL, a.clientCode)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", a.authToken)

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var ordersResp XTSOrdersResponse
	err = json.Unmarshal(body, &ordersResp)
	if err != nil {
		return nil, err
	}

	if !ordersResp.Result {
		return nil, fmt.Errorf("get orders failed: %s", ordersResp.Description)
	}

	// Create orders
	orders := make([]*Order, len(ordersResp.OrderBook))
	for i, xtsOrder := range ordersResp.OrderBook {
		// Map XTS order status to our order status
		status := a.mapXTSOrderStatus(xtsOrder.OrderStatus)

		orders[i] = &Order{
			ID:              xtsOrder.OrderUniqueIdentifier,
			Symbol:          xtsOrder.ExchangeInstrumentID,
			Quantity:        xtsOrder.OrderQuantity,
			Price:           xtsOrder.OrderPrice,
			OrderType:       a.mapXTSOrderType(xtsOrder.OrderType),
			TransactionType: a.mapXTSTransactionType(xtsOrder.OrderSide),
			Status:          status,
			FilledQuantity:  xtsOrder.CumulativeQuantity,
			AveragePrice:    xtsOrder.OrderAverageTradedPrice,
			PlacedAt:        a.parseXTSDateTime(xtsOrder.OrderGeneratedDateTime),
			UpdatedAt:       a.parseXTSDateTime(xtsOrder.LastUpdateDateTime),
			Validity:        a.mapXTSValidityType(xtsOrder.TimeInForce),
			TriggerPrice:    xtsOrder.OrderStopPrice,
			Exchange:        xtsOrder.ExchangeSegment,
			Product:         a.mapXTSProductType(xtsOrder.ProductType),
			BrokerOrderID:   xtsOrder.OrderUniqueIdentifier,
			Message:         xtsOrder.CancelRejectReason,
		}
	}

	return orders, nil
}

// mapXTSOrderType maps XTS order type to our order type
func (a *XTSBrokerAdapter) mapXTSOrderType(xtsOrderType string) OrderType {
	switch xtsOrderType {
	case "MARKET":
		return Market
	case "LIMIT":
		return Limit
	case "STOPLOSS_LIMIT":
		return StopLoss
	case "STOPLOSS_MARKET":
		return StopLossMarket
	default:
		return Limit
	}
}

// mapXTSTransactionType maps XTS transaction type to our transaction type
func (a *XTSBrokerAdapter) mapXTSTransactionType(xtsTransactionType string) TransactionType {
	switch xtsTransactionType {
	case "BUY":
		return Buy
	case "SELL":
		return Sell
	default:
		return Buy
	}
}

// mapXTSProductType maps XTS product type to our product type
func (a *XTSBrokerAdapter) mapXTSProductType(xtsProductType string) ProductType {
	switch xtsProductType {
	case "MIS":
		return Intraday
	case "NRML":
		return Normal
	case "CNC":
		return CashAndCarry
	default:
		return Normal
	}
}

// mapXTSValidityType maps XTS validity type to our validity type
func (a *XTSBrokerAdapter) mapXTSValidityType(xtsValidityType string) ValidityType {
	switch xtsValidityType {
	case "DAY":
		return Day
	case "IOC":
		return IOC
	case "GTC":
		return GTC
	default:
		return Day
	}
}

// parseXTSDateTime parses XTS date time string
func (a *XTSBrokerAdapter) parseXTSDateTime(dateTimeStr string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05", dateTimeStr)
	if err != nil {
		return time.Now()
	}
	return t
}

// BrokerFactory creates broker adapters
type BrokerFactory struct {
	brokers map[string]func(config map[string]string) (BrokerAdapter, error)
}

// NewBrokerFactory creates a new broker factory
func NewBrokerFactory() *BrokerFactory {
	factory := &BrokerFactory{
		brokers: make(map[string]func(config map[string]string) (BrokerAdapter, error)),
	}

	// Register XTS broker
	factory.RegisterBroker("XTS", func(config map[string]string) (BrokerAdapter, error) {
		baseURL, ok := config["baseURL"]
		if !ok {
			return nil, errors.New("baseURL is required for XTS broker")
		}

		apiKey, ok := config["apiKey"]
		if !ok {
			return nil, errors.New("apiKey is required for XTS broker")
		}

		apiSecret, ok := config["apiSecret"]
		if !ok {
			return nil, errors.New("apiSecret is required for XTS broker")
		}

		clientCode, ok := config["clientCode"]
		if !ok {
			return nil, errors.New("clientCode is required for XTS broker")
		}

		isInteractiveStr, ok := config["isInteractive"]
		isInteractive := true
		if ok {
			isInteractive = isInteractiveStr == "true"
		}

		return NewXTSBrokerAdapter(baseURL, apiKey, apiSecret, clientCode, isInteractive), nil
	})

	return factory
}

// RegisterBroker registers a broker with the factory
func (f *BrokerFactory) RegisterBroker(name string, creator func(config map[string]string) (BrokerAdapter, error)) {
	f.brokers[name] = creator
}

// CreateBroker creates a broker adapter
func (f *BrokerFactory) CreateBroker(name string, config map[string]string) (BrokerAdapter, error) {
	creator, ok := f.brokers[name]
	if !ok {
		return nil, fmt.Errorf("broker not found: %s", name)
	}

	return creator(config)
}

// ListBrokers returns a list of available broker names
func (f *BrokerFactory) ListBrokers() []string {
	names := make([]string, 0, len(f.brokers))
	for name := range f.brokers {
		names = append(names, name)
	}
	return names
}
