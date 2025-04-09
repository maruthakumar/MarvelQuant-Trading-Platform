// Package client provides the XTS Client implementation of the broker interface
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/trading-platform/backend/internal/broker/common"
)

// ModifyOrder modifies an existing order with the XTS Client API
func (c *XTSClientImpl) ModifyOrder(order *common.ModifyOrder) (*common.OrderResponse, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if order == nil {
		return nil, errors.New("order is required")
	}
	
	if order.OrderID == "" {
		return nil, errors.New("order ID is required")
	}
	
	url := fmt.Sprintf("%s/interactive/orders", c.baseURL)
	
	// Prepare the request parameters
	params := url.Values{}
	params.Set("appOrderID", order.OrderID)
	params.Set("orderQuantity", fmt.Sprintf("%d", order.OrderQuantity))
	params.Set("orderType", order.OrderType)
	params.Set("limitPrice", fmt.Sprintf("%f", order.LimitPrice))
	params.Set("stopPrice", fmt.Sprintf("%f", order.StopPrice))
	
	// Add clientID parameter if provided and not an investor client
	if order.ClientID != "" && !c.isInvestor {
		params.Set("clientID", order.ClientID)
	}
	
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.token)
	
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
			AppOrderID       string `json:"AppOrderID"`
			OrderGeneratedID string `json:"OrderGeneratedID"`
			OrderStatus      string `json:"OrderStatus"`
			OrderRejected    bool   `json:"OrderRejected"`
			RejectReason     string `json:"RejectReason"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("modify order failed: %s", response.Description)
	}
	
	// Convert the response to the common OrderResponse model
	orderResponse := &common.OrderResponse{
		OrderID:         response.Result.AppOrderID,
		ExchangeOrderID: response.Result.OrderGeneratedID,
		Status:          response.Result.OrderStatus,
	}
	
	if response.Result.OrderRejected {
		orderResponse.Status = "REJECTED"
		orderResponse.RejectionReason = response.Result.RejectReason
	}
	
	return orderResponse, nil
}

// CancelOrder cancels an existing order with the XTS Client API
func (c *XTSClientImpl) CancelOrder(orderID string, clientID string) (*common.OrderResponse, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if orderID == "" {
		return nil, errors.New("order ID is required")
	}
	
	url := fmt.Sprintf("%s/interactive/orders", c.baseURL)
	
	// Prepare the request parameters
	params := url.Values{}
	params.Set("appOrderID", orderID)
	
	// Add clientID parameter if provided and not an investor client
	if clientID != "" && !c.isInvestor {
		params.Set("clientID", clientID)
	}
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.token)
	
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
			AppOrderID       string `json:"AppOrderID"`
			OrderGeneratedID string `json:"OrderGeneratedID"`
			OrderStatus      string `json:"OrderStatus"`
			OrderRejected    bool   `json:"OrderRejected"`
			RejectReason     string `json:"RejectReason"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("cancel order failed: %s", response.Description)
	}
	
	// Convert the response to the common OrderResponse model
	orderResponse := &common.OrderResponse{
		OrderID:         response.Result.AppOrderID,
		ExchangeOrderID: response.Result.OrderGeneratedID,
		Status:          response.Result.OrderStatus,
	}
	
	if response.Result.OrderRejected {
		orderResponse.Status = "REJECTED"
		orderResponse.RejectionReason = response.Result.RejectReason
	}
	
	return orderResponse, nil
}

