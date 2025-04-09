package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// XTSConnector implements DataSourceConnector for XTS API
type XTSConnector struct {
	apiKey      string
	secretKey   string
	source      string
	userID      string
	isConnected bool
	baseURL     string
	token       string
	httpClient  *http.Client
	mutex       sync.RWMutex
	callbacks   map[string][]MarketDataCallback
	callbacksMu sync.RWMutex
}

// NewXTSConnector creates a new XTS connector
func NewXTSConnector(apiKey, secretKey, source, userID, baseURL string) *XTSConnector {
	return &XTSConnector{
		apiKey:     apiKey,
		secretKey:  secretKey,
		source:     source,
		userID:     userID,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		callbacks:  make(map[string][]MarketDataCallback),
	}
}

// Connect connects to the XTS API
func (c *XTSConnector) Connect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return nil
	}

	// Create login request
	reqURL := c.baseURL + "/user/login"
	reqBody := map[string]string{
		"appKey":    c.apiKey,
		"secretKey": c.secretKey,
		"source":    c.source,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	// Send login request
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader(reqJSON))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	// Parse response
	var loginResp struct {
		Result struct {
			Token string `json:"token"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	// Check response status
	if loginResp.Status != "Success" {
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

	// Store token
	c.token = loginResp.Result.Token
	c.isConnected = true

	// TODO: Start WebSocket connection for real-time data

	return nil
}

// Disconnect disconnects from the XTS API
func (c *XTSConnector) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return nil
	}

	// Create logout request
	reqURL := c.baseURL + "/user/logout"
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	// Send logout request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send logout request: %w", err)
	}
	defer resp.Body.Close()

	// Reset connection state
	c.token = ""
	c.isConnected = false

	// TODO: Close WebSocket connection

	return nil
}

// IsConnected checks if the connector is connected
func (c *XTSConnector) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// GetMarketData gets market data for the specified symbols
func (c *XTSConnector) GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error) {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("not connected to XTS API")
	}
	token := c.token
	c.mutex.RUnlock()

	// Create request URL
	reqURL := c.baseURL + "/marketData/quotes"
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create quotes request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// Add query parameters
	q := req.URL.Query()
	q.Add("symbols", strings.Join(symbols, ","))
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send quotes request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read quotes response: %w", err)
	}

	// Parse response
	var quotesResp struct {
		Result struct {
			Quotes []struct {
				Symbol     string  `json:"symbol"`
				Exchange   string  `json:"exchange"`
				LastPrice  float64 `json:"lastPrice"`
				BidPrice   float64 `json:"bidPrice"`
				AskPrice   float64 `json:"askPrice"`
				BidSize    int     `json:"bidSize"`
				AskSize    int     `json:"askSize"`
				Volume     int     `json:"volume"`
				OpenPrice  float64 `json:"openPrice"`
				HighPrice  float64 `json:"highPrice"`
				LowPrice   float64 `json:"lowPrice"`
				ClosePrice float64 `json:"closePrice"`
				Timestamp  int64   `json:"timestamp"`
			} `json:"quotes"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &quotesResp); err != nil {
		return nil, fmt.Errorf("failed to parse quotes response: %w", err)
	}

	// Check response status
	if quotesResp.Status != "Success" {
		return nil, fmt.Errorf("quotes request failed: %s", quotesResp.Message)
	}

	// Convert to MarketData
	result := make(map[string]MarketData)
	for _, quote := range quotesResp.Result.Quotes {
		result[quote.Symbol] = MarketData{
			Symbol:     quote.Symbol,
			Exchange:   quote.Exchange,
			LastPrice:  quote.LastPrice,
			BidPrice:   quote.BidPrice,
			AskPrice:   quote.AskPrice,
			BidSize:    quote.BidSize,
			AskSize:    quote.AskSize,
			Volume:     quote.Volume,
			OpenPrice:  quote.OpenPrice,
			HighPrice:  quote.HighPrice,
			LowPrice:   quote.LowPrice,
			ClosePrice: quote.ClosePrice,
			Timestamp:  time.Unix(quote.Timestamp/1000, (quote.Timestamp%1000)*1000000),
		}
	}

	return result, nil
}

