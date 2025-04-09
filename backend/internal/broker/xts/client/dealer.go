// Package client provides dealer-specific endpoints for the XTS Client implementation
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/trading-platform/backend/internal/broker/common"
)

// GetDealerOrderBook retrieves the dealer order book for the specified client
// This is a dealer-specific endpoint that is only available in the XTS Client API
func (c *XTSClientImpl) GetDealerOrderBook(clientID string) (*common.OrderBook, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if c.isInvestor {
		return nil, errors.New("dealer endpoints are not available for investor clients")
	}
	
	url := fmt.Sprintf("%s/interactive/orders/dealerorderbook", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("Authorization", c.token)
	
	// Add clientID parameter if provided
	if clientID != "" {
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
		return nil, fmt.Errorf("get dealer order book failed: %s", response.Description)
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

// GetDealerTrades retrieves the dealer trades for the specified client
// This is a dealer-specific endpoint that is only available in the XTS Client API
func (c *XTSClientImpl) GetDealerTrades(clientID string) ([]common.OrderDetails, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if c.isInvestor {
		return nil, errors.New("dealer endpoints are not available for investor clients")
	}
	
	url := fmt.Sprintf("%s/interactive/orders/dealertradebook", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("Authorization", c.token)
	
	// Add clientID parameter if provided
	if clientID != "" {
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
			TradeID              string  `json:"TradeID"`
			TradePrice           float64 `json:"TradePrice"`
			TradeQuantity        int     `json:"TradeQuantity"`
			TradeTimestamp       int64   `json:"TradeTimestamp"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("get dealer trades failed: %s", response.Description)
	}
	
	// Convert the response to the common OrderDetails model
	trades := make([]common.OrderDetails, len(response.Result))
	
	for i, trade := range response.Result {
		trades[i] = common.OrderDetails{
			OrderID:              trade.OrderID,
			ExchangeOrderID:      trade.ExchangeOrderID,
			ExchangeSegment:      trade.ExchangeSegment,
			ExchangeInstrumentID: trade.ExchangeInstrumentID,
			OrderSide:            trade.OrderSide,
			OrderType:            trade.OrderType,
			ProductType:          trade.ProductType,
			TimeInForce:          trade.TimeInForce,
			OrderQuantity:        trade.OrderQuantity,
			FilledQuantity:       trade.FilledQuantity,
			RemainingQuantity:    trade.RemainingQuantity,
			LimitPrice:           trade.LimitPrice,
			StopPrice:            trade.StopPrice,
			OrderStatus:          trade.OrderStatus,
			OrderTimestamp:       trade.OrderTimestamp,
			LastUpdateTimestamp:  trade.LastUpdateTimestamp,
			CancelTimestamp:      trade.CancelTimestamp,
			ClientID:             trade.ClientID,
		}
	}
	
	return trades, nil
}

// GetDealerPositions retrieves the dealer positions for the specified client
// This is a dealer-specific endpoint that is only available in the XTS Client API
func (c *XTSClientImpl) GetDealerPositions(clientID string) ([]common.Position, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if c.isInvestor {
		return nil, errors.New("dealer endpoints are not available for investor clients")
	}
	
	url := fmt.Sprintf("%s/interactive/portfolio/dealerpositions", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("Authorization", c.token)
	
	// Add clientID parameter if provided
	if clientID != "" {
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
			ExchangeSegment      string  `json:"ExchangeSegment"`
			ExchangeInstrumentID string  `json:"ExchangeInstrumentID"`
			ProductType          string  `json:"ProductType"`
			Quantity             int     `json:"Quantity"`
			BuyQuantity          int     `json:"BuyQuantity"`
			SellQuantity         int     `json:"SellQuantity"`
			NetQuantity          int     `json:"NetQuantity"`
			AveragePrice         float64 `json:"AveragePrice"`
			LastPrice            float64 `json:"LastPrice"`
			RealizedProfit       float64 `json:"RealizedProfit"`
			UnrealizedProfit     float64 `json:"UnrealizedProfit"`
			ClientID             string  `json:"ClientID"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("get dealer positions failed: %s", response.Description)
	}
	
	// Convert the response to the common Position model
	positions := make([]common.Position, len(response.Result))
	
	for i, position := range response.Result {
		positions[i] = common.Position{
			ExchangeSegment:      position.ExchangeSegment,
			ExchangeInstrumentID: position.ExchangeInstrumentID,
			ProductType:          position.ProductType,
			Quantity:             position.Quantity,
			BuyQuantity:          position.BuyQuantity,
			SellQuantity:         position.SellQuantity,
			NetQuantity:          position.NetQuantity,
			AveragePrice:         position.AveragePrice,
			LastPrice:            position.LastPrice,
			RealizedProfit:       position.RealizedProfit,
			UnrealizedProfit:     position.UnrealizedProfit,
			ClientID:             position.ClientID,
		}
	}
	
	return positions, nil
}

// PlaceDealerOrder places an order on behalf of a client as a dealer
// This is a dealer-specific endpoint that is only available in the XTS Client API
func (c *XTSClientImpl) PlaceDealerOrder(order *common.Order) (*common.OrderResponse, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if c.isInvestor {
		return nil, errors.New("dealer endpoints are not available for investor clients")
	}
	
	if order == nil {
		return nil, errors.New("order is required")
	}
	
	if order.ClientID == "" {
		return nil, errors.New("client ID is required for dealer orders")
	}
	
	// Use the standard PlaceOrder method but ensure clientID is set
	return c.PlaceOrder(order)
}