// GetPositions retrieves the positions for the specified client
func (c *XTSClientImpl) GetPositions(clientID string) ([]common.Position, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	url := fmt.Sprintf("%s/interactive/portfolio/positions", c.baseURL)
	
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
		return nil, fmt.Errorf("get positions failed: %s", response.Description)
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

// GetHoldings retrieves the holdings for the specified client
func (c *XTSClientImpl) GetHoldings(clientID string) ([]common.Holding, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	url := fmt.Sprintf("%s/interactive/portfolio/holdings", c.baseURL)
	
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
			ExchangeSegment      string  `json:"ExchangeSegment"`
			ExchangeInstrumentID string  `json:"ExchangeInstrumentID"`
			TradingSymbol        string  `json:"TradingSymbol"`
			ISIN                 string  `json:"ISIN"`
			Quantity             int     `json:"Quantity"`
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
		return nil, fmt.Errorf("get holdings failed: %s", response.Description)
	}
	
	// Convert the response to the common Holding model
	holdings := make([]common.Holding, len(response.Result))
	
	for i, holding := range response.Result {
		holdings[i] = common.Holding{
			ExchangeSegment:      holding.ExchangeSegment,
			ExchangeInstrumentID: holding.ExchangeInstrumentID,
			TradingSymbol:        holding.TradingSymbol,
			ISIN:                 holding.ISIN,
			Quantity:             holding.Quantity,
			AveragePrice:         holding.AveragePrice,
			LastPrice:            holding.LastPrice,
			RealizedProfit:       holding.RealizedProfit,
			UnrealizedProfit:     holding.UnrealizedProfit,
			ClientID:             holding.ClientID,
		}
	}
	
	return holdings, nil
}

// GetQuote retrieves quotes for the specified symbols
func (c *XTSClientImpl) GetQuote(symbols []string) (map[string]common.Quote, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if len(symbols) == 0 {
		return nil, errors.New("at least one symbol is required")
	}
	
	url := fmt.Sprintf("%s/marketdata/instruments/quotes", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Add("Authorization", c.token)
	
	// Add symbols parameter
	q := req.URL.Query()
	q.Add("instruments", strings.Join(symbols, ","))
	req.URL.RawQuery = q.Encode()
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	var response struct {
		Type        string `json:"type"`
		Code        int    `json:"code"`
		Description string `json:"description"`
		Result      map[string]struct {
			ExchangeSegment      string  `json:"ExchangeSegment"`
			ExchangeInstrumentID string  `json:"ExchangeInstrumentID"`
			TradingSymbol        string  `json:"TradingSymbol"`
			LastPrice            float64 `json:"LastPrice"`
			Open                 float64 `json:"Open"`
			High                 float64 `json:"High"`
			Low                  float64 `json:"Low"`
			Close                float64 `json:"Close"`
			Volume               int64   `json:"Volume"`
			BidPrice             float64 `json:"BidPrice"`
			BidSize              int     `json:"BidSize"`
			AskPrice             float64 `json:"AskPrice"`
			AskSize              int     `json:"AskSize"`
			Timestamp            int64   `json:"Timestamp"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Type != "success" {
		return nil, fmt.Errorf("get quote failed: %s", response.Description)
	}
	
	// Convert the response to the common Quote model
	quotes := make(map[string]common.Quote, len(response.Result))
	
	for symbol, quote := range response.Result {
		quotes[symbol] = common.Quote{
			ExchangeSegment:      quote.ExchangeSegment,
			ExchangeInstrumentID: quote.ExchangeInstrumentID,
			TradingSymbol:        quote.TradingSymbol,
			LastPrice:            quote.LastPrice,
			Open:                 quote.Open,
			High:                 quote.High,
			Low:                  quote.Low,
			Close:                quote.Close,
			Volume:               quote.Volume,
			BidPrice:             quote.BidPrice,
			BidSize:              quote.BidSize,
			AskPrice:             quote.AskPrice,
			AskSize:              quote.AskSize,
			Timestamp:            quote.Timestamp,
		}
	}
	
	return quotes, nil
}

// SubscribeToQuotes subscribes to real-time quotes for the specified symbols
func (c *XTSClientImpl) SubscribeToQuotes(symbols []string) (chan common.Quote, error) {
	if c.token == "" {
		return nil, errors.New("not logged in")
	}
	
	if len(symbols) == 0 {
		return nil, errors.New("at least one symbol is required")
	}
	
	// Create a channel for quotes
	quoteChan := make(chan common.Quote, 100)
	
	// In a real implementation, this would establish a WebSocket connection
	// and stream quotes to the channel. For now, we'll return an error
	// indicating that this is not yet implemented.
	
	// TODO: Implement WebSocket connection for real-time quotes
	
	return quoteChan, errors.New("real-time quotes not yet implemented")
}

// UnsubscribeFromQuotes unsubscribes from real-time quotes for the specified symbols
func (c *XTSClientImpl) UnsubscribeFromQuotes(symbols []string) error {
	if c.token == "" {
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