// GetHistoricalData gets historical data for the specified symbol and interval
func (c *XTSConnector) GetHistoricalData(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error) {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("not connected to XTS API")
	}
	token := c.token
	c.mutex.RUnlock()

	// Create request URL
	reqURL := c.baseURL + "/marketData/history"
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create history request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// Add query parameters
	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("interval", interval)
	q.Add("from", strconv.FormatInt(from.Unix()*1000, 10))
	q.Add("to", strconv.FormatInt(to.Unix()*1000, 10))
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send history request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read history response: %w", err)
	}

	// Parse response
	var historyResp struct {
		Result struct {
			Candles []struct {
				Timestamp int64   `json:"timestamp"`
				Open      float64 `json:"open"`
				High      float64 `json:"high"`
				Low       float64 `json:"low"`
				Close     float64 `json:"close"`
				Volume    int     `json:"volume"`
			} `json:"candles"`
		} `json:"result"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &historyResp); err != nil {
		return nil, fmt.Errorf("failed to parse history response: %w", err)
	}

	// Check response status
	if historyResp.Status != "Success" {
		return nil, fmt.Errorf("history request failed: %s", historyResp.Message)
	}

	// Convert to OHLCV
	result := make([]OHLCV, 0, len(historyResp.Result.Candles))
	for _, candle := range historyResp.Result.Candles {
		result = append(result, OHLCV{
			Symbol:    symbol,
			Interval:  interval,
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
			Timestamp: time.Unix(candle.Timestamp/1000, (candle.Timestamp%1000)*1000000),
		})
	}

	return result, nil
}

// SubscribeToMarketData subscribes to market data for the specified symbols
func (c *XTSConnector) SubscribeToMarketData(ctx context.Context, symbols []string, callback MarketDataCallback) error {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return fmt.Errorf("not connected to XTS API")
	}
	token := c.token
	c.mutex.RUnlock()

	// Register callback
	c.callbacksMu.Lock()
	for _, symbol := range symbols {
		c.callbacks[symbol] = append(c.callbacks[symbol], callback)
	}
	c.callbacksMu.Unlock()

	// Create subscription request
	reqURL := c.baseURL + "/marketData/subscribe"
	reqBody := map[string]interface{}{
		"symbols": symbols,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe request: %w", err)
	}

	// Send subscription request
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create subscribe request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Body = ioutil.NopCloser(bytes.NewReader(reqJSON))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send subscribe request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read subscribe response: %w", err)
	}

	// Parse response
	var subscribeResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &subscribeResp); err != nil {
		return fmt.Errorf("failed to parse subscribe response: %w", err)
	}

	// Check response status
	if subscribeResp.Status != "Success" {
		return fmt.Errorf("subscribe request failed: %s", subscribeResp.Message)
	}

	// TODO: Handle WebSocket messages for subscribed symbols

	return nil
}

// UnsubscribeFromMarketData unsubscribes from market data for the specified symbols
func (c *XTSConnector) UnsubscribeFromMarketData(ctx context.Context, symbols []string) error {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return fmt.Errorf("not connected to XTS API")
	}
	token := c.token
	c.mutex.RUnlock()

	// Unregister callbacks
	c.callbacksMu.Lock()
	for _, symbol := range symbols {
		delete(c.callbacks, symbol)
	}
	c.callbacksMu.Unlock()

	// Create unsubscription request
	reqURL := c.baseURL + "/marketData/unsubscribe"
	reqBody := map[string]interface{}{
		"symbols": symbols,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal unsubscribe request: %w", err)
	}

	// Send unsubscription request
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create unsubscribe request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Body = ioutil.NopCloser(bytes.NewReader(reqJSON))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send unsubscribe request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read unsubscribe response: %w", err)
	}

	// Parse response
	var unsubscribeResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &unsubscribeResp); err != nil {
		return fmt.Errorf("failed to parse unsubscribe response: %w", err)
	}

	// Check response status
	if unsubscribeResp.Status != "Success" {
		return fmt.Errorf("unsubscribe request failed: %s", unsubscribeResp.Message)
	}

	return nil
}

// YahooFinanceConnector implements DataSourceConnector for Yahoo Finance API
type YahooFinanceConnector struct {
	apiKey      string
	isConnected bool
	baseURL     string
	httpClient  *http.Client
	mutex       sync.RWMutex
	callbacks   map[string][]MarketDataCallback
	callbacksMu sync.RWMutex
}

// NewYahooFinanceConnector creates a new Yahoo Finance connector
func NewYahooFinanceConnector(apiKey, baseURL string) *YahooFinanceConnector {
	return &YahooFinanceConnector{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		callbacks:  make(map[string][]MarketDataCallback),
	}
}

// Connect connects to the Yahoo Finance API
func (c *YahooFinanceConnector) Connect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return nil
	}

	// Yahoo Finance API doesn't require a connection, so we just set the flag
	c.isConnected = true
	return nil
}

// Disconnect disconnects from the Yahoo Finance API
func (c *YahooFinanceConnector) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.isConnected = false
	return nil
}

// IsConnected checks if the connector is connected
func (c *YahooFinanceConnector) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// GetMarketData gets market data for the specified symbols
func (c *YahooFinanceConnector) GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error) {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("not connected to Yahoo Finance API")
	}
	c.mutex.RUnlock()

	// Create request URL
	reqURL := c.baseURL + "/v7/finance/quote"
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("symbols", strings.Join(symbols, ","))
	req.URL.RawQuery = q.Encode()

	// Add API key header
	req.Header.Set("X-API-KEY", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send quote request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read quote response: %w", err)
	}

	// Parse response
	var quoteResp struct {
		QuoteResponse struct {
			Result []struct {
				Symbol        string  `json:"symbol"`
				Exchange      string  `json:"exchange"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				Bid           float64 `json:"bid"`
				Ask           float64 `json:"ask"`
				BidSize       int     `json:"bidSize"`
				AskSize       int     `json:"askSize"`
				RegularMarketVolume int     `json:"regularMarketVolume"`
				RegularMarketOpen float64 `json:"regularMarketOpen"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow float64 `json:"regularMarketDayLow"`
				RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
				RegularMarketTime int64   `json:"regularMarketTime"`
			} `json:"result"`
			Error *struct {
				Code string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"quoteResponse"`
	}
	if err := json.Unmarshal(respBody, &quoteResp); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	// Check for errors
	if quoteResp.QuoteResponse.Error != nil {
		return nil, fmt.Errorf("quote request failed: %s - %s", 
			quoteResp.QuoteResponse.Error.Code, 
			quoteResp.QuoteResponse.Error.Description)
	}

	// Convert to MarketData
	result := make(map[string]MarketData)
	for _, quote := range quoteResp.QuoteResponse.Result {
		result[quote.Symbol] = MarketData{
			Symbol:     quote.Symbol,
			E
(Content truncated due to size limit. Use line ranges to read in chunks)