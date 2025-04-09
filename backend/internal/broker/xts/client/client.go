// Package client provides the XTS Client implementation of the broker interface
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/trading-platform/backend/internal/broker/common"
)

// XTSClientImpl implements the BrokerClient interface for XTS Client
type XTSClientImpl struct {
	apiKey       string
	secretKey    string
	source       string
	token        string
	userID       string
	isInvestor   bool
	baseURL      string
	debug        bool
	httpClient   *http.Client
}

// NewXTSClientImpl creates a new XTS Client implementation
func NewXTSClientImpl(config *common.XTSClientConfig) (*XTSClientImpl, error) {
	if config == nil {
		return nil, errors.New("XTS Client configuration is required")
	}

	if config.APIKey == "" || config.SecretKey == "" {
		return nil, errors.New("API key and secret key are required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://developers.symphonyfintech.in" // Default URL
	}

	return &XTSClientImpl{
		apiKey:     config.APIKey,
		secretKey:  config.SecretKey,
		source:     config.Source,
		baseURL:    baseURL,
		debug:      false,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// SetDebug enables or disables debug mode
func (c *XTSClientImpl) SetDebug(debug bool) {
	c.debug = debug
}

// Login authenticates with the XTS Client API
func (c *XTSClientImpl) Login(credentials *common.Credentials) (*common.Session, error) {
	// Use credentials if provided, otherwise use the configured API key and secret key
	apiKey := c.apiKey
	secretKey := c.secretKey
	
	if credentials != nil && credentials.APIKey != "" && credentials.SecretKey != "" {
		apiKey = credentials.APIKey
		secretKey = credentials.SecretKey
	}

	params := url.Values{}
	params.Set("appKey", apiKey)
	params.Set("secretKey", secretKey)
	params.Set("source", c.source)

	url := fmt.Sprintf("%s/interactive/user/session", c.baseURL)
	
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var response struct {
		Type        string `json:"type"`
		Code        int    `json:"code"`
		Description string `json:"description"`
		Result      struct {
			Token            string `json:"token"`
			UserID           string `json:"userID"`
			IsInvestorClient bool   `json:"isInvestorClient"`
			ExpiresIn        int64  `json:"expiresIn"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("login failed: %s", response.Description)
	}
	
	// Store the token and user ID for future requests
	c.token = response.Result.Token
	c.userID = response.Result.UserID
	c.isInvestor = response.Result.IsInvestorClient
	
	// Calculate expiry time
	expiryTime := time.Now().Unix() + response.Result.ExpiresIn
	
	return &common.Session{
		Token:     response.Result.Token,
		UserID:    response.Result.UserID,
		ExpiresAt: expiryTime,
	}, nil
}

// Logout invalidates the current session
func (c *XTSClientImpl) Logout() error {
	if c.token == "" {
		return errors.New("not logged in")
	}
	
	url := fmt.Sprintf("%s/interactive/user/session", c.baseURL)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("Authorization", c.token)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var response struct {
		Type        string `json:"type"`
		Code        int    `json:"code"`
		Description string `json:"description"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return fmt.Errorf("logout failed: %s", response.Description)
	}
	
	// Clear the token and user ID
	c.token = ""
	c.userID = ""
	
	return nil
}

// GetOrderBook retrieves the order book for the specified client
func (c *XTSClientImpl) GetOrderBook(clientID string) (*common.OrderBook, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	url := fmt.Sprintf("%s/interactive/orders", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("Authorization", c.token)
	
	// Add clientID parameter if provided and not an investor client
	if clientID != "" && !c.isInvestor {
		q := req.URL.Query()
		q.Add("clientID", clientID)
		req.URL.RawQuery = q.Encode()
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var response struct {
		Type        string `json:"type"`
		Code        int    `json:"code"`
		Description string `json:"description"`
		Result      []struct {
			OrderID              string  `json:"OrderID"`
			ExchangeOrderID      string  `json:"ExchangeOrderID"`
			ExchangeSegment      string  `json:"ExchangeSegment"`
			ExchangeInstrumentID string  `json:"ExchangeInstrumentID"`
			OrderSide            string  `json:"OrderSide"`
			OrderType            string  `json:"OrderType"`
			ProductType          string  `json:"ProductType"`
			TimeInForce          string  `json:"TimeInForce"`
			OrderQuantity        int     `json:"OrderQuantity"`
			FilledQuantity       int     `json:"FilledQuantity"`
			RemainingQuantity    int     `json:"RemainingQuantity"`
			LimitPrice           float64 `json:"LimitPrice"`
			StopPrice            float64 `json:"StopPrice"`
			OrderStatus          string  `json:"OrderStatus"`
			OrderTimestamp       int64   `json:"OrderTimestamp"`
			LastUpdateTimestamp  int64   `json:"LastUpdateTimestamp"`
			CancelTimestamp      int64   `json:"CancelTimestamp"`
			ClientID             string  `json:"ClientID"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("get order book failed: %s", response.Description)
	}
	
	// Convert the response to the common OrderBook model
	orderBook := &common.OrderBook{
		Orders: make([]common.OrderDetails, len(response.Result)),
	}
	
	for i, order := range response.Result {
		orderBook.Orders[i] = common.OrderDetails{
			OrderID:              order.OrderID,
			ExchangeOrderID:      order.ExchangeOrderID,
			ExchangeSegment:      order.ExchangeSegment,
			ExchangeInstrumentID: order.ExchangeInstrumentID,
			OrderSide:            order.OrderSide,
			OrderType:            order.OrderType,
			ProductType:          order.ProductType,
			TimeInForce:          order.TimeInForce,
			OrderQuantity:        order.OrderQuantity,
			FilledQuantity:       order.FilledQuantity,
			RemainingQuantity:    order.RemainingQuantity,
			LimitPrice:           order.LimitPrice,
			StopPrice:            order.StopPrice,
			OrderStatus:          order.OrderStatus,
			OrderTimestamp:       order.OrderTimestamp,
			LastUpdateTimestamp:  order.LastUpdateTimestamp,
			CancelTimestamp:      order.CancelTimestamp,
			ClientID:             order.ClientID,
		}
	}
	
	return orderBook, nil
}

// Placeholder implementations for other interface methods
// These will be implemented in subsequent steps

func (c *XTSClientImpl) PlaceOrder(order *common.Order) (*common.OrderResponse, error) {
	// Will be implemented in step 002
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) ModifyOrder(order *common.ModifyOrder) (*common.OrderResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) CancelOrder(orderID string, clientID string) (*common.OrderResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) GetPositions(clientID string) ([]common.Position, error) {
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) GetHoldings(clientID string) ([]common.Holding, error) {
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) GetQuote(symbols []string) (map[string]common.Quote, error) {
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) SubscribeToQuotes(symbols []string) (chan common.Quote, error) {
	return nil, errors.New("not implemented")
}

func (c *XTSClientImpl) UnsubscribeFromQuotes(symbols []string) error {
	return errors.New("not implemented")
}
